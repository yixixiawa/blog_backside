package Model

import (
	"time"
)

type Content struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title             string     `gorm:"type:varchar(255);not null" json:"title"`
	Content           string     `gorm:"type:longtext" json:"content"`
	BriefIntroduction string     `gorm:"type:text;column:brief_introduction" json:"brief_introduction"`
	UserID            uint       `gorm:"not null;column:user_id" json:"user_id"`
	CoverImage        string     `gorm:"type:varchar(500);column:cover_image" json:"cover_image"`
	Status            string     `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
	ViewCount         int        `gorm:"default:0;column:view_count" json:"views"`
	Likes             int        `gorm:"default:0" json:"likes"`
	CommentCount      int        `gorm:"default:0;column:comment_count" json:"comment_count"`
	CreatedAt         time.Time  `gorm:"column:create_time" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:update_time" json:"updated_at"`
	PublishedAt       *time.Time `gorm:"column:published_at" json:"published_at"`
	ThumbnailURL      string     `gorm:"type:varchar(500);column:thumbnail_url" json:"thumbnail_url"`

	// 关联关系
	User     User         `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	Tags     []Tag        `gorm:"many2many:content_tag;" json:"tags,omitempty"`
	Comments []Comment    `gorm:"foreignKey:ContentID" json:"comments,omitempty"`
	Files    []FileRecord `gorm:"foreignKey:ContentID" json:"files,omitempty"`
}

func (Content) TableName() string {
	return "content"
}
