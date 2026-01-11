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

//	type AnnouncementHandler struct {
//		Service *services.AnnouncementService
//	}
type AnnouncementHandler struct {
	Service services.AnnouncementServiceInterface
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to obtain current user information"})
		return
	}

	announcement := models.Announcement{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: val.(uuid.UUID),
	}

	if err := h.Service.Create(&announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created",
		"data":    announcement,
	})

}

func (h *AnnouncementHandler) Update(c *gin.Context) {

	val := c.MustGet("currentUserID").(uuid.UUID)

	id := c.Param("id")

	announcement, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The announcement could not be found."})
		return
	}

	if announcement.AuthorID != val {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to modify other people's announcements."})
		return
	}

	var input models.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcement.Title = input.Title
	announcement.Content = input.Content

	if err := h.Service.Update(announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update Failed"})
		return
	}

	c.JSON(http.StatusOK, announcement)

}

func (h *AnnouncementHandler) Delete(c *gin.Context) {
	val := c.MustGet("currentUserID").(uuid.UUID)
	id := c.Param("id")

	announcement, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The announcement could not be found."})
		return
	}

	if announcement.AuthorID != val {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to modify other people's announcements."})
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})

	return

}
