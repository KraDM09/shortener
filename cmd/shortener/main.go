package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/KraDM09/shortener/config"
	"github.com/go-chi/chi"
)

type Link struct {
	Hash string
	URL  string
}

var hashes = []Link{}

func createHash() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	hash := ""

	for i := 0; i < 6; i++ {
		randomNumber := rand.Intn(26)
		char := string(alphabet[randomNumber])

		if rand.Intn(2) == 1 {
			char = strings.ToUpper(char)
		}

		hash = hash + char
	}

	return hash
}

func handler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}

		url := string(body)

		hash := createHash()
		hashes = append(hashes, Link{Hash: hash, URL: url})

		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("http://" + config.FlagBaseShortURL + "/" + hash))
		return
	} else if r.Method == http.MethodGet {
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

	rw.WriteHeader(http.StatusBadRequest)
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
	r := chi.NewRouter()

	r.Post("/", handler)
	r.Get("/{id}", handler)

	fmt.Println("Running server on", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, r)
}
