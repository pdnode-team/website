package handlers

import (
	"errors"
	"log/slog"
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

		slog.Error("Database query error", "handler", "GetOne", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *AnnouncementHandler) GetAll(c *gin.Context) {
	announcements, err := h.Service.GetAll()
	if err != nil {
		slog.Error("Failed to fetch announcements", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, announcements)
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	var input models.CreateAnnouncementRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		slog.Debug("Create announcement bind failed", "err", err) // DEBUG 即可，不干扰生产
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	val, exists := c.Get("currentUserID")

	if !exists {
		slog.Error("Create announcement obtain current user information failed", "err", "Unable to obtain current user information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to obtain current user information"})
		return
	}

	announcement := models.Announcement{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: val.(uuid.UUID),
	}

	if err := h.Service.Create(&announcement); err != nil {
		slog.Error("Create announcement failed", "err", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 在 c.JSON 成功返回之前
	slog.Info("Announcement created",
		"id", announcement.ID,
		"author_id", announcement.AuthorID,
		"title", announcement.Title)

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
		slog.Warn("Unauthorized update attempt",
			"announcement_id", id,
			"user_id", val)
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
		slog.Error("Announcement failed", "id", id, "operator_id", val, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update Failed"})
		return
	}

	slog.Info("Announcement updated", "id", id, "operator_id", val)

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
		slog.Warn("Unauthorized delete attempt", "id", id, "user_id", val)
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to modify other people's announcements."})
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		slog.Error("Delete Announcement Failed", "id", id, "user_id", val, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Delete"})
		return
	}

	slog.Info("Announcement deleted", "id", id, "operator_id", val)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})

	return

}
