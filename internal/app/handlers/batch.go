package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"

	"github.com/KraDM09/shortener/internal/app/models"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/KraDM09/shortener/internal/app/util"
)

type URL struct {
	CorrelationID string `json:"correlation_id"`
	Short         string `json:"short_url"`
}

func BatchHandler(rw http.ResponseWriter, r *http.Request, store storage.Storage) {
	var req models.BatchRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(req) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		_, err := rw.Write([]byte("Пустные батчи запрещены"))
		if err != nil {
			panic(fmt.Errorf("что-то пошло не так %w", err))
		}
		return
	}

	resp := make([]URL, 0, len(req))
	batch := make([]storage.URL, 0, len(req))

	for i := range req {
		hash := util.CreateHash()

		resp = append(resp, URL{
			CorrelationID: req[i].CorrelationID,
			Short:         config.FlagBaseShortURL + "/" + hash,
		})

		batch = append(batch, storage.URL{
			Short:    hash,
			Original: req[i].URL,
		})
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := store.SaveBatch(batch); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
