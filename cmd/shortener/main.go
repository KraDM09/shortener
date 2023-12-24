package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type Link struct {
	Hash string
	Url  string
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
			http.Error(rw, "Ошибка чтения тела запроса", http.StatusInternalServerError)
			return
		}

		url := string(body)

		hash := createHash()
		fmt.Print(hash)
		hashes = append(hashes, Link{Hash: hash, Url: url})

		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(hash))
		return
	} else if r.Method == http.MethodGet {
		path := r.RequestURI
		id := strings.TrimLeft(path, "/")
		var url string

		for _, hash := range hashes {
			if hash.Hash == id {
				url = hash.Url
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{id}", handler)
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
