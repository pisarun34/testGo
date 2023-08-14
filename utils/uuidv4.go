package utils

import (
	"github.com/google/uuid"
)

func GenerateUUIDv4() string {
	uuid := uuid.New()
	return uuid.String()
}
