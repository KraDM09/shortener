package storage

import "errors"

type URL struct {
	Short    string `json:"short_url" db:"short"`
	Original string `json:"original_url" db:"original"`
}

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

type Storage interface {
	Save(hash string, url string) (string, error)
	Get(hash string) string
	SaveBatch(batch []URL) error
}
