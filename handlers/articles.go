package handlers

import (
	"crud/db"
	"crud/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type articleService struct {
	db *db.DB
}

// для того чтобы иметь доступ к методам созданной структуры из main (GetAll, GetById и т.д.)
func NewArticleService(db *db.DB) models.ArticleService {
	return &articleService{
		db: db,
	}
}

func (s articleService) GetAll(c *gin.Context) {
	articles, err := s.db.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (s articleService) GetById(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	article, err := s.db.GetById(i)
	if err != nil {
		if err.Error() == "failed to find article" {
			c.JSON(400, gin.H{"error": err.Error()})
			fmt.Printf("error: %v", err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error: %v", err)
		return
	}
	c.JSON(http.StatusOK, article)
}

func (s articleService) GetByAuthor(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	articles, err := s.db.GetByAuthor(i)
	if err != nil {
		if err.Error() == "failed to find author" {
			c.JSON(400, gin.H{"error": err.Error()})
			fmt.Printf("error: %v", err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error: %v", err)
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (s articleService) CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.BindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}
	if err := s.db.CreateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}
	c.JSON(http.StatusOK, "article was created")
}

func (s articleService) UpdateArticle(c *gin.Context) {
	var article models.Article

	err := c.BindJSON(&article)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}

	err = s.db.UpdateArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}

	c.JSON(http.StatusOK, "article was updated")
}

func (s articleService) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)

	err := s.db.DeleteArticle(i)
	if err != nil {
		if err.Error() == "failed to find article" {
			c.JSON(400, gin.H{"error": err.Error()})
			fmt.Printf("error: %v", err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error: %v", err)
		return
	}

	c.JSON(http.StatusOK, "article was deleted")
}
