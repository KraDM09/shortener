package main

import (
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/router"
	"github.com/KraDM09/shortener/internal/app/server"
	"github.com/KraDM09/shortener/internal/app/storage"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	// обрабатываем аргументы командной строки
	config.ParseFlags()

	//store := &storage.SliceStorage{}
	store := &storage.MapStorage{}

	r := &router.ChiRouter{}

	if err := server.Run(store, r); err != nil {
		panic(err)
	}
}
