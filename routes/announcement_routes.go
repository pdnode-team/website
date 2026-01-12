package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	// 这里是关键：定义别名为 mgin，解决包名冲突
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
	"pdnode.com/website/handlers"
	"pdnode.com/website/middleware"
	"pdnode.com/website/services"
)

func AnnouncementRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 1. 组装依赖
	service := services.NewAnnouncementService(db)
	handler := &handlers.AnnouncementHandler{Service: service}

	// 2. 配置限速策略
	// 策略 A: 针对普通查询（防止爬虫或恶意刷新）
	publicRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  30, // 每分钟 30 次
	}
	publicStore := memory.NewStore()
	publicInstance := limiter.New(publicStore, publicRate)
	publicMiddleware := mgin.NewMiddleware(publicInstance)

	// 策略 B: 针对管理操作（防止误操作或脚本爆破）
	adminRate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  100, // 每小时 100 次修改
	}
	adminStore := memory.NewStore()
	adminInstance := limiter.New(adminStore, adminRate)
	adminMiddleware := mgin.NewMiddleware(adminInstance)

	// 3. 绑定路由
	announcementGroup := rg.Group("/announcements")
	{
		// 公开查询接口：增加限速
		announcementGroup.GET("", publicMiddleware, handler.GetAll)
		announcementGroup.GET("/:id", publicMiddleware, handler.GetOne)

		// 需要登录权限的操作
		protected := announcementGroup.Group("")
		protected.Use(middleware.LoginAuth()) // 先验证登录
		protected.Use(adminMiddleware)        // 再检查管理操作频率
		{
			protected.POST("", handler.Create)
			protected.PUT("/:id", handler.Update)
			protected.DELETE("/:id", handler.Delete)
		}
	}
}
