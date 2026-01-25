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

// TODO: 优化，非单文件
// TODO: 提供升级订阅的选择
// TODO: 删除或者优化日志

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
