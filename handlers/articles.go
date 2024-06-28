package handlers

import (
	"crud/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h BaseHandler) GetAll(c *gin.Context) {
	articles, err := h.db.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Printf("error:%v", err)
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (h BaseHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	article, err := h.db.GetById(i)
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

func (h BaseHandler) GetByAuthor(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	articles, err := h.db.GetByAuthor(i)
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

func (h BaseHandler) CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.BindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Printf("error:%v", err)
		return
	}
	if err := h.db.CreateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Printf("error:%v", err)
		return
	}
	c.JSON(http.StatusOK, "article was created")
}

func (h BaseHandler) UpdateArticle(c *gin.Context) {
	var article models.Article

	err := c.BindJSON(&article)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Printf("error:%v", err)
		return
	}

	err = h.db.UpdateArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Printf("error:%v", err)
		return
	}

	c.JSON(http.StatusOK, "article was updated")
}

func (h BaseHandler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)

	err := h.db.DeleteArticle(i)
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
