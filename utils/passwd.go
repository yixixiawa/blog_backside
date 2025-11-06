package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	// 使用默认的cost（通常是10-12之间）
	// cost越高越安全，但也越慢
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPasswordWithCost 使用自定义cost加密密码
func HashPasswordWithCost(password string, cost int) (string, error) {
	// cost建议范围：10-15
	// cost=10: ~70ms, cost=12: ~250ms, cost=14: ~1s
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
