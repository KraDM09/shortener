package server

import (
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/handlers"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/storage"
	"net/http"
)

func Run(store storage.Storage, r router.Router) error {
	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.SaveNewURLHandler(rw, r, store)
	})
	r.Post("/ping", handlers.PingHandler)
	r.Post("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetURLByHashHandler(rw, r, store)
	})

	return http.ListenAndServe(config.FlagRunAddr, r)
}
