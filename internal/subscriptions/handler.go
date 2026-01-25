package subscriptions

import (
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

	url, err := h.service.CreateCheckoutSession(user, data.Plan)
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
