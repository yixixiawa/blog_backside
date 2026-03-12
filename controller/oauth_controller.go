package controller

import (
	"blog/Model"
	"blog/constants"
	"blog/database"
	"blog/service"
	"blog/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// oauthService 全局OAuth服务实例
var oauthService = service.NewOAuthService()

func GetOAuthPlatforms(c *gin.Context) {
	platforms, err := oauthService.GetEnabledPlatforms()
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "获取平台列表失败"})
		return
	}

	constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{
		"platforms": platforms,
	})
}

// GET /oauth/login/:platform  (例如 /oauth/login/github)
func OAuthLogin(c *gin.Context) {
	platformName := c.Param("platform")

	// 1. 根据平台名获取平台配置
	platform, err := oauthService.GetPlatformByName(platformName)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthNotFound, gin.H{"error": "不支持的OAuth平台"})
		return
	}

	// 2. 创建OAuth2配置
	oauthConfig := oauthService.CreateOAuthConfig(platform)

	// 3. 生成随机state（防止CSRF攻击）
	state := utils.GenerateMixedCode(32)

	// 4. 保存state到数据库（10分钟有效期）
	if err := oauthService.SaveState(state, 0, platform.OAuthID, 10*time.Minute); err != nil {
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "生成认证状态失败"})
		return
	}

	// 5. 生成授权URL并重定向
	authURL := oauthConfig.AuthCodeURL(state)

	// 返回授权URL让前端跳转（或服务端直接重定向）
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "请跳转到授权URL",
		"auth_url": authURL,
	})
}

// ============================================================
// 3. OAuth回调处理 —— GitHub回调到此接口
// ============================================================

// OAuthCallback 处理OAuth回调
// GET /oauth/callback/:platform  (例如 /oauth/callback/github)
func OAuthCallback(c *gin.Context) {
	platformName := c.Param("platform")

	// 1. 获取回调参数
	code := c.Query("code")
	state := c.Query("state")
	callbackError := c.Query("error")

	// 如果GitHub返回了错误（用户拒绝授权等）
	if callbackError != "" {
		errorDesc := c.Query("error_description")
		constants.SendOAuthResponse(c, constants.OAuthCallbackError, gin.H{
			"error":       callbackError,
			"description": errorDesc,
		})
		return
	}

	if code == "" || state == "" {
		constants.SendOAuthResponse(c, constants.OAuthBadRequest, gin.H{"error": "缺少code或state参数"})
		return
	}

	// 2. 验证state（防止CSRF）
	oauthState, err := oauthService.VerifyState(state)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthUnauthorized, gin.H{"error": "无效的state，可能已过期或被篡改"})
		return
	}

	// 3. 获取平台配置
	platform, err := oauthService.GetPlatformByName(platformName)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthNotFound, gin.H{"error": "不支持的OAuth平台"})
		return
	}

	// 4. 创建OAuth2配置
	oauthConfig := oauthService.CreateOAuthConfig(platform)

	// 5. 用code换取access_token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthCallbackError, gin.H{"error": "获取access_token失败: " + err.Error()})
		return
	}

	// 6. 使用access_token获取用户信息
	userInfo, err := fetchGitHubUserInfo(token.AccessToken)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthCallbackError, gin.H{"error": "获取用户信息失败: " + err.Error()})
		return
	}

	// 7. 提取平台用户ID
	platformUserID := extractPlatformUserID(platformName, userInfo)
	if platformUserID == "" {
		constants.SendOAuthResponse(c, constants.OAuthCallbackError, gin.H{"error": "无法获取第三方平台用户ID"})
		return
	}

	// 8. 判断逻辑：已登录用户 → 绑定账号 / 未登录用户 → 登录或注册
	if oauthState.UserID > 0 {
		// 已登录用户绑定第三方账号
		err = oauthService.BindOAuthAccount(oauthState.UserID, platform.OAuthID, platformUserID, userInfo, token)
		if err != nil {
			constants.SendOAuthResponse(c, constants.OAuthConflict, gin.H{"error": err.Error()})
			return
		}
		constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{"message": "第三方账号绑定成功"})
		return
	}

	// 未登录：尝试查找已绑定的用户，没有则自动创建
	user, err := oauthService.GetUserByOAuth(platform.OAuthID, platformUserID)
	if err != nil {
		// 用户不存在，自动创建新用户
		user, err = oauthService.CreateOrUpdateUser(platform, userInfo, platformUserID)
		if err != nil {
			constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "创建用户失败: " + err.Error()})
			return
		}

		// 绑定OAuth账号到新用户
		if err := oauthService.BindOAuthAccount(user.UserID, platform.OAuthID, platformUserID, userInfo, token); err != nil {
			constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "绑定账号失败: " + err.Error()})
			return
		}
	}

	// 9. 为用户生成JWT Token
	jwtToken, err := utils.GenerateToken(int64(user.UserID), user.Username)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "生成token失败"})
		return
	}

	// 10. 存储token到Redis
	tokenKey := fmt.Sprintf("user_token:%d", user.UserID)
	if err := database.SetString(tokenKey, jwtToken, 24*time.Hour); err != nil {
		fmt.Printf("Redis SetString Error: %v\n", err)
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "存储token失败"})
		return
	}

	user.Password = "" // 清除密码
	constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{
		"user":  user,
		"token": jwtToken,
	})
}

// ============================================================
// 4. 已登录用户绑定OAuth账号（需要JWT认证）
// ============================================================

