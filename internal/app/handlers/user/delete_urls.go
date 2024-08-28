package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/storage"
)

func (h *Handler) DeleteUrlsHandler(
	ctx context.Context,
	rw http.ResponseWriter,
	r *http.Request,
	hashChan chan storage.DeleteHash,
	userID string,
) {
	var urls []string
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

	quantity, err := (*h.store).GetQuantityUserShortUrls(ctx, userID, &urls)
	if err != nil {
		http.Error(rw, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	if quantity != len(urls) {
		http.Error(rw, "Не все URL принадлежат пользователю или существуют", http.StatusBadRequest)
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
