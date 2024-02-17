package handlers

import "net/http"

func PingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("pong\n"))
	if err != nil {
		panic(err)
	}
}
