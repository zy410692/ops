package Lib

import (
	"math/rand"
	"strings"
	"time"
)

func GeneratePassword(length int) string {
	// 定义字符集，包括大写字母、小写字母、数字和特殊字符
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_+=<>?/[]{}|"
	var password strings.Builder
	password.Grow(length)

	rand.Seed(time.Now().UnixNano()) // 设置随机种子

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		password.WriteByte(charset[randomIndex])
	}

	return password.String()
}
