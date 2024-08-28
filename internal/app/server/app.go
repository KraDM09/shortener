package server

import (
	"context"
	"net/http"
	"time"

	"github.com/KraDM09/shortener/internal/app/access"
	"github.com/KraDM09/shortener/internal/app/compressor"
	"github.com/KraDM09/shortener/internal/app/handlers"
	"github.com/KraDM09/shortener/internal/app/handlers/user"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/storage"
)

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	store      storage.Storage
	router     router.Router
	logger     logger.Logger
	compressor compressor.Compressor
	access     access.Access

	// канал для отложенного удаления URL
	hashChan chan storage.DeleteHash
}

// newApp принимает на вход внешние зависимости приложения и возвращает новый объект app
func newApp(
	ctx context.Context,
	store storage.Storage,
	router router.Router,
	logger logger.Logger,
	compressor compressor.Compressor,
	access access.Access,
) *app {
	instance := &app{
		store:      store,
		router:     router,
		logger:     logger,
		compressor: compressor,
		access:     access,
		hashChan:   make(chan storage.DeleteHash, 1024), // установим каналу буфер в 1024 строки
	}

	// запустим горутину с фоновым удалением хешей
	go instance.flushHashes(ctx)

	return instance
}

func (a *app) webhook(
	ctx context.Context,
) router.Router {
	a.router.Use(a.logger.RequestLogger)
	a.router.Use(a.compressor.RequestCompressor)
	a.router.Use(a.access.Request)

	handler := handlers.NewHandler(a.store)
	userHandler := user.NewHandler(a.store)

	a.router.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handler.SaveNewURLHandler(ctx, rw, r, GetUserID(r))
	})
	a.router.Get("/ping", func(rw http.ResponseWriter, r *http.Request) {
		handler.PingHandler(ctx, rw, r)
	})
	a.router.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handler.GetURLByHashHandler(ctx, rw, r)
	})
	a.router.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
		handler.ShortenHandler(ctx, rw, r, GetUserID(r))
	})
	a.router.Post("/api/shorten/batch", func(rw http.ResponseWriter, r *http.Request) {
		handler.BatchHandler(ctx, rw, r, GetUserID(r))
	})
	a.router.Get("/api/user/urls", func(rw http.ResponseWriter, r *http.Request) {
		userHandler.UrlsHandler(ctx, rw, r)
	})
	a.router.Delete("/api/user/urls", func(rw http.ResponseWriter, r *http.Request) {
		userHandler.DeleteUrlsHandler(ctx, rw, r, a.hashChan, GetUserID(r))
	})

	return a.router
}

// flushHashes постоянно удаляет несколько хешей из хранилища с определённым интервалом
func (a *app) flushHashes(
	ctx context.Context,
) {
	// будем сохранять хеши, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var deleteHashes []storage.DeleteHash

	for {
		select {
		case hash := <-a.hashChan:
			// добавим сообщение в слайс для последующего сохранения
			deleteHashes = append(deleteHashes, hash)
		case <-ticker.C:
			// подождём, пока придёт хотя бы один хеш на удаление
			if len(deleteHashes) == 0 {
				continue
			}
			// удалим все пришедшие хеши одновременно
			err := a.store.DeleteUrls(ctx, deleteHashes...)
			if err != nil {
				a.logger.Error("cannot delete hashes", "error", err.Error())
				// не будем стирать сообщения, попробуем отправить их чуть позже
				continue
			}
			// сотрём успешно удаленные хеши
			deleteHashes = nil
		}
	}
}
