package Model

import (
	"time"
)

// OAuthPlatform OAuth平台配置表
type OAuthPlatform struct {
	OAuthID      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Platform     string    `gorm:"size:50;not null;uniqueIndex" json:"platform"` // github, google, wechat等
	DisplayName  string    `gorm:"size:100" json:"display_name"`                 // 显示名称
	ClientID     string    `gorm:"size:255;not null" json:"client_id"`
	ClientSecret string    `gorm:"size:255;not null" json:"-"`
	RedirectURL  string    `gorm:"size:255;not null" json:"redirect_url"`
	AuthURL      string    `gorm:"size:255" json:"auth_url"`
	TokenURL     string    `gorm:"size:255" json:"token_url"`
	UserInfoURL  string    `gorm:"size:255" json:"user_info_url"`
	Scopes       string    `gorm:"size:255" json:"scopes"` // 逗号分隔
	IconURL      string    `gorm:"size:255" json:"icon_url"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
	IsEnabled    bool      `gorm:"default:true" json:"is_enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (OAuthPlatform) TableName() string {
	return "oauth_platforms"
}
