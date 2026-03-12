package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"blog/Model"
	"blog/database"
	"blog/utils"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// OAuthService OAuth服务
type OAuthService struct {
	db *gorm.DB
}

// NewOAuthService 创建OAuth服务
func NewOAuthService() *OAuthService {
	return &OAuthService{
		db: database.DB,
	}
}

// GetEnabledPlatforms 获取所有启用的平台
func (s *OAuthService) GetEnabledPlatforms() ([]Model.OAuthPlatform, error) {
	var platforms []Model.OAuthPlatform
	err := s.db.Where("is_enabled = ?", true).
		Order("sort_order asc, created_at asc").
		Find(&platforms).Error
	return platforms, err
}

// GetPlatformByID 根据ID获取平台
func (s *OAuthService) GetPlatformByID(platformID uint) (*Model.OAuthPlatform, error) {
	var platform Model.OAuthPlatform
	err := s.db.First(&platform, platformID).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

// GetPlatformByName 根据名称获取平台
func (s *OAuthService) GetPlatformByName(platform string) (*Model.OAuthPlatform, error) {
	var platformModel Model.OAuthPlatform
	err := s.db.Where("platform = ? AND is_enabled = ?", platform, true).First(&platformModel).Error
	if err != nil {
		return nil, err
	}
	return &platformModel, nil
}

// CreateOAuthConfig 创建OAuth2.0配置
func (s *OAuthService) CreateOAuthConfig(platform *Model.OAuthPlatform) *oauth2.Config {
	scopes := utils.ParseScopes(platform.Scopes)

	return &oauth2.Config{
		ClientID:     platform.ClientID,
		ClientSecret: platform.ClientSecret,
		RedirectURL:  platform.RedirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  platform.AuthURL,
			TokenURL: platform.TokenURL,
		},
	}
}

// SaveState 保存state
func (s *OAuthService) SaveState(state string, userID uint, platformID uint, expiresIn time.Duration) error {
	oauthState := &Model.OAuthState{
		State:      state,
		UserID:     userID,
		PlatformID: platformID,
		ExpiresAt:  time.Now().Add(expiresIn),
	}
	return s.db.Create(oauthState).Error
}

// VerifyState 验证state
func (s *OAuthService) VerifyState(state string) (*Model.OAuthState, error) {
	var oauthState Model.OAuthState
	err := s.db.Where("state = ? AND expires_at > ?", state, time.Now()).
		Preload("Platform").
		First(&oauthState).Error

	if err != nil {
		return nil, errors.New("无效的state或已过期")
	}

	// 使用后删除
	s.db.Delete(&oauthState)
	return &oauthState, nil
}

// BindOAuthAccount 绑定第三方账号
func (s *OAuthService) BindOAuthAccount(userID uint, platformID uint, platformUserID string, userInfo map[string]interface{}, token *oauth2.Token) error {
	// 检查该第三方账号是否已被其他用户绑定
	var count int64
	s.db.Model(&Model.OAuthAccount{}).
		Where("platform_id = ? AND platform_user_id = ?", platformID, platformUserID).
		Count(&count)

	if count > 0 {
		return errors.New("该第三方账号已被其他用户绑定")
	}

	// 将userInfo转为JSON
	rawData, _ := json.Marshal(userInfo)

	// 提取常用信息
	platformUserName, _ := userInfo["name"].(string)
	if platformUserName == "" {
		platformUserName, _ = userInfo["login"].(string) // GitHub使用login
	}
	platformUserEmail, _ := userInfo["email"].(string)
	avatarURL, _ := userInfo["avatar_url"].(string)
	if avatarURL == "" {
		avatarURL, _ = userInfo["picture"].(string) // Google使用picture
	}

	// 创建绑定
	account := &Model.OAuthAccount{
		UserID:            userID,
		PlatformID:        platformID,
		PlatformUserID:    platformUserID,
		PlatformUserName:  platformUserName,
		PlatformUserEmail: platformUserEmail,
		AvatarURL:         avatarURL,
		AccessToken:       token.AccessToken,
		RefreshToken:      token.RefreshToken,
		TokenExpiresAt:    &token.Expiry,
		RawData:           string(rawData),
	}

	return s.db.Create(account).Error
}

// GetUserByOAuth 通过第三方账号获取用户
func (s *OAuthService) GetUserByOAuth(platformID uint, platformUserID string) (*Model.User, error) {
	var account Model.OAuthAccount
	err := s.db.Where("platform_id = ? AND platform_user_id = ?", platformID, platformUserID).
		Preload("User").
		First(&account).Error

	if err != nil {
		return nil, err
	}

	return &account.User, nil
}

// GetUserOAuthAccounts 获取用户绑定的所有第三方账号
func (s *OAuthService) GetUserOAuthAccounts(userID uint) ([]Model.OAuthAccount, error) {
	var accounts []Model.OAuthAccount
	err := s.db.Where("user_id = ?", userID).
		Preload("Platform").
		Find(&accounts).Error

	// 清理敏感信息
	for i := range accounts {
		accounts[i].AccessToken = ""
		accounts[i].RefreshToken = ""
		accounts[i].RawData = ""
	}

	return accounts, err
}

// UnbindOAuthAccount 解绑第三方账号
func (s *OAuthService) UnbindOAuthAccount(userID uint, platformID uint) error {
	// 检查用户是否有其他登录方式
	var count int64
	s.db.Model(&Model.OAuthAccount{}).Where("user_id = ?", userID).Count(&count)

	var user Model.User
	s.db.First(&user, userID)

	// 如果这是最后一个登录方式且用户没有密码，不允许解绑
	if count <= 1 && user.Password == "" {
		return errors.New("无法解绑最后一个登录方式，请先设置密码")
	}

	// 删除绑定
	result := s.db.Where("user_id = ? AND platform_id = ?", userID, platformID).Delete(&Model.OAuthAccount{})
	if result.RowsAffected == 0 {
		return errors.New("未找到绑定关系")
	}

	return nil
}

// CreateOrUpdateUser 创建或更新用户（根据第三方信息）
func (s *OAuthService) CreateOrUpdateUser(platform *Model.OAuthPlatform, userInfo map[string]interface{}, platformUserID string) (*Model.User, error) {
	var user Model.User

	// 提取用户信息
	email, _ := userInfo["email"].(string)
	name, _ := userInfo["name"].(string)
	if name == "" {
		name, _ = userInfo["login"].(string)
	}
	avatar, _ := userInfo["avatar_url"].(string)
	if avatar == "" {
		avatar, _ = userInfo["picture"].(string)
	}

	// 优先通过email查找用户
	if email != "" {
		err := s.db.Where("email = ?", email).First(&user).Error
		if err == nil {
			return &user, nil
		}
	}

	// 创建新用户
	username := fmt.Sprintf("%s_%s", platform.Platform, platformUserID)
	// 如果用户名已存在，添加随机后缀
	var exists bool
	for i := 0; i < 10; i++ {
		var count int64
		s.db.Model(&Model.User{}).Where("username = ?", username).Count(&count)
		if count == 0 {
			exists = false
			break
		}
		username = fmt.Sprintf("%s_%s_%d", platform.Platform, platformUserID, time.Now().UnixNano()%1000)
		exists = true
	}

	if exists {
		return nil, errors.New("无法生成唯一的用户名")
	}

	newUser := Model.User{
		Username: username,
		Email:    email,
		Avatar:   avatar,
		Password: "", // 第三方登录用户初始无密码
	}

	err := s.db.Create(&newUser).Error
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}
