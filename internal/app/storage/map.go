package storage

import "context"

type MapStorage struct{}

var (
	mapHashes  = make(map[string]string)
	userHashes = make(map[string][]URL)
)

func (m MapStorage) Save(hash string, url string, userID string) (string, error) {
	mapHashes[hash] = url
	userHashes[userID] = append(userHashes[userID], URL{
		Short:     hash,
		Original:  url,
		IsDeleted: false,
	})

	return hash, nil
}

func (m MapStorage) Get(hash string) (*URL, error) {
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

func (m MapStorage) SaveBatch(batch []URL, userID string) error {
	for _, record := range batch {
		mapHashes[record.Short] = record.Original
		userHashes[userID] = append(userHashes[userID], record)
	}

	return nil
}

func (m MapStorage) GetUrlsByUserID(userID string) (*[]URL, error) {
	URLs := userHashes[userID]

	return &URLs, nil
}

func (m MapStorage) DeleteUrls(_ context.Context, deleteHashes ...DeleteHash) error {
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
	userID string,
	shortUrls *[]string,
) (int, error) {
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
