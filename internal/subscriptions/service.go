package subscriptions

import (
	"fmt"
	"time"
	"website-pb/config"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/subscription"
)

type SubscriptionService struct {
	app *pocketbase.PocketBase
	cfg *config.Config
}

func NewService(app *pocketbase.PocketBase, cfg *config.Config) *SubscriptionService {
	return &SubscriptionService{app: app, cfg: cfg}
}

// CreateCheckoutSession 处理 Stripe 会话创建
func (s *SubscriptionService) CreateCheckoutSession(user *core.Record, plan string, frontendURL string) (string, error) {
	priceID, exists := s.cfg.PlanToPrice[plan]
	if !exists || priceID == "" {
		return "", fmt.Errorf("invalid plan: %s", plan)
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL:        stripe.String(frontendURL + "/checkout/success?id={CHECKOUT_SESSION_ID}"),
		CancelURL:         stripe.String(frontendURL),
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		ClientReferenceID: stripe.String(user.Id),
		Metadata:          map[string]string{"plan": plan},
		LineItems:         []*stripe.CheckoutSessionLineItemParams{{Price: stripe.String(priceID), Quantity: stripe.Int64(1)}},
	}

	// 关联已有 Stripe Customer ID
	if cusID := user.GetString("stripe_customer_id"); cusID != "" {
		params.Customer = stripe.String(cusID)
	} else {
		params.CustomerEmail = stripe.String(user.Email())
	}

	sess, err := session.New(params)
	return sess.URL, err
}

// HandleCheckoutCompleted 处理支付成功后的数据更新
func (s *SubscriptionService) HandleCheckoutCompleted(sess stripe.CheckoutSession) error {
	// 1. 获取用户
	user, err := s.app.FindRecordById("users", sess.ClientReferenceID)
	if err != nil {
		return err
	}

	// 2. 更新用户信息
	user.Set("stripe_customer_id", sess.Customer.ID)
	user.Set("plan", sess.Metadata["plan"])
	if err := s.app.Save(user); err != nil {
		return err
	}

	// 3. 创建订阅记录
	sub, _ := subscription.Get(sess.Subscription.ID, nil)
	collection, _ := s.app.FindCollectionByNameOrId("subscriptions")
	record := core.NewRecord(collection)
	record.Set("user_id", user.Id)
	record.Set("plan", sess.Metadata["plan"])
	record.Set("expires_at", time.Unix(sub.Items.Data[0].CurrentPeriodEnd, 0).UTC())

	return s.app.Save(record)
}
