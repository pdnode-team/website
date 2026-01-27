package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...
		collection := core.NewBaseCollection("subscriptions")

		collection.Fields.Add(&core.RelationField{
			Name:          "user_id",
			CollectionId:  "_pb_users_auth_", // 目标集合 ID 或名称
			CascadeDelete: false,
			MaxSelect:     1,
		})

		collection.Fields.Add(&core.SelectField{
			Name:      "plan",
			Values:    []string{"starter", "plus", "pro"},
			MaxSelect: 1,
			Required:  true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "stripe_invoice_id",
			Required: true,
		})

		collection.Fields.Add(&core.DateField{
			Name:     "expires_at",
			Required: true,
		})

		// 创建时自动设为当前时间
		collection.Fields.Add(&core.AutodateField{
			Name:     "created",
			OnCreate: true,
		})

		// 创建和更新时都自动同步时间
		collection.Fields.Add(&core.AutodateField{
			Name:     "updated",
			OnCreate: true,
			OnUpdate: true,
		})

		collection.Indexes = append(collection.Indexes,
			"CREATE UNIQUE INDEX `idx_stripe_inv` ON `subscriptions` (`stripe_invoice_id`)",
		)

		return app.Save(collection)
	}, func(app core.App) error {
		// add down queries...
		collection, err := app.FindCollectionByNameOrId("subscriptions")
		if err != nil {
			return nil // 已经删了就不报错
		}
		return app.Delete(collection)
	})
}
