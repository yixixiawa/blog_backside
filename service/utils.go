package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ControllerUtils 控制器工具类
type ControllerUtils struct{}

// 单例实例
var Utils = &ControllerUtils{}

// ValidateEmail 验证邮箱格式
func (u *ControllerUtils) ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// ValidatePassword 验证密码强度
func (u *ControllerUtils) ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "密码长度至少6位"
	}
	if len(password) > 20 {
		return false, "密码长度不能超过20位"
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return false, "密码必须包含大写字母、小写字母和数字"
	}

	return true, ""
}

// ValidateUsername 验证用户名
func (u *ControllerUtils) ValidateUsername(username string) (bool, string) {
	if len(username) < 3 {
		return false, "用户名长度至少3位"
	}
	if len(username) > 20 {
		return false, "用户名长度不能超过20位"
	}

	pattern := `^[a-zA-Z0-9_]+$`
	if !regexp.MustCompile(pattern).MatchString(username) {
		return false, "用户名只能包含字母、数字和下划线"
	}

	return true, ""
}

// GenerateRandomCode 生成随机验证码（使用utils包中的实现）
func (u *ControllerUtils) GenerateRandomCode(length int) string {
	// 导入utils包，使用现有的实现
	// 这里暂时保留简单实现，避免循环依赖
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// ParseIDFromParam 从URL参数中解析ID
func (u *ControllerUtils) ParseIDFromParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("参数 %s 不能为空", paramName)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("参数 %s 必须是有效的数字", paramName)
	}

	return uint(id), nil
}

// GetClientIP 获取客户端IP地址
func (u *ControllerUtils) GetClientIP(c *gin.Context) string {
	// 检查X-Forwarded-For头
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查X-Real-IP头
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		return xri
	}

	// 返回RemoteAddr
	return c.ClientIP()
}

// GetUserAgent 获取用户代理
func (u *ControllerUtils) GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

// BuildLikeQuery 构建LIKE查询条件
func (u *ControllerUtils) BuildLikeQuery(keyword string) string {
	if keyword == "" {
		return ""
	}
	return "%" + strings.TrimSpace(keyword) + "%"
}

// SanitizeString 清理字符串（移除危险字符）
func (u *ControllerUtils) SanitizeString(input string) string {
	// 移除HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	sanitized := re.ReplaceAllString(input, "")

	// 移除SQL注入相关字符
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_"}
	for _, d := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, d, "")
	}

	return strings.TrimSpace(sanitized)
}

// FormatFileSize 格式化文件大小
func (u *ControllerUtils) FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// IsValidImageExtension 验证图片扩展名
func (u *ControllerUtils) IsValidImageExtension(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// GenerateSlug 生成URL友好的slug
func (u *ControllerUtils) GenerateSlug(title string) string {
	// 转为小写
	slug := strings.ToLower(title)

	// 替换空格为连字符
	slug = strings.ReplaceAll(slug, " ", "-")

	// 移除特殊字符
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = re.ReplaceAllString(slug, "")

	// 移除多余的连字符
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// 移除首尾连字符
	slug = strings.Trim(slug, "-")

	return slug
}

// TimeAgo 计算相对时间
func (u *ControllerUtils) TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		return "刚刚"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	}
	if duration < time.Hour*24 {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d小时前", hours)
	}
	if duration < time.Hour*24*7 {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	}
	if duration < time.Hour*24*30 {
		weeks := int(duration.Hours() / 24 / 7)
		return fmt.Sprintf("%d周前", weeks)
	}
	if duration < time.Hour*24*365 {
		months := int(duration.Hours() / 24 / 30)
		return fmt.Sprintf("%d个月前", months)
	}

	years := int(duration.Hours() / 24 / 365)
	return fmt.Sprintf("%d年前", years)
}

// TruncateString 截断字符串
func (u *ControllerUtils) TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}

	// 尝试在单词边界截断
	if length > 3 {
		truncated := s[:length-3]
		if lastSpace := strings.LastIndex(truncated, " "); lastSpace > length/2 {
			return truncated[:lastSpace] + "..."
		}
	}

	return s[:length-3] + "..."
}
