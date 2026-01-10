package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pdnode.com/website/models"
	"pdnode.com/website/routes"
)

// Setup JWT Secret Key

// HashPassword 改进版（带 Pre-hash）

func main() {

	// Setup DB
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Announcement{}, &models.User{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	gin.ForceConsoleColor()
	r := routes.SetupRouter(db)

	SetUpSuperuser()

	err = r.Run()
	if err != nil {
		panic(err)
	}

}
