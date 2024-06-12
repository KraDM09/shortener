package main

import (
	"context"

	"github.com/KraDM09/shortener/internal/app/compressor"
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/server"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/jackc/pgx/v5"
)

func getStorage() storage.Storage {
	if len(config.FlagDatabaseDsn) > 0 {
		conn, err := pgx.Connect(context.Background(), config.FlagDatabaseDsn)
		if err != nil {
			panic(err)
		}

		pg := storage.PG{}.NewStore(conn)
		err = pg.Bootstrap(context.Background())
		if err != nil {
			panic(err)
		}

		return pg
	}

	if len(config.FlagFileStoragePath) > 0 {
		return &storage.FileStorage{}
	}

	return &storage.MapStorage{}
}

// функция main вызывается автоматически при запуске приложения
func main() {
	// обрабатываем аргументы командной строки
	config.ParseFlags()

	store := getStorage()
	r := &router.ChiRouter{}
	log := &logger.ZapLogger{}
	c := &compressor.GzipCompressor{}

	if err := server.Run(store, r, log, c); err != nil {
		panic(err)
	}
}
