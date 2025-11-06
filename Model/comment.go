package Model

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ContentID uint      `gorm:"not null" json:"content_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
