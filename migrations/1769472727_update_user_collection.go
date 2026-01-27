package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...
		collection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		collection.Fields.Add(&core.TextField{
			Name:     "stripe_customer_id",
			Required: false,
		})

		field := collection.Fields.GetByName("name")

		if textField, ok := field.(*core.TextField); ok {
			textField.Required = true
		}

		collection.Indexes = append(collection.Indexes,
			"CREATE UNIQUE INDEX `idx_stripe_customer_id` ON `users` (`stripe_customer_id`) WHERE `stripe_customer_id` != ''",
		)
		return app.Save(collection)
	}, func(app core.App) error {
		// 回滚逻辑：移除字段和索引
		collection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return nil
		}

		collection.Fields.RemoveByName("stripe_customer_id")

		// 移除索引 (从切片中滤掉)
		newIndexes := []string{}
		for _, idx := range collection.Indexes {
			if idx != "CREATE UNIQUE INDEX `idx_stripe_customer_id` ON `users` (`stripe_customer_id`) WHERE `stripe_customer_id` != ''" {
				newIndexes = append(newIndexes, idx)
			}
		}
		collection.Indexes = newIndexes

		return app.Save(collection)
	})
}
