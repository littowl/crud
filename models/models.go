package models

type User struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	Author  User   `json:"author"`
}
