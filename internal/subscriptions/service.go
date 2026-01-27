package subscriptions

import (
	"database/sql"
	"errors"
	"time"
	"website-pb/config"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

var (
	ErrAlreadySubscribed = errors.New("already_have_subscription")
	ErrPlanInvalid       = errors.New("invalid_plan")
)

type SubscriptionService struct {
	app *pocketbase.PocketBase
	cfg *config.Config
}

func NewService(app *pocketbase.PocketBase, cfg *config.Config) *SubscriptionService {
	return &SubscriptionService{app: app, cfg: cfg}
}

func (s *SubscriptionService) CheckValidSubscription(user *core.Record) (*core.Record, error) {

	now := time.Now().UTC().Format("2006-01-02 15:04:05.000Z")

	record := &core.Record{}

	err := s.app.RecordQuery("subscriptions").
		// 1. 基础过滤：用户 ID
		AndWhere(dbx.HashExp{"user_id": user.Id}).
		// 2. 时间过滤：未过期
		AndWhere(dbx.NewExp("expires_at > {:now}", dbx.Params{"now": now})).
		// 3. 排序：将到期时间最远的排在最前面
		OrderBy("expires_at DESC").
		// 4. 只取一条
		Limit(1).
		// 5. 将结果映射到 record 对象
		One(record)

	if err != nil {
		return nil, err // 没找到会返回 sql.ErrNoRows
	}

	return record, nil
}

// CreateCheckoutSession 处理 Stripe 会话创建
func (s *SubscriptionService) CreateCheckoutSession(user *core.Record, plan string, frontendURL string) (string, error) {
	priceID, exists := s.cfg.PlanToPrice[plan]
	if !exists || priceID == "" {
		return "", ErrPlanInvalid
	}

	sub, err := s.CheckValidSubscription(user.Original())

	// 情况 A: 找到了有效订阅 (err == nil)
	if err == nil && sub != nil {
		return "", ErrAlreadySubscribed
	}

	// 情况 B: 发生了真正的数据库错误 (不是“没找到”)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// 结构化日志修复：使用 key-value
		s.app.Logger().Error("Failed to check subscription", "error", err.Error(), "userId", user.Id)
		return "", errors.New("check subscription failed")
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

func (s *SubscriptionService) HandleInvoicePaid(inv stripe.Invoice) error {

	if inv.Customer == nil {
		s.app.Logger().Warn("No customer found: ", inv)
		return errors.New("invoice customer is nil")
	}

	if len(inv.Lines.Data) == 0 {
		s.app.Logger().Warn("invoice has no lines: ", inv)

		return errors.New("invoice has no lines")
	}

	stripeCustomerID := inv.Customer.ID
	user, err := s.app.FindFirstRecordByFilter("users", "stripe_customer_id = {:id}", map[string]any{"id": stripeCustomerID})
	if err != nil {
		return err
	}

	collection, err := s.app.FindCollectionByNameOrId("subscriptions")
	if err != nil {
		s.app.Logger().Error("Failed to find collection: ", err)
		return errors.New("subscriptions collection not found")
	}

	record := core.NewRecord(collection)

	priceID := inv.Lines.Data[0].Pricing.PriceDetails.Price

	priceIDMap := s.cfg.PriceToPlan[priceID]

	if priceIDMap == "" {
		s.app.Logger().Warn("No prices found for invoice ", inv)
		return errors.New("invalid price")
	}

	expiresAt := time.Unix(inv.Lines.Data[0].Period.End, 0).UTC()

	record.Set("user_id", user.Id)
	record.Set("plan", priceIDMap)
	record.Set("stripe_invoice_id", inv.ID)
	record.Set("expires_at", expiresAt)

	return s.app.Save(record)
}

func (s *SubscriptionService) HandleCheckoutCompleted(sess stripe.CheckoutSession) error {
	user, err := s.app.FindRecordById("users", sess.ClientReferenceID)
	if err != nil {
		return err
	}

	// 2. 更新用户信息
	user.Set("stripe_customer_id", sess.Customer.ID)

	return s.app.Save(user)
}
