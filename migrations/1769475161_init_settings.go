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

		err := app.Save(settings)
		if err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
