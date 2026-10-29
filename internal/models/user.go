package models

// User đại diện cho bảng users trong database
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // Thêm tag json:"-" để ẩn password khi marshal JSON
}
