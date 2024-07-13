package storage

type MapStorage struct{}

var (
	mapHashes  = make(map[string]string)
	userHashes = make(map[string][]URL)
)

func (m MapStorage) Save(hash string, url string, userID string) (string, error) {
	mapHashes[hash] = url
	userHashes[userID] = append(userHashes[userID], URL{
		Short:    hash,
		Original: url,
	})

	return hash, nil
}

func (m MapStorage) Get(hash string) (string, error) {
	url := mapHashes[hash]
	return url, nil
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
