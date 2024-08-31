package storage

import (
	"context"
	"sync"
)

type SliceStorage struct{}

type Link struct {
	Hash      string
	URL       string
	UserID    string
	IsDeleted bool
}

var (
	hashes  []Link
	sliceMu sync.RWMutex
)

func (s SliceStorage) Save(
	_ context.Context,
	hash string,
	url string,
	userID string,
) (string, error) {
	sliceMu.RLock()
	defer sliceMu.RUnlock()

	hashes = append(hashes, Link{
		Hash:      hash,
		URL:       url,
		UserID:    userID,
		IsDeleted: false,
	})

	return hash, nil
}

func (s SliceStorage) Get(
	_ context.Context,
	hash string,
) (*URL, error) {
	sliceMu.RLock()
	defer sliceMu.RUnlock()

	var url URL

	for _, h := range hashes {
		if h.Hash == hash {
			url = URL{
				Short:     h.Hash,
				Original:  h.URL,
				UserID:    h.UserID,
				IsDeleted: h.IsDeleted,
			}
			break
		}
	}

	return &url, nil
}

func (s SliceStorage) SaveBatch(
	_ context.Context,
	batch []URL,
	userID string,
) error {
	sliceMu.Lock()
	defer sliceMu.Unlock()

	for _, record := range batch {
		hashes = append(hashes, Link{
			Hash:      record.Short,
			URL:       record.Original,
			UserID:    userID,
			IsDeleted: false,
		})
	}
	return nil
}

func (s SliceStorage) GetUrlsByUserID(
	_ context.Context,
	userID string,
) (*[]URL, error) {
	sliceMu.RLock()
	defer sliceMu.RUnlock()

	URLs := make([]URL, 0)

	for _, h := range hashes {
		if h.UserID != userID {
			continue
		}

		URLs = append(URLs, URL{
			Short:    h.Hash,
			Original: h.URL,
		})
	}

	return &URLs, nil
}

func (s SliceStorage) DeleteUrls(
	ctx context.Context,
	deleteHashes ...DeleteHash,
) error {
	sliceMu.Lock()
	defer sliceMu.Unlock()

	for _, hash := range deleteHashes {
		link, err := s.Get(ctx, hash.Short)
		if err != nil {
			return err
		}

		if link == nil || link.IsDeleted {
			continue
		}

		if link.UserID == hash.UserID {
			link.IsDeleted = true
		}
	}

	return nil
}

func (s SliceStorage) GetQuantityUserShortUrls(
	ctx context.Context,
	userID string,
	shortUrls *[]string,
) (int, error) {
	sliceMu.RLock()
	defer sliceMu.RUnlock()

	quantity := 0

	for _, short := range *shortUrls {
		link, err := s.Get(ctx, short)
		if err != nil {
			return 0, err
		}

		if link == nil || link.UserID == userID || link.IsDeleted {
			continue
		}

		quantity++
	}

	return quantity, nil
}
