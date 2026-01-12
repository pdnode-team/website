package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pdnode.com/website/config"
	"pdnode.com/website/middleware"
	"pdnode.com/website/models"
	"pdnode.com/website/routes"
	"pdnode.com/website/utils"
)

// Setup JWT Secret Key

// HashPassword 改进版（带 Pre-hash）

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup DB
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})

	config.InitLogger()

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode) // 这样 [GIN-debug] 就全消失了
	}

	slog.Info("Environment Check",
		"ENV_VAR", os.Getenv("ENV"),
		"PORT_VAR", os.Getenv("PORT"),
	)

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Announcement{}, &models.User{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// 4. 创建 Gin 实例 (不使用 Default 以免日志冲突)
	r := gin.New()
	gin.ForceConsoleColor()

	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery()) // 建议加上，防止 panic 导致程序退出

	SetUpSuperuser()

	utils.InitAuth()

	routes.SetupRouter(r, db)

	// 8. 启动
	if err := r.Run(); err != nil {
		panic(err)
	}

}
