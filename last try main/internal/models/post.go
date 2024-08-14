package models

import "time"

type Post struct {
	PostID    int       `json:"post_id" gorm:"primary_key; not null"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	User_id   int       `json:"user_id"`
}
