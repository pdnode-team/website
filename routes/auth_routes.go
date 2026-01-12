package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
	"pdnode.com/website/handlers"
	"pdnode.com/website/middleware"
	"pdnode.com/website/services"
)

func AuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 1. 依赖注入
	service := services.NewAuthService(db)
	handler := &handlers.AuthHandler{Service: service}

	// 2. 定义登录限速：防止暴力破解 (例如：每分钟允许尝试 5 次)
	loginRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}
	loginStore := memory.NewStore()
	loginInstance := limiter.New(loginStore, loginRate)
	loginLimitMiddleware := mgin.NewMiddleware(loginInstance)

	// 3. 路由绑定
	authGroup := rg.Group("/")
	{
		// 登录接口：必须加限速！
		authGroup.POST("/login", loginLimitMiddleware, handler.Login)

		// 注册接口：由超级管理员控制
		protected := authGroup.Group("")
		protected.Use(middleware.SuperuserAuth()) // 整个组开启超管验证
		{
			// 这里不需要再写 .Use(middleware.SuperuserAuth())，因为父组已经包含了
			protected.POST("/register", handler.Register)
		}
	}
}
