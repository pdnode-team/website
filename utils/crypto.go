package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pepper := GetSuperuserToken()

	// 1. 先进行 Pre-hash
	combinedHash := preHash(password, pepper)

	// 2. 将 SHA-256 的结果交给 Bcrypt
	// 注意：这里不需要再 append pepper 了，因为 preHash 已经混入了
	bytes, err := bcrypt.GenerateFromPassword(combinedHash, 12)
	return string(bytes), err
}

// CheckPasswordHash 改进版（带 Pre-hash）
func CheckPasswordHash(password, hash string) bool {
	pepper := GetSuperuserToken()

	// 1. 用同样的逻辑生成 Pre-hash
	combinedHash := preHash(password, pepper)

	// 2. 验证
	err := bcrypt.CompareHashAndPassword([]byte(hash), combinedHash)
	return err == nil
}
func preHash(password string, pepper []byte) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write(pepper)
	// 返回十六进制字符串，长度固定为 64 位，远低于 Bcrypt 的 72 位限制
	return []byte(hex.EncodeToString(h.Sum(nil)))
}
