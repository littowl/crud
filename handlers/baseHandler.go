package handlers

import "crud/db"

type BaseHandler struct {
	db *db.DB
}

// для того, чтобы иметь доступ к методам созданной структуры DB (GetAll и т.д.)
func NewBaseHandler(db *db.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}
