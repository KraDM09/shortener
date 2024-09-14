package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handler) GetURLByHashHandler(
	ctx context.Context,
	rw http.ResponseWriter,
	r *http.Request,
) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		http.Error(rw, "Не удалось распарсить адрес", http.StatusBadRequest)
		return
	}

	id := strings.TrimLeft(parsedURL.Path, "/")
	URL, err := (*h.store).Get(ctx, id)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Не удалось получить адрес %s", err.Error()), http.StatusBadRequest)
		return
	}

	if URL == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if URL.IsDeleted {
		rw.WriteHeader(http.StatusGone)
		return
	}

	rw.Header().Set("Location", URL.Original)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
