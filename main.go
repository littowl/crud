package main

import (
	"context"
	"crud/db"
	"crud/handlers"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lpernett/godotenv"
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

	h := handlers.NewBaseHandler(db.NewDB(pool))

	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("", h.GetAll)
		v1.GET(":id", h.GetById)
		v1.POST("", h.CreateArticle)
		v1.PUT("", h.UpdateArticle)
		v1.DELETE(":id", h.DeleteArticle)
	}

	r.Run()

}
