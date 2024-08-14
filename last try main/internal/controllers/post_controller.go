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

func GetPosts(c *gin.Context) {
	var posts []models.Post
	result := database.DB.Find(&posts)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Posts not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, posts)
}

func GetPost(c *gin.Context) {
	// Get the id of the requested Post
	postId := c.Param("id")
	postID, err := strconv.Atoi(postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	result := database.DB.First(&post, postID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, post)
}

func CreatePost(c *gin.Context) {
	// Get the id of the user who is creating it
	var user models.User

	if value, exists := c.Get("currentUser"); exists {
		user = value.(models.User)
	}

	userID := user.UserID

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Populate the post with the id of the user who is creating it
	post.CreatedAt = time.Now()
	post.User_id = int(userID)

	database.DB.Create(&post)

	c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
	var currentUser models.User
	if value, exists := c.Get("currentUser"); exists {
		currentUser = value.(models.User)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the post ID from the URL parameter
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Retrieve the post from the database
	var post models.Post
	result := database.DB.First(&post, postID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	// Check if the current user is the owner of the post
	if uint(post.User_id) != currentUser.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this post"})
		return
	}

	// Bind the request body to the post struct
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated post to the database
	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated post as a response
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	// Get the current user from the context
	var currentUser models.User
	if value, exists := c.Get("currentUser"); exists {
		currentUser = value.(models.User)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the post ID from the URL parameter
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Retrieve the post from the database
	var post models.Post
	result := database.DB.First(&post, postID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	// Check if the current user is the owner of the post
	if uint(post.User_id) != currentUser.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this post"})
		return
	}

	// Delete the post
	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
