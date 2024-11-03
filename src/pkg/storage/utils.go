package storage

import (
	"path/filepath"

	"github.com/pablor21/goms/pkg/uuid"
)

func UniqueFileName(name string) string {
	ext := filepath.Ext(name)
	return uuid.GenerateUUIDWithoutHyphens() + ext
}
