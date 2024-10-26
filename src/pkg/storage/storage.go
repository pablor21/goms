package storage

import (
	"io/fs"
	"os"

	"github.com/pablor21/goms/pkg/logger"
	"github.com/pablor21/goms/pkg/storage/config"
	"github.com/spf13/afero"
)

var storage *Storage

type StorageName string

type File struct {
	afero.File
}

type FileInfo struct {
}

const (
	StorageNameDefault      StorageName = "default"
	StorageNameDefaultCache StorageName = "default_cache"
	StorageNamePublic       StorageName = "public"
	StorageNamePublicCache  StorageName = "public_cache"
	StorageNamePrivate      StorageName = "private"
	StorageNamePrivateCache StorageName = "private_cache"
	StorageNameTmp          StorageName = "tmp"
)

func (s StorageName) String() string {
	return string(s)
}

func ToStorageName(name string) StorageName {
	switch name {
	case "default":
		return StorageNameDefault
	case "default_cache":
		return StorageNameDefaultCache
	case "public":
		return StorageNamePublic
	case "public_cache":
		return StorageNamePublicCache
	case "private":
		return StorageNamePrivate
	case "private_cache":
		return StorageNamePrivateCache
	case "tmp":
		return StorageNameTmp
	default:
		return StorageNameDefault
	}
}

type Storage struct {
	Storages map[StorageName]afero.Fs
}

func NewStorage() *Storage {
	return &Storage{
		Storages: make(map[StorageName]afero.Fs),
	}
}

func InitStorage(cfg config.StorageConfig) *Storage {
	storage = NewStorage()
	for k, v := range cfg {
		logger.Info().Msgf("Adding bucket %s to storage", k)
		storage.AddBucket(ToStorageName(k), v.Uri)
	}
	return storage
}

func GetStorage() *Storage {
	return storage
}

func (s *Storage) AddBucket(name StorageName, f string) {
	err := os.MkdirAll(f, os.ModePerm)
	if err != nil {
		panic(err)
	}
	af := afero.NewBasePathFs(afero.NewOsFs(), f)

	s.Storages[name] = af
}

func (s *Storage) GetBucket(name StorageName) afero.Fs {
	return s.Storages[name]
}

func (s *Storage) ReadDir(fs afero.Fs, name string) ([]fs.FileInfo, error) {
	return afero.ReadDir(fs, name)
}
