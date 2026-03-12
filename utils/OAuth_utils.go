package utils

import (
	"strings"
)

// ParseScopes 解析scope字符串为切片
func ParseScopes(scopes string) []string {
	if scopes == "" {
		return []string{}
	}
	return strings.Split(scopes, ",")
}

// GetPlatformAuthURL 获取平台授权URL
func GetPlatformAuthURL(platform string) string {
	switch platform {
	case "github":
		return "https://github.com/login/oauth/authorize"
	case "google":
		return "https://accounts.google.com/o/oauth2/auth"
	case "wechat":
		return "https://open.weixin.qq.com/connect/qrconnect"
	default:
		return ""
	}
}

// GetPlatformTokenURL 获取平台token URL
func GetPlatformTokenURL(platform string) string {
	switch platform {
	case "github":
		return "https://github.com/login/oauth/access_token"
	case "google":
		return "https://oauth2.googleapis.com/token"
	case "wechat":
		return "https://api.weixin.qq.com/sns/oauth2/access_token"
	default:
		return ""
	}
}

// GetPlatformUserInfoURL 获取平台用户信息URL
func GetPlatformUserInfoURL(platform string) string {
	switch platform {
	case "github":
		return "https://api.github.com/user"
	case "google":
		return "https://www.googleapis.com/oauth2/v2/userinfo"
	case "wechat":
		return "https://api.weixin.qq.com/sns/userinfo"
	default:
		return ""
	}
}

// GetPlatformScopes 获取平台默认scope
func GetPlatformScopes(platform string) string {
	switch platform {
	case "github":
		return "user:email"
	case "google":
		return "https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/userinfo.profile"
	case "wechat":
		return "snsapi_login"
	default:
		return ""
	}
}
