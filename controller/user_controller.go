package controller

import (
	"blog/Model"
	"blog/constants"
	"blog/database"
	"blog/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用户注册
func UserRegister(c *gin.Context) {
	var user Model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	// 检查用户名是否已存在
	if err := database.DB.Where("username = ?", user.Username).First(&Model.User{}).Error; err == nil {
		constants.SendResponse(c, constants.UserConflict, nil)
		return
	}

	// 使用密码服务加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}
	user.Password = hashedPassword
	user.IsAdmin = false

	// 保存到数据库
	if err := database.DB.Create(&user).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}

	user.Password = "" // 清除密码
	constants.SendResponse(c, constants.UserSuccess, user)
}

// 用户登录
func UserLogin(c *gin.Context) {
	var loginReq Model.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	var user Model.User
	if err := database.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		constants.SendResponse(c, constants.UserLoginError, nil)
		return
	}

	// 使用密码服务验证密码
	if !utils.CheckPassword(loginReq.Password, user.Password) {
		constants.SendResponse(c, constants.UserLoginError, nil)
		return
	}

	// 生成JWT Token（确保 GenerateToken 接收 (int64, string)）
	token, err := utils.GenerateToken(int64(user.UserID), user.Username)
	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "生成token失败"})
		return
	}

	// 存储 token 到 Redis，键名示例：user_token:<userID>
	tokenKey := fmt.Sprintf("user_token:%d", user.UserID)
	// 设置过期时间，与 token 的有效期保持一致（这里示例 24 小时）
	if err := database.SetString(tokenKey, token, 24*time.Hour); err != nil {
		fmt.Printf("Redis SetString Error: %v\n", err) // 👈 新增：打印具体错误到控制台
		constants.SendResponse(c, constants.UserRedisError, gin.H{"error": "存储token失败"})
		return
	}

	user.Password = "" // 清除密码
	constants.SendResponse(c, constants.UserSuccess, gin.H{
		"user":  user,
		"token": token,
	})
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	// 如果使用了会话或令牌，在这里清除它们
	constants.SendResponse(c, constants.UserSuccess, nil)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	type PasswordChange struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var pwdChange PasswordChange
	if err := c.ShouldBindJSON(&pwdChange); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}

	// 从认证中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{
			"error": "未登录或会话已过期",
		})
		return
	}

	var user Model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 验证旧密码
	if utils.CheckPassword(pwdChange.OldPassword, user.Password) {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{
			"error": "原密码错误",
		})
		return
	}

	// 加密新密码
	hashedNewPassword, err := utils.HashPassword(pwdChange.NewPassword)
	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{
			"error": "密码加密失败",
		})
		return
	}

	// 更新密码
	user.Password = hashedNewPassword
	if err := database.DB.Save(&user).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{
			"error": "密码更新失败",
		})
		return
	}

	constants.SendResponse(c, constants.UserSuccess, gin.H{
		"message": "密码修改成功",
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	var user Model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, nil)
		return
	}

	user.Password = "" // 清除密码
	constants.SendResponse(c, constants.UserSuccess, user)
}

// checkIsAdmin 检查用户是否为管理员
func checkIsAdmin(userID uint) bool {
	var user Model.User
	if err := database.DB.Select("is_admin").First(&user, userID).Error; err != nil {
		return false
	}
	return user.IsAdmin
}

// ListUsers 获取所有用户列表（仅限管理员）
func ListUsers(c *gin.Context) {
	// 获取当前登录用户ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, nil)
		return
	}

	// 检查管理员权限
	if !checkIsAdmin(uint(currentUserID.(int64))) {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var users []Model.User
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	database.DB.Model(&Model.User{}).Count(&total)

	if err := database.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	constants.SendResponse(c, constants.UserSuccess, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(c *gin.Context) {
	id := c.Param("id")
	targetUserID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	// 获取当前登录用户ID
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, nil)
		return
	}
	currentUserID := uint(currentUserIDVal.(int64))

	// 权限检查：只有管理员或用户自己可以修改
	if currentUserID != uint(targetUserID) {
		if !checkIsAdmin(currentUserID) {
			constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权修改此用户信息"})
			return
		}
	}

	var user Model.User
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, nil)
		return
	}

	var updateData struct {
		Email   string `json:"email"`
		Avatar  string `json:"avatar"`
		IsAdmin *bool  `json:"is_admin"` // 使用指针以区分是否传递了该字段
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	updates := make(map[string]interface{})
	if updateData.Email != "" {
		updates["email"] = updateData.Email
	}
	if updateData.Avatar != "" {
		updates["avatar"] = updateData.Avatar
	}

	// 只有管理员可以修改 IsAdmin 状态
	if updateData.IsAdmin != nil {
		if checkIsAdmin(currentUserID) {
			updates["is_admin"] = *updateData.IsAdmin
		} else {
			// 如果非管理员尝试修改权限，可以选择忽略或报错，这里选择忽略并记录日志或仅忽略
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
			constants.SendResponse(c, constants.UserSystemError, nil)
			return
		}
	}

	user.Password = "" // 清除密码
	constants.SendResponse(c, constants.UserSuccess, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	targetUserID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	// 获取当前登录用户ID
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, nil)
		return
	}
	currentUserID := uint(currentUserIDVal.(int64))

	// 权限检查：只有管理员可以删除用户（或者用户注销自己，视需求而定）
	// 这里假设只有管理员可以删除其他用户，用户可以删除自己
	if currentUserID != uint(targetUserID) {
		if !checkIsAdmin(currentUserID) {
			constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权删除此用户"})
			return
		}
	}

	if err := database.DB.Delete(&Model.User{}, targetUserID).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}

	constants.SendResponse(c, constants.UserSuccess, nil)
}
