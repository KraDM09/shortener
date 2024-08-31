package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router interface {
	Post(pattern string, fn http.HandlerFunc)
	Get(pattern string, fn http.HandlerFunc)
	Delete(pattern string, fn http.HandlerFunc)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
	Use(middlewares ...func(http.Handler) http.Handler)
	Group(fn func(r chi.Router))
}
