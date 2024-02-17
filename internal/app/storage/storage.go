package storage

type Storage interface {
	Save(hash string, url string)
	Get(hash string) string
}
