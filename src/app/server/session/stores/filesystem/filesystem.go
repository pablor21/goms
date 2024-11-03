package filesystem

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pablor21/goms/pkg/logger"
	"github.com/pablor21/goms/pkg/storage"
	"github.com/spf13/afero"
)

type FileSystemStore struct {
	sesManager *scs.SessionManager
	bucket     afero.Fs
	path       string
	storeName  string
	ticker     *time.Ticker
}

func New(config map[string]interface{}, sm *scs.SessionManager) *FileSystemStore {
	storeName, ok := config["storage_name"].(string)

	if !ok {
		panic("storage_name is required")
	}

	path, ok := config["path"].(string)
	if !ok {
		panic("path is required")
	}

	bucket := storage.GetStorage().GetBucket(storage.StorageName(storeName))

	if bucket == nil {
		panic(fmt.Sprintf("bucket %s not found", storeName))
	}

	if err := bucket.MkdirAll(path, 0755); err != nil {
		panic(err)
	}

	fs := &FileSystemStore{
		sesManager: sm,
		path:       path,
		storeName:  storeName,
		bucket:     bucket,
	}

	cleanupInterval, ok := config["cleanup_interval"].(int)

	if ok && cleanupInterval > 0 {
		fs.InitCleanupInterval(time.Duration(cleanupInterval) * time.Second)
	}

	return fs
}

func (fs *FileSystemStore) InitCleanupInterval(interval time.Duration) {
	if interval == 0 {
		return
	}

	if fs.ticker != nil {
		fs.ticker.Stop()
	}

	fs.ticker = time.NewTicker(interval)
	go func() {
		for range fs.ticker.C {
			fs.Cleanup()
		}
	}()
}

func (fs *FileSystemStore) StopCleanup() {
	logger.Debug().Msg("Stopping session cleanup")
	if fs.ticker != nil {
		fs.ticker.Stop()
	}
}

func (fs *FileSystemStore) Cleanup() {
	logger.Debug().Msg("Running session cleanup")

	files, err := afero.ReadDir(fs.bucket, fs.path)
	if err != nil {
		logger.Error().Err(err).Msg("Error reading file storage directory")
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if time.Now().After(file.ModTime().Add(fs.sesManager.Lifetime)) {
			fs.Delete(file.Name())
		}
	}
}

func (fs *FileSystemStore) Delete(token string) (err error) {
	return fs.bucket.Remove(path.Join(fs.path, token))
}

func (fs *FileSystemStore) Find(token string) (b []byte, found bool, err error) {
	f, err := fs.bucket.Open(path.Join(fs.path, token))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	b, err = afero.ReadAll(f)
	if err != nil {
		return nil, false, err
	}

	return b, true, nil
}

func (fs *FileSystemStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	return afero.WriteFile(fs.bucket, path.Join(fs.path, token), b, 0644)
}
