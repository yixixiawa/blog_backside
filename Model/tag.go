package Model

import "time"

type Tag struct {
	TagID        uint      `gorm:"primaryKey;autoIncrement" json:"tag_id"` // 改为 uint
	TagName      string    `gorm:"type:varchar(30);not null;uniqueIndex" json:"tag_name"`
	TagAlias     string    `gorm:"type:varchar(30);default:''" json:"tag_alias"`
	Description  string    `gorm:"type:varchar(255);default:''" json:"description"`
	CreateTime   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"update_time"`
	IsActive     bool      `gorm:"default:true" json:"is_active"` // SQLite 用 bool
	Icon         string    `gorm:"type:varchar(255);default:''" json:"icon"`
	Color        string    `gorm:"type:varchar(50);default:''" json:"color"`
	DisplayOrder int       `gorm:"type:integer;default:0" json:"display_order"`
}
