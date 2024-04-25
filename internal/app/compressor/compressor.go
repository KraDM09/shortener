package compressor

import "net/http"

type Compressor interface {
	RequestCompressor(h http.Handler) http.Handler
}
