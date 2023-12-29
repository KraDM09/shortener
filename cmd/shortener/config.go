package main

import (
	"flag"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagBaseShortURL string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8888", "address and port to run server")

	flag.StringVar(&flagBaseShortURL, "b", "http://localhost:8000", "short url address")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
