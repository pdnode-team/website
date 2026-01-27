package migrations

import (
	"website-pb/config"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...

		settings := app.Settings()

		config.InitRateLimitRule(settings)

		settings.Logs.MaxDays = 90

		return app.Save(settings)
	}, func(app core.App) error {
		// add down queries...
		println("[Init Settings] Cannot Down the migration")
		return nil
	})
}
