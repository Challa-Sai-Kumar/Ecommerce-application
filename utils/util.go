package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewID() string {
	id := uuid.New().String()
	newID := strings.ReplaceAll(id, "-", "")
	return newID
}
