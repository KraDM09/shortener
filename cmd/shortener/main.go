package main

import (
	"github.com/KraDM09/shortener/internal/app/compressor"
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/server"
	"github.com/KraDM09/shortener/internal/app/storage"
)

func getStorage() storage.Storage {
	if len(config.FlagDatabaseDsn) > 0 {
		db := &storage.Database{}

		err := db.Migrate()
		if err != nil {
			panic(err)
		}

		return db
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
