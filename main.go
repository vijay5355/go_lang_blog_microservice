package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	router.POST("/posts", createPost)
	router.GET("/posts", getAllPosts)
	router.GET("/posts/:id", getPost)
	router.PUT("/posts/:id", updatePost)
	router.DELETE("/posts/:id", deletePost)
	router.GET("/posts/count", getPostsCount)
	fmt.Println("Server running at http://localhost:8080")

	router.Run(":8080")
}
