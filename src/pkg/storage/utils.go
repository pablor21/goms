package storage

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func UniqueFileName(name string) string {
	ext := filepath.Ext(name)
	return strings.ReplaceAll(uuid.New().String(), "-", "") + ext
}
