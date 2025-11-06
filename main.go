package main

import (
	"fmt"
	"sqlite_test/database"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("应用程序启动...")

	// 初始化数据库
	database.InitSQLite()
	database.InitRedis()
	// 自动构建/迁移
	// 初始化路由
	router := gin.Default()

	fmt.Println("应用程序启动完成")

	fmt.Println("应用程序启动完成")
	router.Run(":18800")
}
