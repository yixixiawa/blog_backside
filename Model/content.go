package Model

import "time"

type Content struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title             string     `gorm:"type:varchar(255);not null" json:"title"`
	Content           string     `gorm:"type:text" json:"content"` // SQLite 不支持 longtext
	BriefIntroduction string     `gorm:"type:text" json:"brief_introduction"`
	UserID            uint       `gorm:"not null" json:"user_id"`
	CoverImage        string     `gorm:"type:varchar(500)" json:"cover_image"`
	Status            string     `gorm:"type:varchar(20);default:'draft';check:status IN ('draft','published','archived')" json:"status"` // 替代 enum
	ViewCount         int        `gorm:"default:0" json:"views"`
	Likes             int        `gorm:"default:0" json:"likes"`
	CommentCount      int        `gorm:"default:0" json:"comment_count"`
	CreatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	PublishedAt       *time.Time `json:"published_at"`
	ThumbnailURL      string     `gorm:"type:varchar(500)" json:"thumbnail_url"`

	// 关联关系
	User     User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tags     []Tag        `gorm:"many2many:content_tags;" json:"tags,omitempty"`
	Comments []Comment    `gorm:"foreignKey:ContentID" json:"comments,omitempty"`
	Files    []FileRecord `gorm:"foreignKey:ContentID" json:"files,omitempty"`
}
