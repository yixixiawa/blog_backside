package service

import (
	"fmt"
	"sqlite_test/database"
	"sqlite_test/utils"
	"time"
)

type EmailService struct{}

var defaultEmailService = &EmailService{}

// GenerateAndStoreVerificationCode 生成并存储验证码
func GenerateAndStoreVerificationCode(email string) (string, error) {
	// 检查是否存在未过期的验证码
	exists, _ := CheckVerificationCodeExists(email)
	if exists {
		return "", fmt.Errorf("验证码已发送，请等待之前的验证码过期")
	}

	// 生成新验证码
	code := utils.GenerateMixedCode(10)

	// 使用Redis存储验证码
	key := fmt.Sprintf("email:verify:%s", email)
	err := database.SetString(key, code, 5*time.Minute)
	if err != nil {
		return "", fmt.Errorf("存储验证码失败: %v", err)
	}

	return code, nil
}

// VerifyEmailCode 验证邮箱验证码
func VerifyEmailCode(email, code string) (bool, error) {
	key := fmt.Sprintf("email:verify:%s", email)

	storedCode, err := database.GetString(key)
	if err != nil {
		return false, nil // 验证码不存在或已过期
	}

	if storedCode == code {
		// 验证成功后立即删除验证码
		_ = database.Delete(key)
		return true, nil
	}

	return false, nil
}

// CheckVerificationCodeExists 检查验证码是否存在
func CheckVerificationCodeExists(email string) (bool, error) {
	key := fmt.Sprintf("email:verify:%s", email)
	return database.Exists(key)
}

// GetVerificationCodeTTL 获取验证码剩余有效期
func GetVerificationCodeTTL(email string) (time.Duration, error) {
	key := fmt.Sprintf("email:verify:%s", email)
	return database.GetTTL(key)
}

// DeleteVerificationCode 删除验证码
func DeleteVerificationCode(email string) error {
	key := fmt.Sprintf("email:verify:%s", email)
	return database.Delete(key)
}