// OAuthBind 已登录用户主动绑定第三方账号
// GET /oauth/bind/:platform
func OAuthBind(c *gin.Context) {
	platformName := c.Param("platform")

	// 获取当前登录用户ID
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		constants.SendOAuthResponse(c, constants.OAuthUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	currentUserID := uint(currentUserIDVal.(int64))

	// 获取平台配置
	platform, err := oauthService.GetPlatformByName(platformName)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthNotFound, gin.H{"error": "不支持的OAuth平台"})
		return
	}

	// 创建OAuth2配置
	oauthConfig := oauthService.CreateOAuthConfig(platform)

	// 生成state（带上用户ID，回调时用于绑定）
	state := utils.GenerateMixedCode(32)
	if err := oauthService.SaveState(state, currentUserID, platform.OAuthID, 10*time.Minute); err != nil {
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "生成认证状态失败"})
		return
	}

	authURL := oauthConfig.AuthCodeURL(state)
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "请跳转到授权URL完成绑定",
		"auth_url": authURL,
	})
}

// ============================================================
// 5. 解绑OAuth账号（需要JWT认证）
// ============================================================

// OAuthUnbind 解绑第三方账号
// DELETE /oauth/unbind/:platform
func OAuthUnbind(c *gin.Context) {
	platformName := c.Param("platform")

	// 获取当前登录用户ID
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		constants.SendOAuthResponse(c, constants.OAuthUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	currentUserID := uint(currentUserIDVal.(int64))

	// 获取平台
	platform, err := oauthService.GetPlatformByName(platformName)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthNotFound, gin.H{"error": "不支持的OAuth平台"})
		return
	}

	// 解绑
	if err := oauthService.UnbindOAuthAccount(currentUserID, platform.OAuthID); err != nil {
		constants.SendOAuthResponse(c, constants.OAuthBadRequest, gin.H{"error": err.Error()})
		return
	}

	constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{"message": "解绑成功"})
}

// ============================================================
// 6. 获取当前用户绑定的OAuth账号列表（需要JWT认证）
// ============================================================

// GetUserOAuthAccounts 获取当前用户绑定的第三方账号
// GET /oauth/accounts
func GetUserOAuthAccounts(c *gin.Context) {
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		constants.SendOAuthResponse(c, constants.OAuthUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	currentUserID := uint(currentUserIDVal.(int64))

	accounts, err := oauthService.GetUserOAuthAccounts(currentUserID)
	if err != nil {
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "获取绑定列表失败"})
		return
	}

	constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{
		"accounts": accounts,
	})
}

// ============================================================
// 7. 初始化OAuth平台配置（管理员接口 / 首次部署时使用）
// ============================================================

// InitGitHubPlatform 初始化GitHub OAuth平台配置
// POST /oauth/admin/init-github
func InitGitHubPlatform(c *gin.Context) {
	// 检查管理员权限
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		constants.SendOAuthResponse(c, constants.OAuthUnauthorized, nil)
		return
	}
	if !checkIsAdmin(uint(currentUserIDVal.(int64))) {
		constants.SendOAuthResponse(c, constants.OAuthForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var req struct {
		ClientID     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
		RedirectURL  string `json:"redirect_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		constants.SendOAuthResponse(c, constants.OAuthBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 检查是否已存在
	var count int64
	database.DB.Model(&Model.OAuthPlatform{}).Where("platform = ?", "github").Count(&count)
	if count > 0 {
		// 更新已有配置
		if err := database.DB.Model(&Model.OAuthPlatform{}).Where("platform = ?", "github").Updates(map[string]interface{}{
			"client_id":     req.ClientID,
			"client_secret": req.ClientSecret,
			"redirect_url":  req.RedirectURL,
		}).Error; err != nil {
			constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "更新平台配置失败"})
			return
		}
		constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{"message": "GitHub OAuth配置已更新"})
		return
	}

	// 创建新配置
	platform := Model.OAuthPlatform{
		Platform:     "github",
		DisplayName:  "GitHub",
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURL:  req.RedirectURL,
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		UserInfoURL:  "https://api.github.com/user",
		Scopes:       "user:email",
		IconURL:      "https://github.githubassets.com/favicons/favicon.svg",
		SortOrder:    1,
		IsEnabled:    true,
	}

	if err := database.DB.Create(&platform).Error; err != nil {
		constants.SendOAuthResponse(c, constants.OAuthSystemError, gin.H{"error": "创建平台配置失败"})
		return
	}

	constants.SendOAuthResponse(c, constants.OAuthSuccess, gin.H{
		"message":  "GitHub OAuth平台配置成功",
		"platform": platform,
	})
}

// ============================================================
// 内部辅助函数
// ============================================================

// fetchGitHubUserInfo 调用GitHub API获取用户信息
func fetchGitHubUserInfo(accessToken string) (map[string]interface{}, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求GitHub API失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API返回错误(%d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	// 额外获取用户邮箱（GitHub用户可能设置邮箱为私有）
	email, _ := userInfo["email"].(string)
	if email == "" {
		emailInfo, err := fetchGitHubUserEmails(accessToken)
		if err == nil && emailInfo != "" {
			userInfo["email"] = emailInfo
		}
	}

	return userInfo, nil
}

// fetchGitHubUserEmails 获取GitHub用户的邮箱列表
func fetchGitHubUserEmails(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取邮箱列表失败: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	// 优先返回主邮箱
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// 返回第一个已验证的邮箱
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", nil
}

// extractPlatformUserID 从用户信息中提取平台用户ID
func extractPlatformUserID(platform string, userInfo map[string]interface{}) string {
	switch platform {
	case "github":
		// GitHub的用户ID是数字，转为字符串
		if id, ok := userInfo["id"].(float64); ok {
			return fmt.Sprintf("%.0f", id)
		}
	case "google":
		if id, ok := userInfo["id"].(string); ok {
			return id
		}
	}
	return ""
}
