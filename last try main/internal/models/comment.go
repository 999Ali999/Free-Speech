package models

import "time"

type Comment struct {
	CommentId int       `json:"comment_id" gorm:"primary_key; not null"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UserId    int       `json:"user_id"`
	PostId    int       `json:"post_id"`
}
