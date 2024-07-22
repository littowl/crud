package models

type User struct {
	Id   int    `json:"user_id" db:"user_id"`
	Name string `json:"user_name,omitempty" db:"user_name"`
}

type Article struct {
	Id      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content,omitempty" db:"content"`
	Author  User   `json:"author"`
}

type Auth struct {
	Id       int
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
	Hash     string
}
