package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"pdnode.com/website/utils"
)

func SuperuserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Super-Token")

		// 校验 Token
		if token == "" || token != string(utils.GetSuperuserToken()) {
			c.JSON(401, gin.H{"error": "X-Super-Token does not provide or has an incorrect token."})

			c.Abort()
			return
		}

		c.Next()
	}
}

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if !strings.HasPrefix(token, "Bearer ") || token == "" {
			c.JSON(401, gin.H{
				"error": "A Bearer type token is required.",
			})
			c.Abort()
			return
		}

		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

		if token == "" {
			c.JSON(401, gin.H{"error": "Token is empty."})
			c.Abort()
			return
		}

		userID, err := utils.VerifyToken(token)

		if err != nil {
			// 如果有错误，说明验证失败
			c.JSON(401, gin.H{
				"error": err.Error(), // 这里的 err.Error() 会得到 "invalid token" 等信息
			})
			c.Abort()
			return
		}

		c.Set("currentUserID", userID)
		c.Next()
	}
}
