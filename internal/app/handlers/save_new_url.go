package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/KraDM09/shortener/internal/app/util"
)

func SaveNewURLHandler(
	rw http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
	userID string,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	URL := string(body)

	_, err = url.Parse(strings.TrimSpace(URL))
	if err != nil {
		http.Error(rw, fmt.Sprintf("Не удалось распарсить URL= %s err= %s", URL, err.Error()), http.StatusBadRequest)
		return
	}

	hash := util.CreateHash()
	short, err := store.Save(hash, URL, userID)

	switch {
	case errors.Is(err, storage.ErrConflict):
		hash = short
		rw.WriteHeader(http.StatusConflict)
	case err != nil:
		http.Error(rw, fmt.Sprintf("Не удалось сохранить URL= %s hash= %s err= %s", URL, hash, err.Error()), http.StatusInternalServerError)
		return
	case short == "":
		http.Error(rw, fmt.Sprintf("Не удалось сохранить URL= %s hash= %s", URL, hash), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(config.FlagBaseShortURL + "/" + hash))
	if err != nil {
		http.Error(rw, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}
}
