package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

var chiRouter = chi.NewRouter()

type ChiRouter struct{}

func (router ChiRouter) Post(pattern string, fn http.HandlerFunc) {
	chiRouter.Post(pattern, fn)
}

func (router ChiRouter) Delete(pattern string, fn http.HandlerFunc) {
	chiRouter.Delete(pattern, fn)
}

func (router ChiRouter) Get(pattern string, fn http.HandlerFunc) {
	chiRouter.Get(pattern, fn)
}

func (router ChiRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	chiRouter.ServeHTTP(rw, r)
}

func (router ChiRouter) Use(middlewares ...func(http.Handler) http.Handler) {
	chiRouter.Use(middlewares...)
}
