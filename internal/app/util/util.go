package util

import (
	"math/rand"
	"strings"

	UUID "github.com/KraDM09/shortener/internal/app/util/uuid"
)

func CreateHash() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz1234567890"
	hash := ""

	for i := 0; i < 6; i++ {
		randomNumber := rand.Intn(36)
		char := string(alphabet[randomNumber])

		if rand.Intn(2) == 1 {
			char = strings.ToUpper(char)
		}

		hash = hash + char
	}

	return hash
}

func CreateUUID() string {
	uuid := &UUID.GoogleUUID{}

	return uuid.New()
}
