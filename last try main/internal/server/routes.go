package server

import (
	"blogging/internal/controllers"
	"blogging/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(middlewares.NoCacheMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	// Auth
	r.POST("/signup", controllers.CreateUser)
	r.POST("/login", controllers.Login)

	// User
	r.GET("/users/profile", middlewares.CheckAuth, controllers.GetCurrentUserInfo)

	// Posts
	r.GET("/posts", controllers.GetPosts)
	r.GET("/posts/:id", controllers.GetPost)

	r.POST("/posts", middlewares.CheckAuth, controllers.CreatePost)       // protected route
	r.PUT("/posts/:id", middlewares.CheckAuth, controllers.UpdatePost)    // protected route
	r.DELETE("/posts/:id", middlewares.CheckAuth, controllers.DeletePost) // protected route

	// Comments
	r.GET("/posts/:id/comments", controllers.GetCommentsForSpecificPost)

	r.POST("/posts/:id/comments", middlewares.CheckAuth, controllers.CreateComment) // protected route
	r.DELETE("/comments/:id", middlewares.CheckAuth, controllers.DeleteComment)     // protected route

	// Likes and Dislikes for posts
	r.GET("/posts/:id/likes", controllers.GetLikesForSpecificPost)
	r.GET("/posts/:id/dislikes", controllers.GetDislikesForSpecificPost)

	r.POST("/posts/:id/likes", middlewares.CheckAuth, controllers.CreateLikeForSpecificPost)       // protected route
	r.POST("/posts/:id/dislikes", middlewares.CheckAuth, controllers.CreateDislikeForSpecificPost) // protected route

	return r
}
