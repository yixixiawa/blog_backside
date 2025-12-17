package utils

import (
	"blog/constants"
	"blog/database"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 强制验证 token（用于需要登录的路由）
func JWTAuthMiddleware() gin.HandlerFunc {
	return jwtAuthMiddleware(true)
}

// JWTAuthOptionalMiddleware 可选验证 token（用于公开路由，但若带 token 则解析设置上下文）
func JWTAuthOptionalMiddleware() gin.HandlerFunc {
	return jwtAuthMiddleware(false)
}

func jwtAuthMiddleware(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if required {
				constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "missing Authorization header"})
				c.Abort()
				return
			}
			// 非必须：直接放行，不设置 user_id
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			if required {
				constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "invalid Authorization header"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		token := parts[1]
		claims, err := ParseToken(token)
		if err != nil {
			if required {
				constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// 可选：校验 token 是否在 Redis 中（用于支持主动登出 / 单设备登录等）
		tokenKey := fmt.Sprintf("user_token:%d", claims.UserID)
		stored, err := database.GetString(tokenKey)
		if err != nil || stored != token {
			if required {
				constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "token expired or revoked"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// 将用户信息放入上下文，供 handler 使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
