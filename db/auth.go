package db

import (
	"context"
	"crud/models"
	"fmt"
)

func (db DB) Register(a models.Auth) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), "INSERT INTO auth(login, pass_hash) VALUES ($1, $2) RETURNING id", &a.Login, &a.Hash).Scan(&a.Id)
	if err != nil {
		return fmt.Errorf("unable to register new user: %v", err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO users(user_id, user_name) VALUES ($1, $2)", &a.Id, &a.Username)
	if err != nil {
		return fmt.Errorf("unable to insert new user: %v", err)
	}

	return nil
}

func (db DB) GetUser(a models.Auth) (models.Auth, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return models.Auth{}, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), "SELECT id, pass_hash FROM auth WHERE login = $1", &a.Login).Scan(&a.Id, &a.Hash)
	if err != nil {
		return models.Auth{}, fmt.Errorf("failed to find user with this login: %v", err)
	}

	return a, nil
}
