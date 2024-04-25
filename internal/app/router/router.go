package router

import "net/http"

type Router interface {
	Post(pattern string, fn http.HandlerFunc)
	Get(pattern string, fn http.HandlerFunc)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
	Use(middlewares ...func(http.Handler) http.Handler)
}
