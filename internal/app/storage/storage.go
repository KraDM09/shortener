package storage

import (
	"context"
	"errors"
)

type URL struct {
	Short     string `json:"short_url" db:"short"`
	Original  string `json:"original_url" db:"original"`
	IsDeleted bool   `json:"is_deleted" db:"is_deleted"`
}

// DeleteHash описывает объект хеша на удаление
type DeleteHash struct {
	UserID string `json:"user_id"`
	Short  string `json:"short_url"`
}

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

type Storage interface {
	Save(
		hash string,
		url string,
		userID string,
	) (string, error)

	Get(
		hash string,
	) (*URL, error)

	SaveBatch(
		batch []URL,
		userID string,
	) error

	GetUrlsByUserID(
		userID string,
	) (*[]URL, error)

	DeleteUrls(
		ctx context.Context,
		deleteHashes ...DeleteHash,
	) error

	GetQuantityUserShortUrls(
		userID string,
		shortUrls *[]string,
	) (int, error)
}
