package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KraDM09/shortener/internal/constants"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/models"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/KraDM09/shortener/internal/app/util"
)

func ShortenHandler(
	rw http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
) {
	var req models.ShortenRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash := util.CreateHash()
	userID := r.Context().Value(constants.ContextUserIDKey).(string)
	short, err := store.Save(hash, req.URL, userID)

	rw.Header().Set("Content-Type", "application/json")

	if errors.Is(err, storage.ErrConflict) {
		hash = short
		rw.WriteHeader(http.StatusConflict)
	}

	resp := models.ShortenResponse{
		Result: config.FlagBaseShortURL + "/" + hash,
	}

	rw.WriteHeader(http.StatusCreated)

	// сериализуем ответ сервера
	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
