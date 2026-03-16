package security

import (
	"github.com/gin-gonic/gin"
)

// SetupCORSMiddleware 返回 Gin CORS 中间件
// 配置允许跨域请求的来源、方法和请求头
func SetupCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源（可根据环境修改为具体域名）
		origin := c.Request.Header.Get("Origin")

		// 安全的做法：只允许特定域名，防止 CSRF 攻击
		// 生产环境请将 allowedOrigins 改为你实际的前端域名
		allowedOrigins := map[string]bool{
			"http://localhost:23357": true,
			"http://yixixiawa.xyz":   true,
		}

		// 检查 Origin 是否在白名单中
		if origin != "" && allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// 允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 允许客户端访问的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin")

		// 浏览器可以缓存预检请求的结果（单位：秒，这里设为 24 小时）
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 允许携带认证信息（如 Cookie、Authorization Header）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
