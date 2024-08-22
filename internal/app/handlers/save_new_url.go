package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/KraDM09/shortener/internal/app/util"
)

func (h *Handler) SaveNewURLHandler(
	rw http.ResponseWriter,
	r *http.Request,
	userID string,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	URL := string(body)

	hash := util.CreateHash()
	short, err := (*h.store).Save(hash, URL, userID)

	switch {
	case errors.Is(err, storage.ErrConflict):
		hash = short
		rw.WriteHeader(http.StatusConflict)
	case err != nil:
		http.Error(rw, fmt.Sprintf("Не удалось сохранить URL= %s hash= %s err= %s", URL, hash, err), http.StatusInternalServerError)
		return
	case short == "":
		http.Error(rw, fmt.Sprintf("Не удалось сохранить URL= %s hash= %s", URL, hash), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(config.FlagBaseShortURL + "/" + hash))
	if err != nil {
		http.Error(rw, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}
}
