package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pdnode.com/website/handlers"
	"pdnode.com/website/middleware"
	"pdnode.com/website/services"
)

func AuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 依赖注入：DB -> Service -> Handler
	authService := &services.AuthService{DB: db}
	authHandler := &handlers.AuthHandler{Service: authService}

	rg.POST("/login", authHandler.Login)
	protected := rg.Group("")
	protected.Use(middleware.SuperuserAuth())
	{
		protected.POST("/register", authHandler.Register).Use(middleware.SuperuserAuth())

	}
}
