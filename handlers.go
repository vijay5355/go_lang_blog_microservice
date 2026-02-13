package main

// package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const fileName = "posts.json"

type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func loadPosts() ([]Post, error) {

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := os.WriteFile(fileName, []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var posts []Post
	err = json.Unmarshal(data, &posts)

	return posts, err
}
func savePosts(posts []Post) error {

	data, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func createPost(c *gin.Context) {

	var newPost Post

	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if newPost.Title == "" || newPost.Content == "" || newPost.Author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields"})
		return
	}

	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}

	now := getTime()

	newPost.ID = generateID()
	newPost.CreatedAt = now
	newPost.UpdatedAt = now

	posts = append(posts, newPost)

	if err := savePosts(posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File write error"})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

func getAllPosts(c *gin.Context) {

	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func getPostsCount(c *gin.Context) {
	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
	}
	var count int = len(posts)
	c.JSON(http.StatusOK, gin.H{
		"The total no of posts is given by": count,
	})
}

func getPost(c *gin.Context) {

	id := c.Param("id")

	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}

	for _, post := range posts {

		if post.ID == id {
			c.JSON(http.StatusOK, post)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
}

func updatePost(c *gin.Context) {

	id := c.Param("id")

	var updatedPost Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}
	for i, post := range posts {

		if post.ID == id {
			post.Title = updatedPost.Title
			post.Content = updatedPost.Content
			post.Author = updatedPost.Author
			post.UpdatedAt = getTime()

			posts[i] = post
			if err := savePosts(posts); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "File write error"})
				return
			}

			c.JSON(http.StatusOK, post)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
}

func deletePost(c *gin.Context) {

	id := c.Param("id")
	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}

	for i, post := range posts {

		if post.ID == id {
			posts = append(posts[:i], posts[i+1:]...)
			if err := savePosts(posts); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "File write error"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
}
