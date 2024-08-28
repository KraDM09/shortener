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

//go:generate mockery --name=Storage
type Storage interface {
	Save(
		ctx context.Context,
		hash string,
		url string,
		userID string,
	) (string, error)

	Get(
		ctx context.Context,
		hash string,
	) (*URL, error)

	SaveBatch(
		ctx context.Context,
		batch []URL,
		userID string,
	) error

	GetUrlsByUserID(
		ctx context.Context,
		userID string,
	) (*[]URL, error)

	DeleteUrls(
		ctx context.Context,
		deleteHashes ...DeleteHash,
	) error

	GetQuantityUserShortUrls(
		ctx context.Context,
		userID string,
		shortUrls *[]string,
	) (int, error)
}
