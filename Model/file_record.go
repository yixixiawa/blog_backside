package Model

import "time"

type FileRecord struct {
	FileID       uint      `gorm:"primaryKey;autoIncrement" json:"file_id"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	StorageName  string    `gorm:"type:varchar(255);not null" json:"storage_name"`
	FilePath     string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileURL      string    `gorm:"type:varchar(500);not null" json:"file_url"`
	FileSize     int64     `gorm:"not null" json:"file_size"`
	FileType     string    `gorm:"type:varchar(100)" json:"file_type"`
	ContentID    *uint     `json:"content_id"`
	UploadTime   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"upload_time"`
	Status       string    `gorm:"type:varchar(20);default:'active';check:status IN ('active','deleted')" json:"status"` // 替代 enum

	// 关联关系
	Content *Content `gorm:"foreignKey:ContentID" json:"content"`
}
