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

func getStorage(
	ctx context.Context,
) (storage.Storage, error) {
	if len(config.FlagDatabaseDsn) > 0 {
		pg, err := storage.PG{}.NewStore(ctx)
		if err != nil {
			return nil, err
		}

		err = pg.Bootstrap(ctx)
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

	ctx := context.Background()
	store, err := getStorage(ctx)
	if err != nil {
		panic(fmt.Errorf("не удалось получить доступ к хранилищу %w", err))
	}

	r := &router.ChiRouter{}
	log := &logger.ZapLogger{}
	c := &compressor.GzipCompressor{}
	a := &access.Cookie{}

	if err := server.Run(ctx, store, r, log, c, a); err != nil {
		panic(fmt.Errorf("ошибка во время старта сервиса %w", err))
	}
}
