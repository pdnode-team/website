package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup JWT Secret Key
var jwtKey = []byte("your_secret_key_pdnode")

// HashPassword 1. 加密密码（注册时使用）
func HashPassword(password string) (string, error) {
	// 强度系数默认为 10，数值越大越慢（越安全）
	combined := append([]byte(password), GetSuperuserToken()...)
	bytes, err := bcrypt.GenerateFromPassword(combined, bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 2. 比对密码（登录时使用）
func CheckPasswordHash(password, hash string) bool {
	combined := append([]byte(password), GetSuperuserToken()...)
	err := bcrypt.CompareHashAndPassword([]byte(hash), combined)
	return err == nil
}

func main() {

	// Setup DB
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Announcement{}, &User{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	gin.ForceConsoleColor()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "Pdnode Website API is running",
		})
	})

	r.GET("/announcements/:id", func(c *gin.Context) {
		id := c.Param("id")

		var a Announcement

		result := db.First(&a, id)

		if result.Error != nil {

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{
					"message": "Announcement not found",
				})
			}

			c.JSON(500, gin.H{
				"message": "Something went wrong",
			})

			println("[Server Error] Query Error: " + result.Error.Error())

			return
		}

		c.JSON(200, a)
		return

	})

	r.GET("/announcements", func(c *gin.Context) {
		var a []Announcement
		result := db.Order("created_at desc").Find(&a)
		if result.Error != nil {
			c.JSON(500, gin.H{
				"error": result.Error.Error(),
			})
			return
		}

		c.JSON(200, a)

	})

	r.POST("/login", func(c *gin.Context) {
		var input LoginRequest

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var user User
		const fakeBCryptHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgNIhp.qfTMvYeJLvAbpE52EPIvG"

		result := db.Where("email = ?", input.Email).First(&user)

		if result.Error != nil {
			CheckPasswordHash(input.Password, fakeBCryptHash)
			c.JSON(404, gin.H{"error": "Incorrect credentials"})
			return
		}

		isMatch := CheckPasswordHash(input.Password, user.Password)

		if !isMatch {
			c.JSON(401, gin.H{"error": "Incorrect credentials"})
			return
		}

		token, err := GenerateToken(user.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Your credentials are correct, but the access key generation is incorrect."})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
			"token":   token, // 把这个发给前端
		})

	})

	r.POST("/register", func(c *gin.Context) {
		var input RegisterRequest

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		superUserToken := c.GetHeader("X-Super-Token")

		if superUserToken == "" {
			c.JSON(403, gin.H{"error": "Missing token"})
			return
		}

		if string(GetSuperuserToken()) != superUserToken {
			c.JSON(403, gin.H{"error": "Invalid token"})
			return
		}

		var count int64
		// 只统计数量，不查询具体内容，速度更快
		db.Model(&User{}).Where("email = ?", input.Email).Count(&count)

		if count > 0 {
			c.JSON(400, gin.H{"error": "Email already taken"})
			return
		}

		newUser := User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password, // 记得用你之前的加密逻辑
		}

		db.Create(&newUser)

		c.JSON(201, gin.H{
			"message": "success",
		})
		return

	})

	SetUpSuperuser()

	err = r.Run()
	if err != nil {
		panic(err)
	}

}
