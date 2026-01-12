package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pdnode.com/website/models"
	"pdnode.com/website/utils" // 确保你的 utils 包路径正确
)

var authTestDB *gorm.DB

func init() {
	err := os.Chdir("..")
	if err != nil {
		panic("无法切换到根目录: " + err.Error())
	}
	// 初始化内存数据库
	authTestDB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// 迁移用户表
	_ = authTestDB.AutoMigrate(&models.User{})
}

func cleanUserTable(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

// TestAuthService_Register 包含所有注册分支
func TestAuthService_Register(t *testing.T) {
	svc := &AuthService{DB: authTestDB}

	t.Run("register_success", func(t *testing.T) {
		cleanUserTable(authTestDB)
		req := models.RegisterRequest{
			Name:     "Tester",
			Email:    "new@test.com",
			Password: "password123",
		}

		err := svc.Register(req)
		assert.NoError(t, err)

		// 验证数据库确实存进去了
		var user models.User
		err = authTestDB.Where("email = ?", "new@test.com").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "Tester", user.Name)
	})

	t.Run("register_email_taken", func(t *testing.T) {
		cleanUserTable(authTestDB)
		// 先存入一个用户
		authTestDB.Create(&models.User{Email: "taken@test.com", Name: "Exist"})

		req := models.RegisterRequest{
			Name:  "NewUser",
			Email: "taken@test.com",
		}

		err := svc.Register(req)
		assert.Error(t, err)
		assert.Equal(t, "email already taken", err.Error())
	})
}

// TestAuthService_Login 包含所有登录分支
func TestAuthService_Login(t *testing.T) {
	svc := &AuthService{DB: authTestDB}
	rawPassword := "123456"
	// 这里必须先加密，因为你的业务代码会用 CheckPasswordHash 比对
	hashedPassword, _ := utils.HashPassword(rawPassword)

	t.Run("login_success", func(t *testing.T) {
		cleanUserTable(authTestDB)
		user := models.User{
			Email:    "login@test.com",
			Password: hashedPassword,
		}
		err := authTestDB.Session(&gorm.Session{SkipHooks: true}).Create(&user).Error
		require.NoError(t, err)

		u, token, err := svc.Login("login@test.com", rawPassword)

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.NotEmpty(t, token)
		assert.Equal(t, "login@test.com", u.Email)
	})

	t.Run("login_user_not_found", func(t *testing.T) {
		cleanUserTable(authTestDB)
		// 数据库为空，直接查一个不存在的 email
		u, token, err := svc.Login("nobody@test.com", "any_pass")

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Empty(t, token)
		assert.Equal(t, "incorrect credentials", err.Error())
	})

	t.Run("login_wrong_password", func(t *testing.T) {
		cleanUserTable(authTestDB)
		// 存入正确用户
		authTestDB.Create(&models.User{Email: "wrongpass@test.com", Password: hashedPassword})

		// 尝试用错误的密码登录
		u, token, err := svc.Login("wrongpass@test.com", "wrong_password_here")

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Empty(t, token)
		assert.Equal(t, "incorrect credentials", err.Error())
	})
}
