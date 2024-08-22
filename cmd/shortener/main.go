package main

import (
	"context"
	"fmt"

	"github.com/KraDM09/shortener/internal/app/access"

	"github.com/KraDM09/shortener/internal/app/compressor"
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/server"
	"github.com/KraDM09/shortener/internal/app/storage"
)

func getStorage() (storage.Storage, error) {
	if len(config.FlagDatabaseDsn) > 0 {
		pg, err := storage.PG{}.NewStore(context.Background())
		if err != nil {
			return nil, err
		}

		err = pg.Bootstrap()
		if err != nil {
			return nil, err
		}

		return pg, nil
	}

	if len(config.FlagFileStoragePath) > 0 {
		return &storage.FileStorage{}, nil
	}

	return &storage.MapStorage{}, nil
}

// функция main вызывается автоматически при запуске приложения
func main() {
	// обрабатываем аргументы командной строки
	config.ParseFlags()

	store, err := getStorage()
	if err != nil {
		panic(fmt.Errorf("не удалось получить доступ к хранилищу %w", err))
	}

	r := &router.ChiRouter{}
	log := &logger.ZapLogger{}
	c := &compressor.GzipCompressor{}
	a := &access.Cookie{}

	if err := server.Run(store, r, log, c, a); err != nil {
		panic(fmt.Errorf("ошибка во время старта сервиса %w", err))
	}
}
