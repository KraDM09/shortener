package storage

type SliceStorage struct{}

type Link struct {
	Hash   string
	URL    string
	UserID string
}

var hashes []Link

func (s SliceStorage) Save(hash string, url string, userID string) (string, error) {
	hashes = append(hashes, Link{
		Hash:   hash,
		URL:    url,
		UserID: userID,
	})

	return hash, nil
}

func (s SliceStorage) Get(hash string) (string, error) {
	var url string

	for _, h := range hashes {
		if h.Hash == hash {
			url = h.URL
			break
		}
	}

	return url, nil
}

func (s SliceStorage) SaveBatch(batch []URL, userID string) error {
	for _, record := range batch {
		hashes = append(hashes, Link{
			Hash:   record.Short,
			URL:    record.Original,
			UserID: userID,
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
			UserID:   h.UserID,
		})
	}

	return &URLs, nil
}
