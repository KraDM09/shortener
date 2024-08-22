package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KraDM09/shortener/internal/app/storage/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/KraDM09/shortener/internal/app/util"

	"github.com/KraDM09/shortener/internal/constants"
	"golang.org/x/net/context"

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
	userID   = "dabff768-c23d-4f8a-825d-7af2089ec901"
	handler  = handlers.NewHandler(store)
	ctx      = context.Background()
)

func testGetURLByHash(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, endpoint, nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.GetURLByHashHandler(writer, request)
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
			handler.SaveNewURLHandler(writer, request, util.CreateUUID())
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

		context := context.WithValue(ctx, constants.ContextUserIDKey, userID)

		request := httptest.NewRequest(http.MethodPost, config.FlagBaseShortURL+"/api/shorten", bytes.NewBufferString(string(jsonData)))
		request = request.WithContext(context)

		w := httptest.NewRecorder()
		h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler.ShortenHandler(writer, request, util.CreateUUID())
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

var shortURL string

func Test_save_new_url(t *testing.T) {
	t.Run("SaveNewUrl", func(t *testing.T) {
		storageProvider := new(mocks.Storage)
		storageProvider.
			On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything, nil)

		reqBody := bytes.NewBufferString(url)
		req := httptest.NewRequest(http.MethodPost, "/save", reqBody)
		rr := httptest.NewRecorder()

		h := handlers.NewHandler(storageProvider)
		h.SaveNewURLHandler(rr, req, "user-id")

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.NotEmpty(t, rr.Body.String(), "Короткий URL не должен быть пуст")

		shortURL = strings.TrimPrefix(rr.Body.String(), "/")
		assert.Equal(t, 6, len(shortURL), "Длина короткого URL должна быть 6 символов")
		storageProvider.AssertExpectations(t)
	})
}
