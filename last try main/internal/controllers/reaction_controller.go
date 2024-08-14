package controllers

import (
	"blogging/internal/database"
	"blogging/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLikesForSpecificPost(c *gin.Context) {
	id := c.Param("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
	}

	var count int64
	result := database.DB.Model(&models.Reaction{}).Where("post_id = ? AND reaction_type = ?", postId, "like").Count(&count)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "This post has no likes yet."})
		return
	}

	var reaction models.Reaction
	result = database.DB.Where("post_id = ? AND reaction_type = ?", postId, "like").Find(&reaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, reaction)
}

func GetDislikesForSpecificPost(c *gin.Context) {
	id := c.Param("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
	}

	var count int64
	result := database.DB.Model(&models.Reaction{}).Where("post_id = ? AND reaction_type = ?", postId, "dislike").Count(&count)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "This post has no dislikes yet."})
		return
	}

	var reaction models.Reaction
	result = database.DB.Where("post_id = ? AND reaction_type = ?", postId, "dislike").Find(&reaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, reaction)
}

func CreateLikeForSpecificPost(c *gin.Context) {
	var user models.User

	if value, exists := c.Get("currentUser"); exists {
		user = value.(models.User)
	}
	userId := user.UserID

	id := c.Param("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}

	var reaction models.Reaction
	err = c.ShouldBindJSON(&reaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reaction.ReactionType = "like"
	reaction.UserId = int(userId)
	reaction.PostId = &postId
	reaction.CommentId = nil

	// same like check!
	var existingReaction models.Reaction
	if err := database.DB.Where("user_id = ? AND post_id = ? AND reaction_type = ?", userId, postId, "like").First(&existingReaction).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has already reacted to this post"})
		return
	}

	reaction.CreatedAt = time.Now()

	if err := database.DB.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reaction)
}

func CreateDislikeForSpecificPost(c *gin.Context) {
	var user models.User

	if value, exists := c.Get("currentUser"); exists {
		user = value.(models.User)
	}
	userId := user.UserID

	id := c.Param("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}

	var reaction models.Reaction
	err = c.ShouldBindJSON(&reaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reaction.ReactionType = "dislike"
	reaction.UserId = int(userId)
	reaction.PostId = &postId
	reaction.CommentId = nil

	// same dislike check!
	var existingReaction models.Reaction
	if err := database.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&existingReaction).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has already reacted to this post"})
		return
	}

	reaction.CreatedAt = time.Now()

	if err := database.DB.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reaction)
}
