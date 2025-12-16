package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sqlite_test/Model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 创建sqlite数据库
func InitSQLite() {
	// 使用绝对路径
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("获取工作目录失败:", err)
	}
	// 创建文件夹
	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("创建data目录失败:", err)
	}
	// 生成数据库文件路径
	dbFile := filepath.Join(dataDir, "user.db")
	fmt.Printf("数据库文件完整路径: %s\n", dbFile)

	// 测试文件权限
	if err := testFilePermissions(dbFile); err != nil {
		log.Fatal("文件权限测试失败:", err)
	}

	var openErr error
	DB, openErr = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if openErr != nil {
		log.Printf("数据库连接详细错误: %v", openErr)
		log.Printf("错误类型: %T", openErr)
		panic("连接数据库失败，详细错误: " + openErr.Error())
	}

	// 自动迁移所有表结构
	err = DB.AutoMigrate(
		&Model.User{},
		&Model.Tag{},
		&Model.Content{},
		&Model.FileRecord{},
		&Model.ContentTag{},
		&Model.ContentFile{},
		&Model.Comment{},
		&Model.EmailVerify{},
	)
	if err != nil {
		log.Fatal("自动迁移表结构失败:", err)
	}

	fmt.Println("成功连接到数据库并创建所有表结构")
}

// 测试文件权限
func testFilePermissions(dbFile string) error {
	// 尝试创建测试文件
	testFile := dbFile + ".test"
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("无法创建测试文件: %v", err)
	}
	file.Close()

	// 删除测试文件
	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("无法删除测试文件: %v", err)
	}

	fmt.Println("文件权限测试通过")
	return nil
}
