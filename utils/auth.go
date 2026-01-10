package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte("your_secret_key_pdnode")

type UserClaims struct {
	UserID               uuid.UUID `json:"user_id"`
	Email                string    `json:"email"`
	jwt.RegisteredClaims           // 包含过期时间、发行人等基础字段
}

func GenerateToken(userID uuid.UUID) (string, error) {
	// 1. 设置过期时间（比如 24 小时后）
	expirationTime := time.Now().Add(24 * time.Hour)

	// 2. 创建 Payload (Claims)
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签名时间
		},
	}

	// 3. 使用指定的签名方法创建 Token 对象 (常用 HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. 使用密钥进行数字签名，生成最终字符串
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func VerifyToken(tokenString string) (uuid.UUID, error) {
	// 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 1. 验证签名算法是否为你预期的（防止安全漏洞）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 2. 返回你的密钥
		return jwtKey, nil
	})

	// 3. 判断是否合法
	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("failed to parse claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid format: %v", err)
	}

	return userID, nil
}
