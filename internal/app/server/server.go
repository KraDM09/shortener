package server

import (
	"net/http"

	"github.com/KraDM09/shortener/internal/app/compressor"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/handlers"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/storage"
)

func Run(store storage.Storage, r router.Router, logger logger.Logger, compressor compressor.Compressor) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	r.Use(logger.RequestLogger)
	r.Use(compressor.RequestCompressor)

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.SaveNewURLHandler(rw, r, store)
	})
	r.Get("/ping", handlers.PingHandler)
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetURLByHashHandler(rw, r, store)
	})
	r.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
		handlers.ShortenHandler(rw, r, store)
	})
	r.Post("/api/shorten/batch", func(rw http.ResponseWriter, r *http.Request) {
		handlers.BatchHandler(rw, r, store)
	})

	logger.Info("Running server", "address", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, r)
}
