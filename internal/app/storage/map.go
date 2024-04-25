package storage

type MapStorage struct{}

var mapHashes = make(map[string]string)

func (m MapStorage) Save(hash string, url string) {
	mapHashes[hash] = url
}

func (m MapStorage) Get(hash string) string {
	url := mapHashes[hash]
	return url
}
