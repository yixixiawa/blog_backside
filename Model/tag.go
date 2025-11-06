package Model

import (
	"time"
)

type Tag struct {
	TagID        int       `gorm:"primaryKey;autoIncrement" json:"tag_id"`
	TagName      string    `gorm:"type:varchar(30);not null;uniqueIndex" json:"tag_name"`
	TagAlias     *string   `gorm:"type:varchar(30)" json:"tag_alias"`
	Description  *string   `gorm:"type:varchar(255)" json:"description"`
	CreateTime   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"update_time"`
	IsActive     bool      `gorm:"type:tinyint(1);not null;default:1" json:"is_active"`
	Icon         *string   `gorm:"type:varchar(255)" json:"icon"`
	Color        *string   `gorm:"type:varchar(50)" json:"color"`
	DisplayOrder int       `gorm:"type:int;not null;default:0" json:"display_order"`

	// 关联关系
	Contents []Content `gorm:"many2many:content_tag" json:"contents"`
}

// TableName 指定表名为 tag
func (Tag) TableName() string {
	return "tag"
}
