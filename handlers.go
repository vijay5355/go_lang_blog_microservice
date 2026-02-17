package main

// package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const fileName = "posts.json"

var db *sql.DB

type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func initDB() {
	var err error
	// connStr:= "postgres://postgres:HelloWorld@localhost:5432/blogdb?sslmode=disable"
	connStr := "postgres://postgres:HelloWorld@localhost:5432/blogdb?sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established")
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

// func generateID() string {
// 	return fmt.Sprintf("%d", time.Now().UnixMilli())
// 	// return time.Now().UnixMilli()
// }

func generateID() int64 {
	return time.Now().UnixNano()
}

// func getTime() string {
// 	return time.Now().Format(time.RFC3339)
// }

func createPost(c *gin.Context) {

	var newPost Post
	var created time.Time
	var updated time.Time

	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if newPost.Title == "" || newPost.Content == "" || newPost.Author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields"})
		return
	}

	query := `INSERT INTO posts (id,title,content,author)
	VALUES ($1,$2,$3,$4)
	RETURNING createdat,updatedat`

	newPost.ID = strconv.FormatInt(generateID(), 10)
	if err := db.QueryRowContext(context.Background(),
		query, newPost.ID, newPost.Title, newPost.Content, newPost.Author).
		Scan(&created, &updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newPost.CreatedAt = created.Format(time.RFC3339)
	newPost.UpdatedAt = updated.Format(time.RFC3339)

	posts, err := loadPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}
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

// func getPostsCount(c *gin.Context) {
// 	posts, err := loadPosts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
// 	}
// 	var count int = len(posts)
// 	c.JSON(http.StatusOK, gin.H{
// 		"The total no of posts is given by": count,
// 	})
// }

// func getPost(c *gin.Context) {

// 	id := c.Param("id")

// 	posts, err := loadPosts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
// 		return
// 	}

// 	for _, post := range posts {

// 		if post.ID == id {
// 			c.JSON(http.StatusOK, post)
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
// }

// func updatePost(c *gin.Context) {

// 	id := c.Param("id")

// 	var updatedPost Post
// 	if err := c.ShouldBindJSON(&updatedPost); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}
// 	posts, err := loadPosts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
// 		return
// 	}
// 	for i, post := range posts {

// 		if post.ID == id {
// 			// post.Title = updatedPost.Title
// 			// post.Content = updatedPost.Content
// 			// post.Author = updatedPost.Author

// 			if updatedPost.Title != "" {
// 				post.Title = updatedPost.Title
// 			}
// 			if updatedPost.Content != "" {
// 				post.Content = updatedPost.Content
// 			}
// 			if updatedPost.Author != "" {
// 				post.Author = updatedPost.Author
// 			}

// 			post.UpdatedAt = getTime()

// 			posts[i] = post
// 			if err := savePosts(posts); err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "File write error"})
// 				return
// 			}

// 			c.JSON(http.StatusOK, post)
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
// }

// func deletePost(c *gin.Context) {

// 	id := c.Param("id")
// 	posts, err := loadPosts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
// 		return
// 	}

// 	for i, post := range posts {

// 		if post.ID == id {
// 			posts = append(posts[:i], posts[i+1:]...)
// 			if err := savePosts(posts); err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "File write error"})
// 				return
// 			}

// 			c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
// }
