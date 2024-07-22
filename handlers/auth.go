package handlers

import (
	"crud/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h BaseHandler) Register(c *gin.Context) {
	var a models.Auth

	err := c.BindJSON(&a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Printf("error:%v", err)
		return
	}

	if a.Username == "" || a.Login == "" || a.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, login or password is empty"})
		fmt.Printf("error:%v", "username, login or password is empty")
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Printf("failed to generate hash from password: %v", err)
		return
	}

	a.Hash = string(passHash)

	err = h.db.Register(a)
	if err != nil {

		if err.Error() == "user with this login already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Printf("error:%v", err)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("%v", err)
		return
	}

	c.JSON(http.StatusOK, "user was registered")
}

func (h BaseHandler) Login(c *gin.Context) {
	var a models.Auth

	err := c.BindJSON(&a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Printf("error:%v", err.Error())
		return
	}

	user, err := h.db.GetUser(a)
	if err != nil {
		if err.Error() == "failed to find user with this login" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Printf("error:%v", err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(a.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["Uid"] = user.Id
	claims["Login"] = user.Login
	claims["Exp"] = time.Now().Add(time.Hour * 2).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error:%v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
