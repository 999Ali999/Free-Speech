package controllers

import (
	"blogging/internal/database"
	"blogging/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	database.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.UserID != 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "username already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	database.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func Login(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	database.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password",
			"error": err.Error()})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.UserID,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to generate token", "error": err.Error()})
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func GetCurrentUserInfo(c *gin.Context) {
	user, _ := c.Get("currentUser")

	c.JSON(200, gin.H{
		"Current user": user,
	})
}
