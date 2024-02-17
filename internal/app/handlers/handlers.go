package handlers

import (
	"github.com/KraDM09/shortener/internal/app/storage"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/util"
)

type Handler struct {
	store storage.Storage
}

func PingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("pong\n"))
	if err != nil {
		panic(err)
	}
}

func SaveNewURLHandler(rw http.ResponseWriter, r *http.Request, store storage.Storage) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(rw, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	URL := string(body)

	hash := util.CreateHash()
	store.Save(hash, URL)

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write([]byte(config.FlagBaseShortURL + "/" + hash))
	if err != nil {
		panic(err)
	}
}

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
