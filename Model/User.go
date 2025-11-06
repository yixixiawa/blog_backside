package Model

type User struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username string `gorm:"size:30;not null;uniqueIndex" json:"username"`
	Password string `gorm:"size:100;not null" json:"password"`
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`
	IsAdmin  bool   `gorm:"not null;default:false" json:"is_admin"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
