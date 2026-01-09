package main

import (
  "time"

  "github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
)

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
