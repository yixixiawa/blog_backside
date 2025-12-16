package AutoMigrate

import (
	"fmt"
	"log"
	"sqlite_test/Model"
	"sqlite_test/database"
)

func Generation_sql() {
	fmt.Println("=== 开始数据库迁移 ===")

	// 测试数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	fmt.Println("数据库连接正常")

	// 所有模型列表
	models := []struct {
		model interface{}
		name  string
	}{
		{model: &Model.User{}, name: "User"},
		{model: &Model.Tag{}, name: "Tag"},
		{model: &Model.Content{}, name: "Content"},
		{model: &Model.FileRecord{}, name: "FileRecord"},
		{model: &Model.ContentTag{}, name: "ContentTag"},
		{model: &Model.ContentFile{}, name: "ContentFile"},
		{model: &Model.Comment{}, name: "Comment"},
		{model: &Model.EmailVerify{}, name: "EmailVerify"},
	}

	successCount := 0
	for _, m := range models {
		fmt.Printf("\n--- 迁移 %s 表 ---\n", m.name)

		// 执行迁移
		if err := database.DB.AutoMigrate(m.model); err != nil {
			fmt.Printf("❌ %s 表迁移失败: %v\n", m.name, err)

			// 尝试获取更详细的错误信息
			if err := database.DB.Exec("SELECT 1").Error; err != nil {
				fmt.Printf("数据库连接异常: %v\n", err)
			}
		} else {
			fmt.Printf("✅ %s 表迁移成功\n", m.name)
			successCount++
		}
	}

	fmt.Printf("\n=== 迁移完成: %d/%d 表成功 ===\n", successCount, len(models))
}
