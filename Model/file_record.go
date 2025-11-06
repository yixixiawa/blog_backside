package Model

import (
	"time"
)

type FileRecord struct {
	FileID       uint      `gorm:"primaryKey;autoIncrement" json:"file_id"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	StorageName  string    `gorm:"type:varchar(255);not null" json:"storage_name"`
	FilePath     string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileURL      string    `gorm:"type:varchar(500);not null" json:"file_url"`
	FileSize     int64     `gorm:"type:bigint;not null" json:"file_size"`
	FileType     *string   `gorm:"type:varchar(100)" json:"file_type"`
	ContentID    *uint     `gorm:"type:bigint" json:"content_id"`
	UploadTime   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"upload_time"`
	Status       string    `gorm:"type:enum('active','deleted');not null;default:'active'" json:"status"`

	// 关联关系
	Content *Content `gorm:"foreignKey:ContentID" json:"content"`
}

// TableName 指定表名为 file_record
func (FileRecord) TableName() string {
	return "file_record"
}
