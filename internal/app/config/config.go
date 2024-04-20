package config

import (
	"flag"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var (
	FlagRunAddr         string
	FlagBaseShortURL    string
	FlagLogLevel        string
	FlagFileStoragePath string
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию

	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagBaseShortURL, "b", "http://localhost:8080", "short url address")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&FlagFileStoragePath, "f", "/tmp/short-url-db.json", "file storage path")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		FlagRunAddr = serverAddress
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		FlagBaseShortURL = baseURL
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FlagFileStoragePath = envFileStoragePath
	}
}
