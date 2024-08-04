package access

import (
	"net/http"
)

type Access interface {
	Request(next http.Handler) http.Handler
}
