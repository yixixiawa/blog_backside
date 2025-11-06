package AutoMigrate

import (
	"fmt"
	"sqlite_test/Model"
	"sqlite_test/database"
)

// 定义模型初始化顺序
type migrationStep struct {
	model    interface{}
	name     string
	depends  []string
	migrated bool
}

func Generation_sql() {
	fmt.Println("开始自动生成SQL表结构")

	// 定义迁移步骤及其依赖关系
	steps := []migrationStep{
		{model: &Model.User{}, name: "User"},
		{model: &Model.Tag{}, name: "Tag"},
		{model: &Model.Content{}, name: "Content", depends: []string{"User"}},
		{model: &Model.ContentTag{}, name: "ContentTag", depends: []string{"Content", "Tag"}},
		{model: &Model.Comment{}, name: "Comment", depends: []string{"User", "Content"}},
		{model: &Model.FileRecord{}, name: "FileRecord", depends: []string{"User"}},
	}

	// 执行迁移，直到所有模型都已迁移或无法继续
	for {
		migratedSomething := false
		allMigrated := true

		for i := range steps {
			if steps[i].migrated {
				continue
			}

			canMigrate := true
			// 检查依赖是否已经迁移
			for _, dep := range steps[i].depends {
				dependencyMigrated := false
				for _, s := range steps {
					if s.name == dep && s.migrated {
						dependencyMigrated = true
						break
					}
				}
				if !dependencyMigrated {
					canMigrate = false
					break
				}
			}

			if canMigrate {
				fmt.Printf("正在生成 %s 表结构...\n", steps[i].name)
				if err := database.DB.AutoMigrate(steps[i].model); err != nil {
					fmt.Printf("自动生成 %s 表结构失败: %v\n", steps[i].name, err)
					return
				}
				steps[i].migrated = true
				migratedSomething = true
				fmt.Printf("成功生成 %s 表结构\n", steps[i].name)
			} else {
				allMigrated = false
			}
		}

		// 如果所有模型都已迁移或者本轮没有进行任何迁移，则退出
		if allMigrated || !migratedSomething {
			break
		}
	}

	// 验证所有模型是否都已迁移
	for _, step := range steps {
		if !step.migrated {
			fmt.Printf("警告：%s 表结构未能生成，可能存在循环依赖\n", step.name)
		}
	}

	fmt.Println("所有表结构自动生成完成")
}
