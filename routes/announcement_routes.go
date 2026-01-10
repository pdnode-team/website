package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pdnode.com/website/handlers"
	"pdnode.com/website/middleware"
	"pdnode.com/website/services"
)

func AnnouncementRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 组装依赖
	service := &services.AnnouncementService{DB: db}
	handler := &handlers.AnnouncementHandler{Service: service}

	// 绑定路由
	rg.GET("/announcements/:id", handler.GetOne)
	rg.GET("/announcements", handler.GetAll)
	protected := rg.Group("")
	protected.Use(middleware.LoginAuth())
	{
		// 这个组里的所有路由都会经过 LoginAuth
		protected.POST("/announcements", handler.Create)
	}
}
