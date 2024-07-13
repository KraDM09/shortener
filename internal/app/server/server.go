package server

import (
	"net/http"

	"github.com/KraDM09/shortener/internal/constants"

	"github.com/KraDM09/shortener/internal/app/access"
	"github.com/KraDM09/shortener/internal/app/handlers/user"

	"github.com/KraDM09/shortener/internal/app/compressor"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/handlers"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/storage"
)

func Run(
	store storage.Storage,
	r router.Router,
	logger logger.Logger,
	compressor compressor.Compressor,
	access access.Access,
) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	r.Use(logger.RequestLogger)
	r.Use(compressor.RequestCompressor)
	r.Use(access.Request)

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.SaveNewURLHandler(rw, r, store, GetUserId(r))
	})
	r.Get("/ping", handlers.PingHandler)
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetURLByHashHandler(rw, r, store)
	})
	r.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
		handlers.ShortenHandler(rw, r, store, GetUserId(r))
	})
	r.Post("/api/shorten/batch", func(rw http.ResponseWriter, r *http.Request) {
		handlers.BatchHandler(rw, r, store, GetUserId(r))
	})
	r.Get("/api/user/urls", func(rw http.ResponseWriter, r *http.Request) {
		user.UrlsHandler(rw, r, store)
	})

	logger.Info("Running server", "address", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, r)
}

func GetUserId(r *http.Request) string {
	return r.Context().Value(constants.ContextUserIDKey).(string)
}
