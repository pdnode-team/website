package services

import (
	"gorm.io/gorm"
	"pdnode.com/website/models"
)

//type AnnouncementService struct {
//	DB *gorm.DB
//}

type AnnouncementServiceInterface interface {
	GetByID(id string) (*models.Announcement, error)
	GetAll() ([]models.Announcement, error)
	Create(announcement *models.Announcement) error
	Delete(id string) error
	Update(announcement *models.Announcement) error
}

// AnnouncementService 定义具体的结构体（实现者）
type AnnouncementService struct {
	DB *gorm.DB
}

// NewAnnouncementService 3. 构造函数：返回接口类型
func NewAnnouncementService(db *gorm.DB) AnnouncementServiceInterface {
	return &AnnouncementService{DB: db}
}

// GetByID 根据 ID 获取单个公告
func (s *AnnouncementService) GetByID(id string) (*models.Announcement, error) {
	var a models.Announcement
	if err := s.DB.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
	//return nil, nil
}

// GetAll 获取所有公告（按时间倒序）
func (s *AnnouncementService) GetAll() ([]models.Announcement, error) {
	var announcements []models.Announcement
	if err := s.DB.Order("created_at desc").Find(&announcements).Error; err != nil {
		return nil, err
	}
	return announcements, nil
}

func (s *AnnouncementService) Create(announcement *models.Announcement) error {
	if err := s.DB.Create(announcement).Error; err != nil {
		return err
	}
	return nil

}

func (s *AnnouncementService) Delete(id string) error {
	return s.DB.Delete(&models.Announcement{}, "id = ?", id).Error
}

func (s *AnnouncementService) Update(announcement *models.Announcement) error {
	return s.DB.Model(announcement).Updates(models.Announcement{
		Title:   announcement.Title,
		Content: announcement.Content,
	}).Error
}
