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

	rows, err := conn.Query(context.Background(), "SELECT id, title, author FROM articles")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Article
		err := rows.Scan(&a.Id, &a.Title, &a.Author)
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

	row := conn.QueryRow(context.Background(), "SELECT * FROM articles WHERE id = $1", id)

	err = row.Scan(&a.Id, &a.Title, &a.Content, &a.Author)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to find article")
	}

	return a, nil
}

func (db DB) GetByAuthor(id int) ([]models.Article, error) {
	var a []models.Article
	var rowId int
	var article models.Article

	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), "SELECT id FROM users_articles WHERE author_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from users_articles: %v", err)
	}
	defer rows.Close()

	// for rows.Next() {
	// 	err = rows.Scan(&rowId)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to scan row from users_articles: %v", err)
	// 	}

	// 	articleRows, err := conn.Query(context.Background(), "SELECT id, title, author FROM articles WHERE id = $1", id)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to retrieve data from articles: %v", err)
	// 	}
	// 	defer articleRows.Close()

	// 	for articleRows.Next() {
	// 		err = articleRows.Scan(&article.Id, &article.Title, &article.Author)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("unable to scan row from articles: %v", err)
	// 		}

	// 		a = append(a, article)
	// 	}
	// }

	return a, nil
}

func (db DB) CreateArticle(a models.Article) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "INSERT INTO articles(title, content, author) VALUES ($1, $2, $3)", a.Title, a.Content, a.Author)
	if err != nil {
		return fmt.Errorf("unable to insert: %v", err)
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
