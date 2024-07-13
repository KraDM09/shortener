package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/KraDM09/shortener/internal/constants"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/KraDM09/shortener/internal/app/util"
)

func SaveNewURLHandler(rw http.ResponseWriter, r *http.Request, store storage.Storage) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	URL := string(body)

	value := r.Context().Value(constants.ContextUserIDKey)

	if value == nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := value.(string)
	hash := util.CreateHash()
	short, err := store.Save(hash, URL, userID)

	switch {
	case errors.Is(err, storage.ErrConflict):
		hash = short
		rw.WriteHeader(http.StatusConflict)
	case err != nil:
		http.Error(rw, "Не удалось сохранить URL", http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(config.FlagBaseShortURL + "/" + hash))
	if err != nil {
		http.Error(rw, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}
}
