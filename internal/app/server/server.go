package server

import (
	"net/http"

	"github.com/KraDM09/shortener/internal/constants"

	"github.com/KraDM09/shortener/internal/app/access"
	"github.com/KraDM09/shortener/internal/app/compressor"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/logger"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/storage"
)

func Run(
	store storage.Storage,
	r router.Router,
	logger logger.Logger,
	compressor compressor.Compressor,
	access access.Access,
) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	// создаём экземпляр приложения, передавая внешние зависимости
	instance := newApp(store, r, logger, compressor, access)

	instance.logger.Info("Running server", "address", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, instance.webhook())
}

func GetUserID(r *http.Request) string {
	return r.Context().Value(constants.ContextUserIDKey).(string)
}
