package handlers

import (
	"github.com/gin-gonic/gin"
	"pdnode.com/website/models"
	"pdnode.com/website/services"
	"pdnode.com/website/utils"
)

type AuthHandler struct {
	Service *services.AuthService
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success", "token": token, "user": user})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Handler 负责处理 Header 这种 HTTP 特有的东西
	if c.GetHeader("X-Super-Token") != string(utils.GetSuperuserToken()) {
		c.JSON(403, gin.H{"error": "Invalid token"})
		return
	}

	if err := h.Service.Register(input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "success"})
}
