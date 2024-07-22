package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	Login string
	Uid   int
	Exp   int64
	jwt.MapClaims
}

func validateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		msg = err.Error()
		return &SignedDetails{}, msg
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return &SignedDetails{}, "error: "
	}
	fmt.Print("\n claims.exp: ", claims.Login, claims.Exp, " end.") // ошибка потому что в указатель на структуру не записались значения
	if claims.Exp < time.Now().Local().Unix() {
		return &SignedDetails{}, "token is expired"
	}

	return claims, msg
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Auth")
		if clientToken == "" {
			c.JSON(401, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}
		claims, err := validateToken(clientToken)
		fmt.Print(claims, claims.Login)
		if err != "" {
			if strings.Contains(err, "token is expired") {
				c.JSON(401, gin.H{"error": "Auth token is expired"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			fmt.Print("err after validate token: ", err)
			return
		}
		fmt.Print("\n no err \n")
		c.Set("Login", claims.Login)
		c.Set("Uid", claims.Uid)
		c.Next()

	}
}
