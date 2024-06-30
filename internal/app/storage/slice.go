package storage

type SliceStorage struct{}

type Link struct {
	Hash string
	URL  string
}

var hashes []Link

func (s SliceStorage) Save(hash string, url string) (string, error) {
	hashes = append(hashes, Link{Hash: hash, URL: url})

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

func (s SliceStorage) SaveBatch(batch []URL) error {
	for _, record := range batch {
		hashes = append(hashes, Link{
			Hash: record.Short,
			URL:  record.Original,
		})
	}
	return nil
}
