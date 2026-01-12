package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require" // 💡 推荐用 require：如果失败直接中断测试
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pdnode.com/website/models"
)

// 全局测试变量，减少重复初始化
var testDB *gorm.DB

func init() {
	// 在测试包加载时初始化一次
	testDB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = testDB.AutoMigrate(&models.Announcement{})
}

// 每次测试前清理表，保证数据隔离
func cleanData(db *gorm.DB) {
	db.Exec("DELETE FROM announcements")
}

func TestAnnouncementService_AllInOne(t *testing.T) {
	svc := &AnnouncementService{DB: testDB}

	t.Run("Create & Get", func(t *testing.T) {
		cleanData(testDB)
		a := &models.Announcement{Title: "Service"}

		err := svc.Create(a)
		require.NoError(t, err) // require 失败会终止当前 Run，防止后面访问 nil 指针 panic

		found, err := svc.GetByID(fmt.Sprintf("%d", a.ID))
		assert.NoError(t, err)
		assert.Equal(t, "Service", found.Title)
	})

	t.Run("Delete", func(t *testing.T) {
		cleanData(testDB)
		// 1. 快速造数据
		a := models.Announcement{Title: "待删除"}
		testDB.Create(&a)

		// 2. 执行删除
		err := svc.Delete(fmt.Sprintf("%d", a.ID))
		assert.NoError(t, err)

		// 3. 验证 (更优雅的写法)
		var count int64
		testDB.Model(&models.Announcement{}).Where("id = ?", a.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
	t.Run("GetAll_Success_And_Order", func(t *testing.T) {
		cleanData(testDB)
		// 1. 构造多条数据，验证排序 (desc)
		testDB.Create(&models.Announcement{Title: "旧公告", Content: "1"})
		testDB.Create(&models.Announcement{Title: "新公告", Content: "2"})

		announcements, err := svc.GetAll()

		assert.NoError(t, err)
		assert.Len(t, announcements, 2)
		// 验证排序：第一条应该是最后创建的那条
		assert.Equal(t, "新公告", announcements[0].Title)
	})

	t.Run("GetAll_DB_Error", func(t *testing.T) {
		// 1. 将表名改掉，让业务代码找不到表
		err := testDB.Migrator().RenameTable(&models.Announcement{}, "temp_announcements")
		if err != nil {
			return
		}

		// 2. 确保测试结束后把名字改回来
		defer func(migrator gorm.Migrator, oldName, newName interface{}) {
			err := migrator.RenameTable(oldName, newName)
			if err != nil {

			}
		}(testDB.Migrator(), "temp_announcements", &models.Announcement{})

		announcements, err := svc.GetAll()

		// 3. 此时 Find 会因为找不到表而报错
		assert.Error(t, err)
		assert.Nil(t, announcements)
	})
	t.Run("Create_DB_Error", func(t *testing.T) {
		cleanData(testDB)

		// 故意手动创建一个重复的 ID，触发主键冲突错误
		a1 := &models.Announcement{ID: 1, Title: "First"}
		testDB.Create(a1)

		a2 := &models.Announcement{ID: 1, Title: "Second"} // 相同的 ID: 1
		err := svc.Create(a2)

		assert.Error(t, err) // 这里会覆盖 if err != nil 分支
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		cleanData(testDB)
		// 尝试查询一个不存在的 ID (999)
		found, err := svc.GetByID("999")

		assert.Error(t, err)
		assert.Nil(t, found)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Update_Success", func(t *testing.T) {
		cleanData(testDB)
		// 1. 先创建原始数据
		a := models.Announcement{Title: "原标题", Content: "原内容"}
		testDB.Create(&a)

		// 2. 构造更新对象
		updateInfo := &models.Announcement{
			ID:      a.ID,
			Title:   "修改后的标题",
			Content: "修改后的内容",
		}

		// 3. 执行更新
		err := svc.Update(updateInfo)
		assert.NoError(t, err)

		// 4. 从数据库重新读取验证
		var updated models.Announcement
		testDB.First(&updated, a.ID)
		assert.Equal(t, "修改后的标题", updated.Title)
		assert.Equal(t, "修改后的内容", updated.Content)
	})

	t.Run("NewAnnouncementService_Factory", func(t *testing.T) {
		// 测试工厂函数是否正确返回接口
		factorySvc := NewAnnouncementService(testDB)
		assert.NotNil(t, factorySvc)

		// 验证它确实是 *AnnouncementService 类型
		_, ok := factorySvc.(*AnnouncementService)
		assert.True(t, ok)
	})
}
