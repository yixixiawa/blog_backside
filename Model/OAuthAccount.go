package Model

import (
	"time"
)

// OAuthAccount 第三方账号关联表
type OAuthAccount struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"`      // 外键约束：级联更新和删除
	PlatformID        uint       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"platform_id"` // 外键约束：级联更新，限制删除
	PlatformUserID    string     `gorm:"size:100;not null;index:idx_platform_user,unique" json:"platform_user_id"`         // 第三方平台的用户ID
	PlatformUserName  string     `gorm:"size:100" json:"platform_user_name"`                                               // 第三方平台的用户名
	PlatformUserEmail string     `gorm:"size:100" json:"platform_user_email"`                                              // 第三方平台的邮箱
	AvatarURL         string     `gorm:"size:500" json:"avatar_url"`                                                       // 第三方头像
	AccessToken       string     `gorm:"size:500" json:"-"`                                                                // 不返回给前端
	RefreshToken      string     `gorm:"size:500" json:"-"`
	TokenExpiresAt    *time.Time `json:"token_expires_at"`
	RawData           string     `gorm:"type:text" json:"-"` // 第三方返回的原始数据
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// 关联关系
	User     User          `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	Platform OAuthPlatform `gorm:"foreignKey:PlatformID;references:OAuthID" json:"platform,omitempty"`
}

// TableName 指定表名
func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}
