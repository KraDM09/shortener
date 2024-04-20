package storage

import (
	"bufio"
	"encoding/json"
	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/util"
	"log"
	"os"
)

type FileStorage struct {
}

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

type FileRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s FileStorage) Save(hash string, url string) {
	// сериализуем структуру в JSON формат
	data, err := json.Marshal(FileRow{
		UUID:        util.CreateUUID(),
		ShortURL:    hash,
		OriginalURL: url,
	})

	if err != nil {
		return
	}

	file, err := os.OpenFile(config.FlagFileStoragePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}

	defer file.Close()

	data = append(data, '\n')

	if _, err := file.Write(data); err != nil {
		return
	}
}

func (s FileStorage) Get(hash string) string {
	var url string

	file, err := os.Open(config.FlagFileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Создаем сканер для файла
	scanner := bufio.NewScanner(file)

	var row = FileRow{}

	// Читаем файл построчно
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &row)

		if err != nil {
			return url
		}

		if row.ShortURL == hash {
			url = row.OriginalURL
			break
		}

	}

	return url
}
