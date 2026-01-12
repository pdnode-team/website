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
	service := services.NewAuthService(db)

	handler := &handlers.AuthHandler{Service: service}

	rg.POST("/login", handler.Login)
	protected := rg.Group("")
	protected.Use(middleware.SuperuserAuth())
	{
		protected.POST("/register", handler.Register).Use(middleware.SuperuserAuth())

	}
}
