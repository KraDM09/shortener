package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/KraDM09/shortener/internal/app/storage"
)

func GetURLByHashHandler(rw http.ResponseWriter, r *http.Request, store storage.Storage) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		http.Error(rw, "Не удалось распарсить адрес", http.StatusBadRequest)
		return
	}

	id := strings.TrimLeft(parsedURL.Path, "/")
	URL, err := store.Get(id)
	if err != nil {
		http.Error(rw, "Не удалось получить адрес", http.StatusBadRequest)
		return
	}

	if URL != "" {
		rw.Header().Set("Location", URL)
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
