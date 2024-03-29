package router

import (
	"github.com/go-chi/chi"
	"net/http"
)

var chiRouter = chi.NewRouter()

type ChiRouter struct {
}

func (router ChiRouter) Post(pattern string, fn http.HandlerFunc) {
	chiRouter.Post(pattern, fn)
}

func (router ChiRouter) Get(pattern string, fn http.HandlerFunc) {
	chiRouter.Get(pattern, fn)
}

func (router ChiRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	chiRouter.ServeHTTP(rw, r)
}
