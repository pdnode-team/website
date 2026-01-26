package subscriptions

import (
	"website-pb/config"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterRoutes(app *pocketbase.PocketBase, se *core.ServeEvent, cfg *config.Config) {
	service := NewService(app, cfg)
	handler := &SubscriptionHandler{service: service}

	// 路由注册
	se.Router.POST("/api/webhook/stripe", handler.StripeWebhook)
	se.Router.POST("/api/checkout/subscription", handler.Checkout).Bind(apis.RequireAuth())
	se.Router.GET("/api/subscription", handler.CheckSubscription).Bind(apis.RequireAuth())
}
