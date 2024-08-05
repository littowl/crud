package handlers

import (
	"crud/db"
	"crud/models"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	db    *db.DB
	cache *cache.Cache
}

func NewAuthService(db *db.DB, cache *cache.Cache) models.AuthService {
	return &authService{
		db:    db,
		cache: cache,
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateHash() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s authService) Register(c *gin.Context) {
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

	user, err := s.db.GetUser(models.Auth{Login: a.Login})
	if user != (models.Auth{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with this email is already existed"})
		fmt.Printf("%v", err)
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Printf("failed to generate hash from password: %v", err)
		return
	}

	a.Hash = string(passHash)
	a.Link = fmt.Sprintf("127.0.0.1:5000/auth/verify?login=%s&hash=%v", a.Login, generateHash())

	// Sender data.
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	// Receiver email address.
	to := []string{
		a.Login,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Confirmation\n"
	message := []byte(subject + "\n" + a.Link)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	s.cache.Set(a.Login, a, time.Hour*2)
	c.JSON(http.StatusOK, "Email Sent Successfully!")
}

func (s authService) Verify(c *gin.Context) {
	login, ok := c.GetQuery("login")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect link"})
		fmt.Printf("error while getting login: %v", "incorrect link")
		return
	}

	data, found := s.cache.Get(login)

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link has expired"})
		fmt.Printf("error: %v", "link has expired")
		return
	}
	fmt.Print(c.Request.Host+c.Request.RequestURI, data.(models.Auth).Link)
	if (c.Request.Host + c.Request.RequestURI) != data.(models.Auth).Link {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect link"})
		fmt.Printf("error while hash check: %v", "incorrect link")
		return
	}

	err := s.db.Register(data.(models.Auth))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("%v", err)
		return
	}

	s.cache.Delete(login)
	c.JSON(http.StatusOK, "user was veryfied")
}

func (s authService) Login(c *gin.Context) {
	var a models.Auth

	err := c.BindJSON(&a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Printf("error:%v", err.Error())
		return
	}

	user, err := s.db.GetUser(a)
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
