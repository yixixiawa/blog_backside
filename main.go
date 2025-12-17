package main

import (
	"blog/AutoMigrate"
	"blog/controller"
	"blog/database"
	"log"
)

func main() {
	// 初始化DB/Redis（修改 Init 函数以返回 error 更可靠）
	database.InitSQLite()
	database.InitRedis()

	// 自动迁移
	AutoMigrate.Generation_sql()

	// 使用 controller 提供的引擎（已注册路由）
	router := controller.InitializeServer()

	// 可选：打印已注册路由用于调试
	for _, r := range router.Routes() {
		log.Printf("%s %s\n", r.Method, r.Path)
	}

	// 启动
	if err := router.Run(":18800"); err != nil {
		log.Fatal(err)
	}
}
