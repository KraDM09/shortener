package handlers

import (
	"context"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) PingHandler(
	ctx context.Context,
	rw http.ResponseWriter,
	_ *http.Request,
) {
	conn, err := pgx.Connect(ctx, config.FlagDatabaseDsn)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)

	rw.WriteHeader(http.StatusOK)
}
