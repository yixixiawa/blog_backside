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
	UploadTime   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"upload_time"`
	Status       string    `gorm:"type:varchar(20);default:'active';check:status IN ('active','deleted')" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联关系 - 多对多通过 content_files 表
	Contents     []Content     `gorm:"many2many:content_files;foreignKey:FileID;joinForeignKey:FileID;References:ID;joinReferences:ContentID" json:"contents,omitempty"`
	ContentFiles []ContentFile `gorm:"foreignKey:FileID;references:FileID" json:"content_files,omitempty"`
}

// TableName 指定表名
func (FileRecord) TableName() string {
	return "file_records"
}
