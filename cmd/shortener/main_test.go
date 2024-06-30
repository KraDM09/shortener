package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/handlers"
	"github.com/KraDM09/shortener/internal/app/models"
	"github.com/KraDM09/shortener/internal/app/storage"

	"github.com/stretchr/testify/assert"
)

var (
	endpoint string
	url      = "https://practicum.yandex.ru/profile/go-advanced/"
	store    = &storage.MapStorage{}
)

func testGetURLByHash(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, endpoint, nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.GetURLByHashHandler(writer, request, store)
	})
	h(w, request)

	result := w.Result()
	result.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
	assert.Equal(t, url, result.Header.Get("Location"))
}

func Test_handler(t *testing.T) {
	t.Run("SaveNewURL", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, config.FlagBaseShortURL+"/", bytes.NewBufferString(url))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handlers.SaveNewURLHandler(writer, request, store)
		})
		h(w, request)

		result := w.Result()
		result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "text/plain", result.Header.Get("Content-Type"))

		body, _ := io.ReadAll(result.Body)
		endpoint = string(body)
	})

	t.Run("testGetURLByHash", testGetURLByHash)
}

func Test_handler2(t *testing.T) {
	t.Run("Shorten", func(t *testing.T) {
		req := models.ShortenRequest{
			URL: url,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			panic(fmt.Errorf("что-то пошло не так %w", err))
		}

		request := httptest.NewRequest(http.MethodPost, config.FlagBaseShortURL+"/api/shorten", bytes.NewBufferString(string(jsonData)))
		w := httptest.NewRecorder()
		h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handlers.ShortenHandler(writer, request, store)
		})
		h(w, request)

		result := w.Result()

		// Чтение тела ответа
		resultBody, err := io.ReadAll(result.Body)
		if err != nil {
			panic(fmt.Errorf("что-то пошло не так %w", err))
		}

		result.Body.Close()

		resp := models.ShortenResponse{}

		err = json.Unmarshal(resultBody, &resp)
		if err != nil {
			panic(fmt.Errorf("что-то пошло не так %w", err))
		}

		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

		endpoint = resp.Result
	})

	t.Run("testGetURLByHash", testGetURLByHash)
}
