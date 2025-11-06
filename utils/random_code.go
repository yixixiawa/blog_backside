package utils

import (
	"math/rand"
	"time"
)

const (
	// 数字
	NumberChars = "0123456789"
	// 字母
	LetterChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 数字+字母
	MixedChars = NumberChars + LetterChars
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateMixedCode 生成数字+字母验证码（项目唯一导出验证码生成器）
// 保留此函数供全项目使用，移除其他不需要的导出函数以减少混淆。
func GenerateMixedCode(length int) string {
	return generateRandomCode(length, MixedChars)
}

// generateRandomCode 内部生成随机码函数
func generateRandomCode(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
