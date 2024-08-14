package models

type AuthInput struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
