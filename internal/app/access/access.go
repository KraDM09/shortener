package access

import (
	"net/http"
)

type Access interface {
	SaveUserID(next http.Handler) http.Handler
	Control(next http.Handler) http.Handler
}
