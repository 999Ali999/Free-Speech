package models

import (
	"time"
)

type Reaction struct {
	ReactionId   int       `json:"reaction_id" gorm:"primary_key; not null"`
	CreatedAt    time.Time `json:"created_at"`
	ReactionType string    `gorm:"type:varchar(20);check:reaction_type IN ('like', 'dislike'); not null"`
	UserId       int       `json:"user_id" gorm:"not null"`
	PostId       *int      `gorm:"null"`
	CommentId    *int      `gorm:"null"`
}
