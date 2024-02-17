package handlers

import (
	"github.com/KraDM09/shortener/internal/app/storage"
	"net/http"
	"net/url"
	"strings"
)

func GetURLByHashHandler(rw http.ResponseWriter, r *http.Request, store storage.Storage) {
	parsedURL, err := url.Parse(r.RequestURI)

	if err != nil {
		http.Error(rw, "Не удалось распарсить адрес", http.StatusBadRequest)
		return
	}

	id := strings.TrimLeft(parsedURL.Path, "/")
	URL := store.Get(id)

	if URL != "" {
		rw.Header().Set("Location", URL)
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
