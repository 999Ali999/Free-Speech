package models

type User struct {
	UserID   uint   `json:"user_id" gorm:"primary_key; not null"`
	Email    string `json:"email" gorm:"unique; not null"`
	Password string `json:"password" gorm:"not null"`
}
