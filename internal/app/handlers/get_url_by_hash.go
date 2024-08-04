package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/KraDM09/shortener/internal/app/storage"
)

func GetURLByHashHandler(
	rw http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
) {
	parsedURL, err := url.Parse(strings.TrimSpace(r.RequestURI))
	if err != nil {
		http.Error(rw, "Не удалось распарсить адрес", http.StatusBadRequest)
		return
	}

	id := strings.TrimLeft(parsedURL.Path, "/")
	URL, err := store.Get(id)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Не удалось получить адрес %s", err.Error()), http.StatusBadRequest)
		return
	}

	if URL == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if URL.IsDeleted {
		rw.WriteHeader(http.StatusGone)
		return
	}

	rw.Header().Set("Location", URL.Original)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
