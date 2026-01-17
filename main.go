package main

import (
	"website-pb/config"
	"website-pb/internal/subscriptions"
	"website-pb/internal/users"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v84"
)

//
//const (
//	PlanStarter = "starter"
//	PlanPro     = "pro"
//	PlanPlus    = "plus"
//)
//
//// 2. 将名称映射到具体的 Stripe ID
//var planToPrice = map[string]string{
//	PlanStarter: "price_1SpKLiCwYVbOisgIfZmtAoWG",
//	PlanPro:     "price_1SpKM5CwYVbOisgIm7Mf0Gn3",
//	PlanPlus:    "price_1SpKLUCwYVbOisgIaGOI6AZa",
//}
//
//// TODO: 优化，非单文件
//// TODO: 提供升级订阅的选择
//// TODO: 删除或者优化日志
//
//type Config struct {
//	FrontendURL   string
//	StripeKey     string
//	StripeSignKey string
//}
//
//func NewConfig() *Config {
//	cfg := &Config{}
//
//	// 设置默认值或从环境变量读取
//	cfg.FrontendURL = getEnv("FRONTEND_URL", "http://localhost:5173")
//	cfg.StripeKey = getEnv("STRIPE_KEY", "")
//	cfg.StripeSignKey = getEnv("STRIPE_SIGN_KEY", "")
//
//	return cfg
//}
//
//func getEnv(key, defaultValue string) string {
//	if value, exists := os.LookupEnv(key); exists {
//		return value
//	}
//	if defaultValue == "" {
//		panic("Environment variable " + key + " is empty and required.")
//	}
//	return defaultValue
//}
//
//func main() {
//	app := pocketbase.New()
//
//	err := godotenv.Load()
//	if err != nil {
//		app.Logger().Warn("No .env file found, relying on system environment variables")
//	}
//
//	config := NewConfig()
//
//	stripe.Key = config.StripeKey
//
//	app.OnRecordUpdateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
//		// 1. 获取数据库中的原始值 (Original() 返回 *models.Record)
//		// 使用 GetString("field") 代替 ["field"]
//
//		if e.Auth != nil && e.Auth.IsSuperuser() {
//			return e.Next()
//		}
//
//		oldID := e.Record.Original().GetString("stripe_customer_id")
//
//		// 2. 获取用户在当前请求中试图写入的新值
//		newID := e.Record.GetString("stripe_customer_id")
//
//		// 3. 对比并强制拦截
//		if newID != oldID {
//			// 如果新旧不一致
//			return e.BadRequestError("You do not have permission to update stripe_customer_id", nil)
//		}
//
//		return e.Next()
//	})
//
//	app.OnRecordCreateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
//		// e.App
//		// e.Collection
//		// e.Record
//		// and all RequestEvent fields...
//		if e.Auth != nil && e.Auth.IsSuperuser() {
//			return e.Next()
//		}
//
//		stripeCustomerId := e.Record.GetString("stripe_customer_id")
//
//		if stripeCustomerId != "" {
//
//			return e.BadRequestError("You do not have permission to enter stripe_customer_id", nil)
//		}
//
//		return e.Next()
//	})
//
//	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
//		// serves static files from the provided public dir (if exists)
//
//		//se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))
//		se.Router.POST("/api/checkout/subscription", func(e *core.RequestEvent) error {
//			authRecord := e.Auth
//
//			if authRecord == nil {
//				return apis.NewUnauthorizedError("Login first", nil)
//			}
//
//			data := struct {
//				DataPlan string `json:"plan"`
//			}{}
//			if err := e.BindBody(&data); err != nil {
//				return err
//			}
//
//			priceID, exists := planToPrice[data.DataPlan]
//			if !exists {
//				return e.BadRequestError("Invalid plan selection", nil)
//			}
//
//			// 获取当前 UTC 时间字符串
//			now := time.Now().UTC().Format("2006-01-02 15:04:05.000Z")
//
//			// 语法：集合名_via_关联字段名
//			// 查找：存在一个订阅记录，其 expires_at 大于现在
//
//			record, err := app.FindFirstRecordByFilter(
//				"subscriptions",
//				"user_id = {:userId} && expires_at > {:now}",
//				dbx.Params{
//					"userId": authRecord.Id,
//					"now":    now,
//				},
//			)
//
//			if err != nil {
//				// 关键点：如果是“没找到记录”，这说明用户目前没有有效订阅，是正常的！
//				if !errors.Is(err, sql.ErrNoRows) {
//					return e.InternalServerError("Database lookup failed", nil)
//				}
//			}
//			if record != nil {
//				return e.BadRequestError("Do not subscribe repeatedly.", nil)
//			}
//
//			params := &stripe.CheckoutSessionParams{
//				SuccessURL: stripe.String(config.FrontendURL + "/checkout/success?session_id={CHECKOUT_SESSION_ID}"),
//				CancelURL:  stripe.String(config.FrontendURL + ""),
//				Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
//				PaymentMethodTypes: stripe.StringSlice([]string{
//					"card",
//				}),
//				LineItems: []*stripe.CheckoutSessionLineItemParams{
//					{
//						Price:    stripe.String(priceID), // 在 Stripe 后台创建的 Price ID
//						Quantity: stripe.Int64(1),
//					},
//				},
//				Metadata: map[string]string{
//					"plan": data.DataPlan, // 存入 "pro"
//				},
//
//				//CustomerEmail:     stripe.String(authRecord.Email()),
//				ClientReferenceID: stripe.String(authRecord.Id),
//			}
//			if authRecord.GetString("stripe_customer_id") != "" {
//				params.Customer = stripe.String(authRecord.GetString("stripe_customer_id"))
//			} else {
//				params.CustomerEmail = stripe.String(authRecord.Email())
//			}
//
//			s, err := session.New(params)
//
//			if err != nil {
//				app.Logger().Error("Stripe session creation failed",
//					"error", err,
//					"userId", authRecord.Id,
//				)
//				return apis.NewInternalServerError("Could not create session", err)
//			}
//
//			return e.JSON(http.StatusOK, map[string]string{
//				"url":    s.URL,
//				"status": "success",
//			})
//
//		}).Bind(apis.RequireAuth())
//		se.Router.POST("/api/webhook/stripe", func(e *core.RequestEvent) error {
//			const MaxBodyBytes = int64(65536)
//			payload, err := io.ReadAll(io.LimitReader(e.Request.Body, MaxBodyBytes))
//			if err != nil {
//				return e.BadRequestError("Read body failed", nil)
//			}
//
//			signature := e.Request.Header.Get("Stripe-Signature")
//			endpointSecret := config.StripeSignKey
//
//			event, err := webhook.ConstructEventWithOptions(payload, signature, endpointSecret, webhook.ConstructEventOptions{
//				IgnoreAPIVersionMismatch: true,
//			})
//			if err != nil {
//				app.Logger().Warn("Stripe webhook signature verification failed", "error", err)
//				return e.BadRequestError("Invalid signature", nil)
//			}
//
//			// --- 修改点：根据事件类型进入不同逻辑 ---
//			switch event.Type {
//			case "checkout.session.completed":
//				var checkoutSession stripe.CheckoutSession
//				if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
//					return e.InternalServerError("Could not unmarshal checkout session", err)
//				}
//
//				pbUserId := checkoutSession.ClientReferenceID
//				stripeCusId := checkoutSession.Customer.ID
//				userPlan := checkoutSession.Metadata["plan"]
//
//				// 只有在这个事件里才检查 metadata
//				if userPlan == "" {
//					app.Logger().Warn("[Stripe][Subscription] No plan info in metadata", "sessionId", checkoutSession.ID)
//					return e.InternalServerError("No plan info in metadata", nil)
//				}
//
//				user, err := app.FindRecordById("users", pbUserId)
//				if err != nil {
//					app.Logger().Warn("[Subscription] User not found", "user_id", pbUserId, "error", err)
//
//					return e.BadRequestError("User not found", nil)
//				}
//
//				if user.Get("stripe_customer_id") == "" {
//					user.Set("stripe_customer_id", stripeCusId)
//				}
//				user.Set("plan", userPlan)
//
//				if err := app.Save(user); err != nil {
//					app.Logger().Error("Failed to update user plan after checkout",
//						"error", err,
//						"userId", pbUserId,
//					)
//					return e.InternalServerError("Could not save user", nil)
//				}
//
//				collection, err := app.FindCollectionByNameOrId("subscriptions")
//
//				if err != nil {
//					app.Logger().Error("[Subscription] Failed to find collection by id", "error", err)
//
//					return e.InternalServerError("Could not find collection", nil)
//				}
//
//				record := core.NewRecord(collection)
//
//				record.Set("user_id", pbUserId)
//				record.Set("plan", userPlan)
//				record.Set("stripe_session_id", checkoutSession.ID)
//				subID := checkoutSession.Subscription.ID
//				sub, err := subscription.Get(subID, nil)
//
//				if err != nil {
//					app.Logger().Error("[Subscription] Failed to get subscription", "error", err)
//					return e.InternalServerError("Could not get subscription", nil)
//				}
//
//				record.Set("expires_at", time.Unix(sub.Items.Data[0].CurrentPeriodEnd, 0).UTC())
//
//				err = app.Save(record)
//				if err != nil {
//					app.Logger().Error("[Subscription] Failed to save record", "error", err)
//					return e.InternalServerError("Could not create record", nil)
//				}
//
//				//record.Set("expires_at", time.Unix(sub.CurrentPeriodEnd, 0).UTC())
//				app.Logger().Info("Subscription activated",
//					"userId", pbUserId,
//					"plan", userPlan,
//				)
//
//			default:
//				// 对于其他所有事件 (charge.succeeded, customer.created 等)
//				// 统统打印并返回 200 OK
//				app.Logger().Debug("[Stripe][Webhook] This event has been ignored.",
//					"event", event.Type,
//				)
//			}
//
//			// 统一返回 200
//			return e.NoContent(http.StatusOK)
//		})
//
//		return se.Next()
//	})
//
//	if err := app.Start(); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	app := pocketbase.New()

	if err := godotenv.Load(); err != nil {
		app.Logger().Warn("No .env file found, using system env")
	}

	cfg := config.New()

	// 2. 设置全局 Stripe Key
	stripe.Key = cfg.StripeKey

	// 3. 注册其他模块的钩子
	users.RegisterHooks(app)

	// 4. 注册路由
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 调用订阅模块，把 app, se 和 cfg 传进去
		subscriptions.RegisterRoutes(app, se, cfg)
		return se.Next()
	})

	err := app.Start()
	if err != nil {
		panic(err)
	}
}
