package storage

import (
	"context"
	"sync"
)

type MapStorage struct{}

var (
	mapMu      sync.RWMutex
	mapHashes  = make(map[string]string)
	userHashes = make(map[string][]URL)
)

func (m MapStorage) Save(
	_ context.Context,
	hash string,
	url string,
	userID string,
) (string, error) {
	mapMu.Lock()
	defer mapMu.Unlock()

	mapHashes[hash] = url
	userHashes[userID] = append(userHashes[userID], URL{
		Short:     hash,
		Original:  url,
		IsDeleted: false,
	})

	return hash, nil
}

func (m MapStorage) Get(
	_ context.Context,
	hash string,
) (*URL, error) {
	mapMu.RLock()
	defer mapMu.RUnlock()

	var url URL

	for _, userLinks := range userHashes {
		link := m.Find(&userLinks, hash)

		if link != nil {
			url = *link
			break
		}

	}
	return &url, nil
}

func (m MapStorage) SaveBatch(
	_ context.Context,
	batch []URL,
	userID string,
) error {
	mapMu.Lock()
	defer mapMu.Unlock()

	for _, record := range batch {
		mapHashes[record.Short] = record.Original
		userHashes[userID] = append(userHashes[userID], record)
	}

	return nil
}

func (m MapStorage) GetUrlsByUserID(
	_ context.Context,
	userID string,
) (*[]URL, error) {
	mapMu.RLock()
	defer mapMu.RUnlock()

	URLs := userHashes[userID]

	return &URLs, nil
}

func (m MapStorage) DeleteUrls(
	_ context.Context,
	deleteHashes ...DeleteHash,
) error {
	mapMu.Lock()
	defer mapMu.Unlock()

	for _, hash := range deleteHashes {
		userLinks := userHashes[hash.UserID]

		link := m.Find(&userLinks, hash.Short)

		if link == nil || link.IsDeleted {
			continue
		}

		link.IsDeleted = true
	}

	return nil
}

func (m MapStorage) Find(links *[]URL, hash string) *URL {
	for _, l := range *links {
		if l.Short == hash {
			return &l
		}
	}
	return nil
}

func (m MapStorage) GetQuantityUserShortUrls(
	_ context.Context,
	userID string,
	shortUrls *[]string,
) (int, error) {
	mapMu.RLock()
	defer mapMu.RUnlock()

	quantity := 0
	userLinks := userHashes[userID]

	for _, short := range *shortUrls {
		link := m.Find(&userLinks, short)

		if link == nil || link.IsDeleted {
			continue
		}

		quantity++
	}

	return quantity, nil
}
