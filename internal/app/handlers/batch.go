package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KraDM09/shortener/internal/app/config"

	"github.com/KraDM09/shortener/internal/app/models"
	"github.com/KraDM09/shortener/internal/app/storage"
	"github.com/KraDM09/shortener/internal/app/util"
)

type URL struct {
	CorrelationID string `json:"correlation_id"`
	Original      string `json:"original_url"`
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
			panic(err)
		}
		return
	}

	resp := make([]URL, len(req))
	batch := make([]storage.URL, len(req))

	for i, original := range req {
		hash := util.CreateHash()

		resp[i] = URL{
			CorrelationID: original.CorrelationID,
			Original:      original.URL,
			Short:         config.FlagBaseShortURL + "/" + hash,
		}

		batch[i] = storage.URL{
			Short:    hash,
			Original: original.URL,
		}
	}

	rw.Header().Set("Content-Type", "application/json")

	if err := store.SaveBatch(batch); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
