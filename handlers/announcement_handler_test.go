package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"pdnode.com/website/models"
)

// 只需要写一次，它是所有测试的基础
type MockAnnouncementService struct{ mock.Mock }

func (m *MockAnnouncementService) GetByID(id string) (*models.Announcement, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Announcement), args.Error(1)
}
func (m *MockAnnouncementService) GetAll() ([]models.Announcement, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Announcement), args.Error(1)
}
func (m *MockAnnouncementService) Create(a *models.Announcement) error { return m.Called(a).Error(0) }
func (m *MockAnnouncementService) Update(a *models.Announcement) error { return m.Called(a).Error(0) }
func (m *MockAnnouncementService) Delete(id string) error              { return m.Called(id).Error(0) }

func testBox() (*MockAnnouncementService, *AnnouncementHandler, *httptest.ResponseRecorder, *gin.Context) {
	m := new(MockAnnouncementService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := &AnnouncementHandler{Service: m}
	return m, h, w, c // 这里返回 4 个，对应 m, h, w, c
}

func TestAnnouncementHandler_GetOne(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		m, h, w, c := testBox()

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		fakeData := &models.Announcement{ID: 1, Title: "你好"}

		m.On("GetByID", "1").Return(fakeData, nil).Once()

		h.GetOne(c)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "你好")
	})

	t.Run("not_found", func(t *testing.T) {
		m, h, w, c := testBox()
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		m.On("GetByID", "1").Return(nil, gorm.ErrRecordNotFound).Once()
		h.GetOne(c)
		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "Announcement not found")

	})

	t.Run("server_error", func(t *testing.T) {
		m, h, w, c := testBox()
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		m.On("GetByID", "1").Return(nil, errors.New("test error")).Once()

		h.GetOne(c)
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Something went wrong")

	})
}

func TestAnnouncementHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, h, w, c := testBox()
		fakeData := []models.Announcement{
			{ID: 1, Title: "公告1"},
			{ID: 2, Title: "公告2"},
		}

		m.On("GetAll").Return(fakeData, nil).Once()

		h.GetAll(c)
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "公告1")
		assert.Contains(t, w.Body.String(), "公告2")
	})
	t.Run("error", func(t *testing.T) {
		m, h, w, c := testBox()

		m.On("GetAll").Return(nil, errors.New("test error")).Once()

		h.GetAll(c)
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "test error")
	})
}

func TestAnnouncementHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, h, w, c := testBox()

		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		input := `{"title": "新公告", "content": "这是内容"}`
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		m.On("Create", mock.AnythingOfType("*models.Announcement")).Return(nil).Once()

		h.Create(c)

		assert.Equal(t, 201, w.Code)
		assert.Contains(t, w.Body.String(), "Created")
		assert.Contains(t, w.Body.String(), "00000000-0000-0000-0000-000000000000")

	})
	t.Run("parameter_error", func(t *testing.T) {
		_, h, w, c := testBox()
		input := `{""}`
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Create(c)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
	t.Run("no_userid_error", func(t *testing.T) {
		_, h, w, c := testBox()
		input := `{"title": "新公告", "content": "这是内容"}`
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Unable to obtain current user information")
	})
	t.Run("create_error", func(t *testing.T) {
		m, h, w, c := testBox()
		input := `{"title": "新公告", "content": "这是内容"}`
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		m.On("Create", mock.AnythingOfType("*models.Announcement")).Return(errors.New("test error")).Once()

		h.Create(c)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "test error")
	})
	//t.Run("error", func(t *testing.T) {
	//	m, h, w, c := testBox()
	//})
}

func TestAnnouncementHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		u, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

		m.On("GetByID", "1").Return(&models.Announcement{ID: 1, AuthorID: u}, nil).Once()

		m.On("Delete", "1").Return(nil).Once()

		h.Delete(c)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "deleted")

	})
	t.Run("not_found", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		m.On("GetByID", "1").Return(nil, errors.New("not found")).Once()

		m.On("Delete", "1").Return(nil).Once()

		h.Delete(c)

		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "not be found")

	})
	t.Run("permission_error", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		u, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")

		m.On("GetByID", "1").Return(&models.Announcement{ID: 1, AuthorID: u}, nil).Once()

		m.On("Delete", "1").Return(nil).Once()

		h.Delete(c)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "do not have permission")

	})
	t.Run("deleted_failed", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		u, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

		m.On("GetByID", "1").Return(&models.Announcement{ID: 1, AuthorID: u}, nil).Once()

		m.On("Delete", "1").Return(errors.New("test error")).Once()

		h.Delete(c)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Cannot Delete")

	})
}

func TestAnnouncementHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		input := `{"title": "你好", "content": "这是内容"}`
		c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		m.On("GetByID", "1").Return(&models.Announcement{ID: 1}, nil).Once()

		m.On("Update", mock.AnythingOfType("*models.Announcement")).Return(nil).Once()

		h.Update(c)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "你好")
	})
	t.Run("not_found", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		input := `{"title": "你好", "content": "这是内容"}`
		c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		m.On("GetByID", "1").Return(nil, errors.New("not found")).Once()

		m.On("Update", mock.AnythingOfType("*models.Announcement")).Return(nil).Once()

		h.Update(c)

		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "not be found")
	})
	t.Run("permission_error", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		u, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		m.On("GetByID", "1").Return(&models.Announcement{ID: 1, AuthorID: u}, nil).Once()

		h.Update(c)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "You do not have permission")

	})
	t.Run("parameter_error", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		input := `{"`
		c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		m.On("GetByID", "1").Return(&models.Announcement{ID: 1, AuthorID: uid}, nil).Once()
		m.On("Update", mock.AnythingOfType("*models.Announcement")).Return(nil).Once()

		h.Update(c)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "EOF")
		assert.Contains(t, w.Body.String(), "error")

	})
	t.Run("update_failed", func(t *testing.T) {
		m, h, w, c := testBox()
		uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		c.Set("currentUserID", uid)

		input := `{"title": "新公告", "content": "这是内容"}`

		c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		m.On("GetByID", "1").Return(&models.Announcement{ID: 1, AuthorID: uid}, nil).Once()
		m.On("Update", mock.AnythingOfType("*models.Announcement")).Return(errors.New("test error")).Once()

		h.Update(c)

		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), "Update Failed")

	})
}
