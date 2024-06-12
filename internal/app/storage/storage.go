package storage

type URL struct {
	Short    string `json:"short_url" db:"short"`
	Original string `json:"original_url" db:"original"`
}

type Storage interface {
	Save(hash string, url string)
	Get(hash string) string
	SaveBatch(batch []URL) error
}
