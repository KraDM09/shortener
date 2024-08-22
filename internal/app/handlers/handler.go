package handlers

import (
	"github.com/KraDM09/shortener/internal/app/storage"
)

type Handler struct {
	store *storage.Storage
}

func NewHandler(
	store storage.Storage,
) *Handler {
	return &Handler{
		store: &store,
	}
}
