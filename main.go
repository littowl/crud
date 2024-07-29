package main

import (
	"context"
	"crud/db"
	"crud/handlers"
	"crud/middleware"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lpernett/godotenv"
	"github.com/patrickmn/go-cache"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	m, err := migrate.New("file://migrations", databaseURL)

	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
		os.Exit(1)
	}

	cache := cache.New(5*time.Hour, 10*time.Hour)

	h := handlers.NewBaseHandler(db.NewDB(pool), cache)

	r := gin.Default()

	commonRoutes := r.Group("/articles")
	{
		commonRoutes.GET("", h.GetAll)
		commonRoutes.GET(":id", h.GetById)
		commonRoutes.GET("/author/:id", h.GetByAuthor)
	}

	routesWithAuth := r.Group("/articles").Use(middleware.Authenticate())
	{
		routesWithAuth.POST("", h.CreateArticle)
		routesWithAuth.PUT("", h.UpdateArticle)
		routesWithAuth.DELETE(":id", h.DeleteArticle)
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", h.Register)
		authRoutes.GET("/verify", h.Verify)
		authRoutes.POST("/login", h.Login)
	}

	r.Run()

}
