package uuid

import "github.com/google/uuid"

type GoogleUUID struct {
}

func (u GoogleUUID) New() string {
	return uuid.New().String()
}
