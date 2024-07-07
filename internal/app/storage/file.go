package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/util"
)

type FileStorage struct{}

type FileRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

func (s FileStorage) Save(hash string, url string, userId string) (string, error) {
	// сериализуем структуру в JSON формат
	data, err := json.Marshal(FileRow{
		UUID:        util.CreateUUID(),
		ShortURL:    hash,
		OriginalURL: url,
		UserID:      userId,
	})
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(config.FlagFileStoragePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return "", err
	}

	defer file.Close()

	data = append(data, '\n')

	if _, err := file.Write(data); err != nil {
		return "", err
	}

	return hash, nil
}

func (s FileStorage) Get(hash string) (string, error) {
	var url string

	file, err := os.Open(config.FlagFileStoragePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Создаем сканер для файла
	scanner := bufio.NewScanner(file)

	row := FileRow{}

	// Читаем файл построчно
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			return "", err
		}

		if row.ShortURL == hash {
			url = row.OriginalURL
			break
		}

	}

	return url, nil
}

func (s FileStorage) SaveBatch(batch []URL, userID string) error {
	file, err := os.OpenFile(config.FlagFileStoragePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, record := range batch {
		data, err := json.Marshal(FileRow{
			UUID:        util.CreateUUID(),
			ShortURL:    record.Short,
			OriginalURL: record.Original,
			UserID:      userID,
		})
		if err != nil {
			return err
		}

		data = append(data, '\n')

		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	return nil
}

func (s FileStorage) GetUrlsByUserID(userID string) (*[]URL, error) {
	return nil, fmt.Errorf("not implemented")
}
