package Model

import (
	"time"
)

// OAuthState OAuth状态表（用于CSRF防护）
type OAuthState struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	State      string    `gorm:"size:100;uniqueIndex;not null" json:"state"`
	UserID     uint      `gorm:"index" json:"user_id"` // 如果用户已登录，记录用户ID
	PlatformID uint      `gorm:"index" json:"platform_id"`
	ExpiresAt  time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联关系
	User     User          `gorm:"foreignKey:UserID;references:UserID" json:"-"`
	Platform OAuthPlatform `gorm:"foreignKey:PlatformID;references:OAuthID" json:"-"`
}

// TableName 指定表名
func (OAuthState) TableName() string {
	return "oauth_states"
}
