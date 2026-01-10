package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"pdnode.com/website/models"
	"pdnode.com/website/services"
)

type AnnouncementHandler struct {
	Service *services.AnnouncementService
}

func (h *AnnouncementHandler) GetOne(c *gin.Context) {
	id := c.Param("id")

	a, err := h.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Announcement not found"})
			return
		}
		// 这里的打印可以保留，方便服务器排查
		println("[Server Error] Query Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *AnnouncementHandler) GetAll(c *gin.Context) {
	announcements, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, announcements)
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	var input models.CreateAnnouncementRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	val, exists := c.Get("currentUserID")

	if !exists {
		c.JSON(500, gin.H{"error": "Unable to obtain current user information"})
		return
	}

	announcement := models.Announcement{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: val.(uuid.UUID),
	}

	if err := h.Service.Create(&announcement); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Created",
		"data":    announcement,
	})

}
