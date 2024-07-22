package db

import (
	"context"
	"crud/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

// для того, чтобы предоставлять доступ к методам структуры DB
func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

func (db DB) GetAll() ([]models.Article, error) {
	var articles []models.Article

	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT articles.id, articles.title, users.* FROM articles JOIN users_articles ON id = article_id JOIN users ON author_id = users.user_id")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Article
		err := rows.Scan(&a.Id, &a.Title, &a.Author.Id, &a.Author.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (db DB) GetById(id int) (models.Article, error) {
	var a models.Article
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return models.Article{}, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(), "SELECT articles.*, users.* FROM articles JOIN users_articles ON id = article_id JOIN users ON author_id = users.user_id WHERE articles.id = $1", id)

	err = row.Scan(&a.Id, &a.Title, &a.Content, &a.Author.Id, &a.Author.Name)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to find article: %v", err)
	}

	return a, nil
}

func (db DB) GetByAuthor(id int) ([]models.Article, error) {
	var articles []models.Article

	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT articles.id, articles.title, users.* FROM articles JOIN users_articles ON id = article_id JOIN users ON author_id = users.user_id WHERE users.user_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from users_articles: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Article
		err := rows.Scan(&a.Id, &a.Title, &a.Author.Id, &a.Author.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		articles = append(articles, a)

	}

	return articles, nil
}

func (db DB) CreateArticle(a models.Article) error {
	var articleId int
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(), "INSERT INTO articles(title, content) VALUES ($1, $2) RETURNING id", a.Title, a.Content)
	err = row.Scan(&articleId)
	if err != nil {
		return fmt.Errorf("unable to insert: %v", err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO users_articles(author_id, article_id) VALUES ($1, $2)", a.Author.Id, articleId)
	if err != nil {
		return fmt.Errorf("unable to insert into mapping table: %v", err)
	}
	return nil
}

func (db DB) UpdateArticle(a models.Article) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	conn.QueryRow(context.Background(), "UPDATE articles SET title = $1, content = $2 WHERE id = $3", a.Title, a.Content, a.Id)

	return nil
}

func (db DB) DeleteArticle(id int) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "DELETE FROM articles WHERE id = $1 RETURNING id;", id)
	if err != nil {
		//если ты тупой, то тебе вернет ошибку пупсик
		return fmt.Errorf("unable to DELITE: %v", err)
	}

	return nil
}
