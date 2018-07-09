package helper

import (
	"log"

	"github.com/google/uuid"
)

func CreateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error: %v\n", err.Error())
	}
	return u.String()
}
