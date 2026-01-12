package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pdnode.com/website/models"
)

// 定义 Mock 对象
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string) (*models.User, string, error) {
	args := m.Called(email, password)
	// 处理 nil 返回值的情况，避免类型断言 Panic
	var u *models.User
	if args.Get(0) != nil {
		u = args.Get(0).(*models.User)
	}
	return u, args.String(1), args.Error(2)
}

func (m *MockAuthService) Register(input models.RegisterRequest) error {
	return m.Called(input).Error(0)
}

// 刷新 testBox
func authTestBox() (*MockAuthService, *AuthHandler, *httptest.ResponseRecorder, *gin.Context) {
	m := new(MockAuthService)
	h := &AuthHandler{Service: m} // 注入 Mock
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return m, h, w, c
}

func TestAuthHandler_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, h, w, c := authTestBox()

		input := `{"name":"Tester", "email":"new@test.com", "password":"password123"}`
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		// ✅ 剧本：直接返回 nil，不会再去读任何文件！
		m.On("Register", mock.Anything).Return(nil).Once()

		h.Register(c)
		assert.Equal(t, 201, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		m, h, w, c := authTestBox()

		input := `{"name":"Tester", "email":"taken@test.com", "password":"password123"}`
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		m.On("Register", mock.Anything).Return(errors.New("email already taken")).Once()

		h.Register(c)
		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "email already taken")
	})
	t.Run("register_bind_error", func(t *testing.T) {
		_, h, w, c := authTestBox()
		c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBufferString(`{invalid}`))
		c.Request.Header.Set("Content-Type", "application/json")
		h.Register(c)
		assert.Equal(t, 400, w.Code)
	})
}
func TestAuthHandler_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, h, w, c := authTestBox()

		// 1. 准备输入
		input := `{"email":"test@test.com", "password":"password123"}`
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		// 2. 模拟 Service 返回成功
		fakeUser := &models.User{Email: "test@test.com"}
		m.On("Login", "test@test.com", "password123").Return(fakeUser, "fake-token", nil).Once()

		// 3. 执行
		h.Login(c)

		// 4. 断言
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "fake-token")
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("bind_error", func(t *testing.T) {
		m, h, w, c := authTestBox()

		// 故意传入非法 JSON 结构，覆盖第一个 if err != nil 分支
		input := `{"email": "test@test.com", "password": }` // 错误的 JSON 格式
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		// 验证没有调用 Service
		m.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
	})

	t.Run("unauthorized_error", func(t *testing.T) {
		m, h, w, c := authTestBox()

		// 传入正确格式，但模拟 Service 认证失败，覆盖第二个 if err != nil 分支
		input := `{"email":"wrong@test.com", "password":"wrong_password"}`
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(input))
		c.Request.Header.Set("Content-Type", "application/json")

		m.On("Login", "wrong@test.com", "wrong_password").
			Return(nil, "", errors.New("incorrect credentials")).Once()

		h.Login(c)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "incorrect credentials")
	})
}
