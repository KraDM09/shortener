package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/KraDM09/shortener/config"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/util"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Link struct {
	Hash string
	URL  string
}

var hashes = []Link{}

// хендлер для /ping
func PingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("pong\n"))
}

func SaveNewURLHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(rw, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	url := string(body)

	hash := util.CreateHash()
	hashes = append(hashes, Link{Hash: hash, URL: url})

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(config.FlagBaseShortURL + "/" + hash))
}

func GetURLByHashHandler(rw http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)

	if err != nil {
		http.Error(rw, "Не удалось распарсить адрес", http.StatusBadRequest)
		return
	}

	id := strings.TrimLeft(parsedURL.Path, "/")
	var url string

	for _, hash := range hashes {
		if hash.Hash == id {
			url = hash.URL
			break
		}
	}

	if url != "" {
		rw.Header().Set("Location", url)
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

// функция main вызывается автоматически при запуске приложения
func main() {
	// обрабатываем аргументы командной строки
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Use(logger.RequestLogger)

	r.Post("/", SaveNewURLHandler)
	r.Post("/ping", PingHandler)
	r.Get("/{id}", GetURLByHashHandler)

	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))
	// оборачиваем хендлер webhook в middleware с логированием
	return http.ListenAndServe(config.FlagRunAddr, r)
}
