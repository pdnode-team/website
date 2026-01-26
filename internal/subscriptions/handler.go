package subscriptions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v84"
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

	event, err := webhook.ConstructEventWithOptions(
		payload,
		e.Request.Header.Get("Stripe-Signature"),
		h.service.cfg.StripeSignKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true, // 忽略版本不一致报错
		},
	)
	if err != nil {
		fmt.Println(err)
		return e.BadRequestError("Invalid signature", nil)
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Println(err)
			return e.BadRequestError("JSON unmarshal failed", nil)
		}

		// 🌟 调用 Service 层处理业务（如更新用户订阅状态、发货等）
		// 传入 e.App (PocketBase 实例) 以便在 Service 里操作数据库
		if err := h.service.HandleCheckoutCompleted(session); err != nil {
			fmt.Println(err)

			return e.InternalServerError("Handle checkout failed", err)
		}
	case "invoice.paid":

		var inv stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &inv)
		if err != nil {
			fmt.Println(err)
			return e.BadRequestError("Parsing invoice failed", err)
		}

		err = h.service.HandleInvoicePaid(inv)
		if err != nil {
			fmt.Println(err)
			return e.InternalServerError("Handle checkout failed", err)
		}

	}

	return e.NoContent(http.StatusOK)
}

func (h *SubscriptionHandler) CheckSubscription(e *core.RequestEvent) error {
	subscription, err := h.service.CheckValidSubscription(e.Auth.Original())

	if errors.Is(err, sql.ErrNoRows) {

		return e.BadRequestError("No subscription", nil)
	}

	if err != nil {

		return e.InternalServerError("Check subscription failed", err)
	}

	return e.JSON(http.StatusOK, subscription)
}
