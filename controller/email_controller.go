package controller

import (
	"sqlite_test/constants"
	"sqlite_test/service"

	"github.com/gin-gonic/gin"
)

// SendVerificationEmail 发送验证邮件
func SendVerificationEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.GenerateAndStoreVerificationCode(req.Email)

	constants.SendResponse(c, constants.Success, gin.H{
		"message":   "验证码已发送",
		"expire_in": 300, // 5分钟有效期
	})

	if ttl, _ := service.GetVerificationCodeTTL(req.Email); ttl > 0 {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error": "请求过于频繁",
			"wait":  int(ttl.Seconds()),
		})
		return
	}
}

// CheckVerificationCode 检查验证码
func CheckVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证验证码
	ok, err := service.VerifyEmailCode(req.Email, req.Code)
	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "验证码验证失败"})
		return
	}
	if !ok {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "验证码错误或已过期"})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{
		"message": "验证成功",
	})
}
