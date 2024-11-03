package uuid

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func GenerateUUIDWithoutHyphens() string {
	return strings.ReplaceAll(GenerateUUID().String(), "-", "")
}
