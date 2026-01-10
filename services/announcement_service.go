package services

import (
	"gorm.io/gorm"
	"pdnode.com/website/models"
)

type AnnouncementService struct {
	DB *gorm.DB
}

// GetByID 根据 ID 获取单个公告
func (s *AnnouncementService) GetByID(id string) (*models.Announcement, error) {
	var a models.Announcement
	if err := s.DB.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
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
