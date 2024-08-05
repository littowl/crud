package models

import "github.com/gin-gonic/gin"

type ArticleService interface {
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	GetByAuthor(c *gin.Context)
	CreateArticle(c *gin.Context)
	UpdateArticle(c *gin.Context)
	DeleteArticle(c *gin.Context)
}

type User struct {
	Id   int    `json:"user_id" db:"user_id"`
	Name string `json:"user_name,omitempty" db:"user_name"`
}

type Article struct {
	Id      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content,omitempty" db:"content"`
	Author  User   `json:"author"`
}

type AuthService interface {
	Register(c *gin.Context)
	Verify(c *gin.Context)
	Login(c *gin.Context)
}

type Auth struct {
	Id       int
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
	Hash     string
	Link     string
}
