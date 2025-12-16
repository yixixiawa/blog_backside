package Model

import "time"

type ContentFile struct {
	ContentID uint      `gorm:"primaryKey;not null" json:"content_id"`
	FileID    uint      `gorm:"primaryKey;not null" json:"file_id"`
	Order     int       `gorm:"default:0" json:"order"`
	Usage     string    `gorm:"type:varchar(50)" json:"usage"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系（用于预加载）
	Content    Content    `gorm:"foreignKey:ContentID;references:ID" json:"content,omitempty"`
	FileRecord FileRecord `gorm:"foreignKey:FileID;references:FileID" json:"file_record,omitempty"`
}

// TableName 指定表名
func (ContentFile) TableName() string {
	return "content_files"
}
