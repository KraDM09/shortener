package user

import (
	"encoding/json"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/models"
	"github.com/KraDM09/shortener/internal/app/storage"
)

func DeleteUrlsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	store storage.Storage,
	hashChan chan storage.DeleteHash,
	userID string,
) {
	var urls models.DeleteUrlsRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&urls); err != nil {
		http.Error(rw, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		http.Error(rw, "Пустные батчи запрещены", http.StatusBadRequest)
		return
	}

	shortUrls := make([]string, 0, len(urls))

	for _, hash := range urls {
		shortUrls = append(shortUrls, hash)
	}

	quantity, err := store.GetQuantityUserShortUrls(userID, &shortUrls)

	if quantity != len(urls) {
		http.Error(rw, "Не все URL принадлежат пользователю или существуют", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	for i := range urls {
		// положим в очередь на удаление
		hashChan <- storage.DeleteHash{
			UserID: userID,
			Short:  urls[i],
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
}
