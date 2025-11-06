package controller

import (
	"fmt"
	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"
	"sqlite_test/utils"
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

	// 使用密码服务加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}
	user.Password = hashedPassword

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
		constants.SendResponse(c, constants.UserNotFound, nil)
		return
	}

	// 使用密码服务验证密码
	if !utils.CheckPassword(loginReq.Password, user.Password) {
		constants.SendResponse(c, constants.UserUnauthorized, nil)
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
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "存储token失败"})
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

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(c *gin.Context) {
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

	var updateData struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	user.Email = updateData.Email
	if err := database.DB.Save(&user).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}

	user.Password = "" // 清除密码
	constants.SendResponse(c, constants.UserSuccess, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, nil)
		return
	}

	if err := database.DB.Delete(&Model.User{}, userID).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, nil)
		return
	}

	constants.SendResponse(c, constants.UserSuccess, nil)
}
