package storage

type MapStorage struct{}

var mapHashes = make(map[string]string)

func (m MapStorage) Save(hash string, url string) (string, error) {
	mapHashes[hash] = url

	return hash, nil
}

func (m MapStorage) Get(hash string) (string, error) {
	url := mapHashes[hash]
	return url, nil
}

func (m MapStorage) SaveBatch(batch []URL) error {
	for _, record := range batch {
		mapHashes[record.Short] = record.Original
	}

	return nil
}
