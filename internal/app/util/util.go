package util

import (
	UUID "github.com/KraDM09/shortener/internal/app/util/uuid"
	"math/rand"
	"strings"
)

func CreateHash() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	hash := ""

	for i := 0; i < 6; i++ {
		randomNumber := rand.Intn(26)
		char := string(alphabet[randomNumber])

		if rand.Intn(2) == 1 {
			char = strings.ToUpper(char)
		}

		hash = hash + char
	}

	return hash
}

func Uuid() string {
	uuid := &UUID.GoogleUUID{}

	return uuid.New()
}
