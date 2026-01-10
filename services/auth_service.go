package services

import (
	"errors"

	"gorm.io/gorm"
	"pdnode.com/website/models"
	"pdnode.com/website/utils"
)

type AuthService struct {
	DB *gorm.DB
}

// Login 返回用户对象、Token和错误
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	// 1. 查找用户
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		// 为了防范计时攻击，即使找不到用户也跑一遍哈希比对
		const fakeHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgNIhp.qfTMvYeJLvAbpE52EPIvG"
		utils.CheckPasswordHash(password, fakeHash)
		return nil, "", errors.New("incorrect credentials")
	}

	// 2. 比对密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("incorrect credentials")
	}

	// 3. 生成 Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, "", errors.New("token generation failed")
	}

	return &user, token, nil
}

// Register 处理注册逻辑
func (s *AuthService) Register(input models.RegisterRequest) error {
	var count int64
	s.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)
	if count > 0 {
		return errors.New("email already taken")
	}

	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	return s.DB.Create(&newUser).Error
}
