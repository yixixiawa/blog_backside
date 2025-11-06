package Model

import "time"

type Comment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CommentText string    `gorm:"type:text;not null" json:"comment_text"` // 重命名避免冲突
	UserID      uint      `gorm:"not null" json:"user_id"`
	ContentID   uint      `gorm:"not null" json:"content_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Content Content `gorm:"foreignKey:ContentID" json:"content"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
