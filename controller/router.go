package controller

import (
	"blog/utils"

	"github.com/gin-gonic/gin"
)

// 认证

func InitRoutes(r *gin.Engine) {
	fileGroup := r.Group("/file")
	{
		// 公开访问
		fileGroup.GET("/listimg", ListImages)           // 获取所有图片
		fileGroup.GET("/content/:id", GetContentImages) // 根据文章ID获取图片

		// 需要认证的上传接口
		authFile := fileGroup.Group("")
		authFile.Use(utils.JWTAuthMiddleware())
		{
			authFile.POST("/uploadimg", uploadimg) // 上传图片
			authFile.POST("/uploadfile", UploadFile)
		}
	}

	// 静态文件服务，用于直接访问上传的图片
	r.Static("/img", GetImageStoragePath())

	// 用户相关路由
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", UserRegister)
		userGroup.POST("/login", UserLogin)

		// 需要认证的路由
		auth := userGroup.Group("")
		auth.Use(utils.JWTAuthMiddleware())
		{
			auth.POST("/logout", UserLogout)
			auth.PUT("/password", ChangePassword)
			auth.GET("/list", ListUsers) // 获取用户列表
			auth.GET("/:id", GetUserInfo)
			auth.PUT("/:id", UpdateUserProfile)
			auth.DELETE("/:id", DeleteUser)
		}
	}

	// 内容相关路由（GET 为公开，其他需要认证）
	contentGroup := r.Group("/content")
	{
		// 公开读接口
		contentGroup.GET("", ListContents)
		contentGroup.GET("/:id", GetContent)

		// 需要认证的写接口
		authContent := contentGroup.Group("/content_auth")
		authContent.Use(utils.JWTAuthMiddleware())
		{
			authContent.POST("", CreateContent)
			authContent.PUT("/:id", UpdateContent)
			authContent.DELETE("/:id", DeleteContent)
			authContent.POST("/:id/tags", AddContentTags)
			authContent.DELETE("/:id/tags/:tagId", RemoveContentTag)
		}
	}

	// 评论相关路由（全部需要认证）
	commentGroup := r.Group("/comment")
	commentGroup.Use(utils.JWTAuthMiddleware())
	{
		commentGroup.POST("", CreateComment)
		commentGroup.GET("/content/:contentId", ListContentComments)
		commentGroup.PUT("/:id", UpdateComment)
		commentGroup.DELETE("/:id", DeleteComment)
	}

	// 标签相关路由（GET 为公开，其他需要认证）
	tagGroup := r.Group("/tag")
	{
		// 公开读接口
		tagGroup.GET("", ListTags)
		tagGroup.GET("/:id", GetTag)

		// 需要认证的写接口
		authTag := tagGroup.Group("")
		authTag.Use(utils.JWTAuthMiddleware())
		{
			authTag.POST("", CreateTag)
			authTag.PUT("/:id", UpdateTag)
			authTag.DELETE("/:id", DeleteTag)
		}
	}

	// 邮件相关路由（全部需要认证）
	emailGroup := r.Group("/email")
	{
		emailGroup.POST("/verify", SendVerificationEmail)
		emailGroup.POST("/verify/check", CheckVerificationCode)
	}
	goodsGroup := r.Group("/goods")
	{
		goodsGroup.GET("/items", search_goods)
	}
}

func SetupMiddlewares(r *gin.Engine) {
	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 日志中间件
	r.Use(gin.Logger())

	// 恢复中间件
	r.Use(gin.Recovery())

	// 可以添加更多中间件
	// r.Use(authMiddleware()) // 认证中间件
	// r.Use(rateLimitMiddleware()) // 限流中间件
}

func InitializeServer() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode) // 生产环境使用

	// 创建Gin引擎
	r := gin.New()

	// 设置中间件
	SetupMiddlewares(r)

	// 初始化路由
	InitRoutes(r)

	return r
}
