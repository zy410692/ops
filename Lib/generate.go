package Lib

import (
	"math/rand"
	"strings"
	"time"
)

func GeneratePassword(length int) string {
	// 将字符集分开，便于确保每种类型的字符都被使用
	const (
		lowerChars = "abcdefghijklmnopqrstuvwxyz"
		upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers    = "0123456789"
		specials   = "!@#$%^&*()-_+=<>?/[]{}|"
	)

	if length < 8 {
		length = 8 // 确保密码长度至少为8，以包含所有类型的字符
	}

	var password strings.Builder
	password.Grow(length)

	rand.Seed(time.Now().UnixNano())

	// 确保至少包含一个特殊字符
	password.WriteByte(specials[rand.Intn(len(specials))])
	// 确保至少包含一个数字
	password.WriteByte(numbers[rand.Intn(len(numbers))])
	// 确保至少包含一个大写字母
	password.WriteByte(upperChars[rand.Intn(len(upperChars))])
	// 确保至少包含一个小写字母
	password.WriteByte(lowerChars[rand.Intn(len(lowerChars))])

	// 合并所有字符集
	allChars := lowerChars + upperChars + numbers + specials

	// 填充剩余长度
	for i := 4; i < length; i++ {
		randomIndex := rand.Intn(len(allChars))
		password.WriteByte(allChars[randomIndex])
	}

	// 将生成的密码转换为字节切片并打乱顺序
	result := []byte(password.String())
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return string(result)
}
