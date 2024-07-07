package handlers

import (
	"context"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/jackc/pgx/v5"
)

func PingHandler(rw http.ResponseWriter, _ *http.Request) {
	conn, err := pgx.Connect(context.Background(), config.FlagDatabaseDsn)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	rw.WriteHeader(http.StatusOK)
}
