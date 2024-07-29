package handlers

import (
	"crud/db"

	"github.com/patrickmn/go-cache"
)

type BaseHandler struct {
	db    *db.DB
	cache *cache.Cache
}

// для того, чтобы иметь доступ к методам созданной структуры DB (GetAll и т.д.)
func NewBaseHandler(db *db.DB, cache *cache.Cache) *BaseHandler {
	return &BaseHandler{
		db:    db,
		cache: cache,
	}
}
