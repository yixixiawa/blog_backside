package Model

import "time"

type ContentTag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ContentID uint      `gorm:"not null;index" json:"content_id"`
	TagID     uint      `gorm:"not null;index" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ContentTag) TableName() string {
	return "content_tags"
}
