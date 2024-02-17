package storage

type SliceStorage struct {
}

type Link struct {
	Hash string
	URL  string
}

var hashes []Link

func (s SliceStorage) Save(hash string, url string) {
	hashes = append(hashes, Link{Hash: hash, URL: url})
}

func (s SliceStorage) Get(hash string) string {
	var url string

	for _, h := range hashes {
		if h.Hash == hash {
			url = h.URL
			break
		}
	}

	return url
}
