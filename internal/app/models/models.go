package models

import UUID "github.com/KraDM09/shortener/internal/app/util/uuid"

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type BatchRequest []Original

type Original struct {
	CorrelationId string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

type BatchResponse []Short

type Short struct {
	CorrelationId UUID.UUID `json:"correlation_id"`
	URL           string    `json:"short_url"`
}
