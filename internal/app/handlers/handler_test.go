package handlers

import (
	"bytes"
	"github.com/KraDM09/shortener/internal/app/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handler(t *testing.T) {
	var endpoint string
	url := "https://practicum.yandex.ru/profile/go-advanced/"

	store := &storage.MapStorage{}

	t.Run("positive test #1", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewBufferString(url))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			SaveNewURLHandler(writer, request, store)
		})
		h(w, request)

		result := w.Result()
		result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "text/plain", result.Header.Get("Content-Type"))

		body, _ := io.ReadAll(result.Body)
		endpoint = string(body)
	})

	t.Run("positive test #2", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, endpoint, nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			GetURLByHashHandler(writer, request, store)
		})
		h(w, request)

		result := w.Result()
		result.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, url, result.Header.Get("Location"))
	})
}
