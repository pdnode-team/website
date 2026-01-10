package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"pdnode.com/website/utils"
)

type Announcement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAnnouncementRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Name     string    `json:"name" binding:"required,min=2,max=12"`
	Email    string    `json:"email" binding:"required,email" gorm:"unique;not null"`
	Password string    `json:"-" binding:"required,min=8,max=32"`
}

func (u *User) BeforeCreate(*gorm.DB) (err error) {
	u.ID = uuid.New()

	// 假设你有一个全局或从配置文件读取 Token 的函数 GetSuperuserToken()
	// 这里的 u.Password 此时还是前端传来的明文
	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err // 如果加密失败，中止创建
	}

	u.Password = hashed // 将明文替换为密文
	return nil
}
func (u *User) BeforeUpdate(*gorm.DB) (err error) {
	// 检查 Password 是否已经是 Bcrypt 哈希格式
	// 如果不是以 $2a$ 或 $2b$ 开头，说明是用户传来的新明文
	if !strings.HasPrefix(u.Password, "$2a$") && !strings.HasPrefix(u.Password, "$2b$") {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=12"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}
