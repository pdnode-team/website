package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"pdnode.com/website/models"
	"pdnode.com/website/services"
)

type AuthHandler struct {
	// 以前是 *services.AuthService，现在改为接口
	Service services.AuthServiceInterface
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.Debug("request bind failed",
			"path", c.Request.URL.Path,
			"error", err.Error(),
			"ip", c.ClientIP(),
		)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		slog.Warn("Login failed",
			"email", input.Email,
			"reason", "invalid credentials", // 明确失败原因
			"ip", c.ClientIP(),
		)
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	slog.Info("user login success",
		"email", user.Email,
		"user_id", user.ID,
		"ip", c.ClientIP(),
	)

	c.JSON(200, gin.H{"message": "success", "token": token, "user": user})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Register(input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "success"})
}
