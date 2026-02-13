package Strings

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(storedHash, password string) bool {
	// 格式检查
	parts := strings.Split(storedHash, "_")
	if len(parts) != 2 {
		return false
	}

	salt := parts[0]
	storedHashPart := parts[1]

	// 长度检查
	if len(salt) != 32 || len(storedHashPart) != 64 {
		return false
	}

	// 重新计算
	computedHash := HashPassword(password + salt)

	// 恒定时间比较
	return subtle.ConstantTimeCompare(
		[]byte(storedHashPart),
		[]byte(computedHash),
	) == 1
}
