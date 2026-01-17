package users

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterHooks(app *pocketbase.PocketBase) {
	// 用户更新拦截
	app.OnRecordUpdateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Auth != nil && e.Auth.IsSuperuser() {
			return e.Next()
		}
		oldID := e.Record.Original().GetString("stripe_customer_id")
		newID := e.Record.GetString("stripe_customer_id")
		if newID != oldID {
			return e.BadRequestError("You do not have permission to update stripe_customer_id", nil)
		}
		return e.Next()
	})

	// 用户创建拦截
	app.OnRecordCreateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Auth != nil && e.Auth.IsSuperuser() {
			return e.Next()
		}
		if e.Record.GetString("stripe_customer_id") != "" {
			return e.BadRequestError("You do not have permission to enter stripe_customer_id", nil)
		}
		return e.Next()
	})
}
