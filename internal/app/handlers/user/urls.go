package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/constants"
)

func (h *Handler) UrlsHandler(
	ctx context.Context,
	rw http.ResponseWriter,
	r *http.Request,
) {
	rw.Header().Set("Content-Type", "application/json")
	value := r.Context().Value(constants.ContextUserIDKey)

	if value == nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := value.(string)
	URLs, err := (*h.store).GetUrlsByUserID(ctx, userID)
	if err != nil {
		http.Error(rw, "Не удалось получить список адресов", http.StatusInternalServerError)
		return
	}

	if len(*URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	for i, URL := range *URLs {
		(*URLs)[i].Short = config.FlagBaseShortURL + "/" + URL.Short
	}

	json.NewEncoder(rw).Encode(URLs)
}
