package server

import (
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/handlers"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/storage"
	"net/http"
)

func Run(store storage.Storage, r router.Router, logger logger.Logger) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	r.Use(logger.RequestLogger)

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.SaveNewURLHandler(rw, r, store)
	})
	r.Post("/ping", handlers.PingHandler)
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetURLByHashHandler(rw, r, store)
	})

	logger.Info("Running server", "address", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, r)
}
