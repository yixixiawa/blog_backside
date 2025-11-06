package Model

// EmailVerify 如果你只想记录邮箱验证历史，可以简化为：
type EmailVerify struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Email      string `gorm:"type:varchar(255);not null" json:"email"`
	VerifyCode string `gorm:"type:varchar(100);not null" json:"verify_code"`
	Purpose    string `gorm:"type:enum('registration','password_reset','email_update');default:'registration'" json:"purpose"`
	IsUsed     bool   `gorm:"type:boolean;not null;default:false" json:"is_used"`
}
