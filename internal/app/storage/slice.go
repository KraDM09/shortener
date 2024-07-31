package storage

import (
	"context"
)

type SliceStorage struct{}

type Link struct {
	Hash      string
	URL       string
	UserID    string
	IsDeleted bool
}

var hashes []Link

func (s SliceStorage) Save(hash string, url string, userID string) (string, error) {
	hashes = append(hashes, Link{
		Hash:      hash,
		URL:       url,
		UserID:    userID,
		IsDeleted: false,
	})

	return hash, nil
}

func (s SliceStorage) Get(hash string) (*URL, error) {
	var url URL

	for _, h := range hashes {
		if h.Hash == hash {
			url = URL{
				Short:     h.Hash,
				Original:  h.URL,
				IsDeleted: h.IsDeleted,
			}
			break
		}
	}

	return &url, nil
}

func (s SliceStorage) SaveBatch(batch []URL, userID string) error {
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

func (s SliceStorage) GetUrlsByUserID(userID string) (*[]URL, error) {
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

func (s SliceStorage) DeleteUrls(_ context.Context, deleteHashes ...DeleteHash) error {
	for _, hash := range deleteHashes {
		link := s.Find(&hashes, hash.Short)

		if link == nil || link.IsDeleted {
			continue
		}

		if link.UserID == hash.UserID {
			link.IsDeleted = true
		}
	}

	return nil
}

func (s SliceStorage) Find(links *[]Link, hash string) *Link {
	for _, l := range *links {
		if l.Hash == hash {
			return &l
		}
	}
	return nil
}

func (s SliceStorage) GetQuantityUserShortUrls(
	userID string,
	shortUrls *[]string,
) (int, error) {
	quantity := 0

	for _, short := range *shortUrls {
		link := s.Find(&hashes, short)

		if link == nil || link.UserID == userID || link.IsDeleted {
			continue
		}

		quantity++
	}

	return quantity, nil
}
