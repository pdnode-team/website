package main

import (
	"os"
	"strings"
	"website-pb/config"
	"website-pb/internal/subscriptions"
	"website-pb/internal/users"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/stripe/stripe-go/v84"
)

// TODO: 禁止重复订阅
// TODO: 防止 /checkout/success 被滥用
// TODO: 日志
// TODO: 创建订阅集合迁移文件
// TODO: 创建用户集合迁移文件
// TODO: 用.env初始化SMTP和设置
// TODO: 添加限速
// TODO: 发送各种邮件

const version string = "v1.0.0-alpha"

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

	// 版本
	app.Logger().Info("Pdnode Website API " + version)

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	}) // 4. 注册路由

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.GET("/{path...}", apis.Static(os.DirFS("./web/build"), true))

		// 调用订阅模块，把 app, se 和 cfg 传进去
		subscriptions.RegisterRoutes(app, se, cfg)
		return se.Next()
	})

	err := app.Start()
	if err != nil {
		panic(err)
	}
}
