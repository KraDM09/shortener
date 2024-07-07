package user

import (
	"encoding/json"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/constants"

	"github.com/KraDM09/shortener/internal/app/storage"
)

func UrlsHandler(rw http.ResponseWriter, r *http.Request, store storage.Storage) {
	userID := r.Context().Value(constants.ContextUserIDKey).(string)
	URLs, err := store.GetUrlsByUserID(userID)
	if err != nil {
		http.Error(rw, "Не удалось получить список адресов", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if len(*URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	for i, URL := range *URLs {
		(*URLs)[i].Short = config.FlagBaseShortURL + "/" + URL.Short
	}

	json.NewEncoder(rw).Encode(URLs)
}
