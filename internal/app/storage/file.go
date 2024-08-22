package storage

import (
	"bufio"
	"encoding/json"
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
	IsDeleted   bool   `json:"is_deleted"`
}

func (s FileStorage) Save(hash string, url string, userID string) (string, error) {
	// сериализуем структуру в JSON формат
	data, err := json.Marshal(FileRow{
		UUID:        util.CreateUUID(),
		ShortURL:    hash,
		OriginalURL: url,
		UserID:      userID,
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

func (s FileStorage) Get(hash string) (*URL, error) {
	var url URL

	file, err := os.Open(config.FlagFileStoragePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Создаем сканер для файла
	scanner := bufio.NewScanner(file)

	row := FileRow{}

	// Читаем файл построчно
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			return nil, err
		}

		if row.ShortURL == hash {
			url = URL{
				Short:     row.ShortURL,
				Original:  row.OriginalURL,
				IsDeleted: row.IsDeleted,
			}
			break
		}

	}

	return &url, nil
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
	URLs := make([]URL, 0)

	file, err := os.Open(config.FlagFileStoragePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	row := FileRow{}

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			return nil, err
		}

		if row.UserID == userID {
			URLs = append(URLs, URL{
				Short:    row.ShortURL,
				Original: row.OriginalURL,
			})
			break
		}

	}

	return &URLs, nil
}

func (s FileStorage) DeleteUrls(deleteHashes ...DeleteHash) error {
	file, err := os.Open(config.FlagFileStoragePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// временный файл для записи
	tempFile, err := os.CreateTemp("./", "tempFile")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)

	row := FileRow{}

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			return err
		}

		if s.Contains(&deleteHashes, row.ShortURL, row.UserID) && !row.IsDeleted {
			row.IsDeleted = true
		}

		updatedRow, err := json.Marshal(row)
		if err != nil {
			return err
		}

		_, err = writer.WriteString(string(updatedRow) + "\n")
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = tempFile.Close()
	if err != nil {
		return err
	}

	// Заменяем исходный файл временным файлом
	if err := os.Rename(tempFile.Name(), config.FlagFileStoragePath); err != nil {
		return err
	}

	return nil
}

func (s FileStorage) Contains(deleteHash *[]DeleteHash, hash string, userID string) bool {
	for _, url := range *deleteHash {
		if url.Short == hash && url.UserID == userID {
			return true
		}
	}
	return false
}

func (s FileStorage) GetQuantityUserShortUrls(
	userID string,
	shortUrls *[]string,
) (int, error) {
	quantity := 0

	file, err := os.Open(config.FlagFileStoragePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	row := FileRow{}

	deleteHash := make([]DeleteHash, 0, len(*shortUrls))
	for _, hash := range *shortUrls {
		deleteHash = append(deleteHash, DeleteHash{
			Short:  hash,
			UserID: userID,
		})
	}

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			return 0, err
		}

		if s.Contains(&deleteHash, row.ShortURL, row.UserID) && !row.IsDeleted {
			quantity++
		}
	}

	return quantity, nil
}
