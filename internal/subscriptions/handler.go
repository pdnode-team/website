package subscriptions

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v84/webhook"
)

type SubscriptionHandler struct {
	service *SubscriptionService
}

func (h *SubscriptionHandler) Checkout(e *core.RequestEvent) error {
	user := e.Auth
	if user == nil {
		return apis.NewUnauthorizedError("Login first", nil)
	}

	var data struct {
		Plan string `json:"plan"`
	}
	if err := e.BindBody(&data); err != nil {
		return err
	}

	// 1. 获取动态的基础地址 (例如 https://example.com 或 http://localhost:8090)
	scheme := "http"
	if e.Request.TLS != nil || e.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := e.Request.Host // 这会自动获取当前访问的域名和端口

	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	url, err := h.service.CreateCheckoutSession(user, data.Plan, baseURL)
	if err != nil {
		return e.BadRequestError(err.Error(), nil)
	}

	return e.JSON(http.StatusOK, map[string]string{"url": url})
}

func (h *SubscriptionHandler) StripeWebhook(e *core.RequestEvent) error {
	payload, err := io.ReadAll(e.Request.Body)
	if err != nil {
		return e.BadRequestError("Read body failed", nil)
	}

	event, err := webhook.ConstructEvent(payload, e.Request.Header.Get("Stripe-Signature"), h.service.cfg.StripeSignKey)
	if err != nil {
		return e.BadRequestError("Invalid signature", nil)
	}

	if event.Type == "checkout.session.completed" {
		// ... 反序列化并调用 h.service.HandleCheckoutCompleted(sess)
	}

	return e.NoContent(http.StatusOK)
}
