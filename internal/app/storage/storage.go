package storage

import "errors"

type URL struct {
	Short    string `json:"short_url" db:"short"`
	Original string `json:"original_url" db:"original"`
	UserID   string `json:"user_id" db:"user_id"`
}

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

type Storage interface {
	Save(hash string, url string, userID string) (string, error)
	Get(hash string) (string, error)
	SaveBatch(batch []URL, userID string) error
	GetUrlsByUserID(userID string) (*[]URL, error)
}
