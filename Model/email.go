package Model

import "time"

type EmailVerify struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email      string    `gorm:"type:varchar(255);not null" json:"email"`
	VerifyCode string    `gorm:"type:varchar(100);not null" json:"verify_code"`
	Purpose    string    `gorm:"type:varchar(50);default:'registration'" json:"purpose"` // 替代 enum
	IsUsed     bool      `gorm:"default:false" json:"is_used"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}
