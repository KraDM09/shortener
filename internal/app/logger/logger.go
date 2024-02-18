package logger

import "net/http"

type Logger interface {
	RequestLogger(next http.Handler) http.Handler
	Initialize(level string) error
	Info(msg string, key string, value string)
}
