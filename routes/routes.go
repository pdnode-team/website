package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置所有的路由映射
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	rg := r.Group("/")

	AnnouncementRoutes(rg, db)
	AuthRoutes(rg, db)

	// 基础路由
	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	rg.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "Pdnode Website API is running",
		})
	})

	return r
}
