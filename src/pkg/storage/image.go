package storage

import (
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/pablor21/goms/pkg/storage/images"
	"github.com/spf13/afero"
)

func GetThumbnailFileName(params ThumbnailWithCacheParams) string {
	if params.Quality == 0 {
		params.Quality = 80
	}

	if params.Fit == "" {
		params.Fit = "contain"
	}
	// get the size
	size := fmt.Sprintf("%s_w%d_h%d_q%d", params.Fit, params.Width, params.Height, params.Quality)
	// if params.Fit == "contain" {
	// 	size = fmt.Sprintf("%d-%d_%d", params.Width, params.Height, params.Quality)
	// }

	// get the file name
	cacheFilename := fmt.Sprintf("%s/%s.%s", params.Path, size, params.Format)

	return cacheFilename
}

func GetThumbnailCacheFile(params ThumbnailWithCacheParams) (file afero.File, err error) {
	fileName := GetThumbnailFileName(params)
	cacheStorage := GetCacheStorage(params.StorageName)
	file, err = cacheStorage.Open(fileName)
	return
}

func GetThumnailCacheStats(params ThumbnailWithCacheParams) (fs.FileInfo, error) {
	fileName := GetThumbnailFileName(params)
	cacheStorage := GetCacheStorage(params.StorageName)
	info, err := cacheStorage.Stat(fileName)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func GetCacheStorage(storageName StorageName) afero.Fs {
	if storageName == "" {
		storageName = StorageNameDefaultCache
	}
	switch storageName {
	case StorageNamePublic:
		return GetStorage().GetBucket(StorageNamePublicCache)
	case StorageNamePrivate:
		return GetStorage().GetBucket(StorageNamePrivateCache)
	default:
		return GetStorage().GetBucket(StorageNameDefaultCache)
	}
}

func CacheFileExists(params ThumbnailWithCacheParams) (bool, error) {
	fileName := GetThumbnailFileName(params)
	cacheStorage := GetCacheStorage(params.StorageName)
	exists, err := afero.Exists(cacheStorage, fileName)
	return exists, err
}

func DeleteCache(storageName StorageName, uri string) error {
	cacheStorage := GetCacheStorage(storageName)
	// get the folder
	dir := path.Dir(uri)

	return cacheStorage.RemoveAll(dir)
}

type ThumbnailWithCacheParams struct {
	images.ThumbnailParams
	images.ImageExportParams
	StorageName StorageName
	IgnoreCache bool
	Path        string
}

type StorageImage interface {
	images.Image
	StorageName() StorageName
	GetBucket() afero.Fs
	GetCacheStorage(storageName StorageName) afero.Fs
	FileName() string
	ThumbnailWithCache(params ThumbnailWithCacheParams) (io.ReadSeekCloser, fs.FileInfo, error)
	CacheFileExists(params ThumbnailWithCacheParams) (bool, error)
	GetThumbnailCacheFile(params ThumbnailWithCacheParams) (file afero.File, err error)
	GetThumbnailFileName(params ThumbnailWithCacheParams) string
	GetThumnailCacheStats(params ThumbnailWithCacheParams) (fs.FileInfo, error)
	DeleteCache() error
}

type storageImageImpl struct {
	images.Image
	fileName    string
	storageName StorageName
	reader      io.ReadCloser
}

func NewStorageImage(fileName string, storageName StorageName) StorageImage {
	return &storageImageImpl{
		fileName:    fileName,
		storageName: storageName,
	}
}

func NewStorageImageFromReader(reader io.ReadCloser, fileName string, storageName StorageName) StorageImage {
	return &storageImageImpl{
		fileName:    fileName,
		storageName: storageName,
		reader:      reader,
	}
}

func (i *storageImageImpl) StorageName() StorageName {
	return i.storageName
}

func (i *storageImageImpl) FileName() string {
	return i.fileName

}

func (i *storageImageImpl) SetFileName(fileName string) {
	i.fileName = fileName
}

func (i *storageImageImpl) GetThumbnailFileName(params ThumbnailWithCacheParams) string {
	return GetThumbnailFileName(params)
}

func (i *storageImageImpl) GetBucket() afero.Fs {
	return GetStorage().GetBucket(i.storageName)
}

func (i *storageImageImpl) GetCacheStorage(storageName StorageName) afero.Fs {
	if storageName == "" {
		storageName = i.StorageName()
	}
	switch storageName {
	case StorageNamePublic:
		return GetStorage().GetBucket(StorageNamePublicCache)
	case StorageNamePrivate:
		return GetStorage().GetBucket(StorageNamePrivateCache)
	default:
		return GetStorage().GetBucket(StorageNameDefaultCache)
	}
}

func (i *storageImageImpl) CacheFileExists(params ThumbnailWithCacheParams) (bool, error) {
	return CacheFileExists(params)
}

func (i *storageImageImpl) Load() error {
	if i.Image != nil {
		return nil
	}
	if i.reader == nil {
		reader, err := i.GetBucket().Open(i.FileName())
		if err != nil {
			return err
		}
		i.reader = reader
	}

	image := images.NewImageFile(i.reader)
	err := image.Load()
	if err != nil {
		return err
	}
	i.Image = image
	return nil
}

func (i *storageImageImpl) Thumbnail(params images.ThumbnailParams) (images.Image, error) {
	err := i.Load()
	if err != nil {
		return nil, err
	}
	return i.Image.Thumbnail(params)
}

func (i *storageImageImpl) GetInfo() (images.ImageMetadata, error) {
	err := i.Load()
	if err != nil {
		return images.ImageMetadata{}, err
	}
	return i.Image.GetInfo()
}

func (i *storageImageImpl) Close() error {
	if i.Image != nil {
		i.Image.Close()
		i.Image = nil
	}
	if i.reader != nil {
		// i.file.Close()
		// i.file = nil
		i.reader.Close()
		i.reader = nil
	}
	return nil
}

func (i *storageImageImpl) GetThumbnailCacheFile(params ThumbnailWithCacheParams) (file afero.File, err error) {
	return GetThumbnailCacheFile(params)
}

func (i *storageImageImpl) DeleteCache() error {
	return DeleteCache(i.storageName, i.FileName())
}

func (i *storageImageImpl) ThumbnailWithCache(params ThumbnailWithCacheParams) (io.ReadSeekCloser, fs.FileInfo, error) {
	// check if the cache is ignored
	if !params.IgnoreCache {
		// get the cache file
		file, err := i.GetThumbnailCacheFile(params)
		if err == nil {
			info, err := file.Stat()
			if err != nil {
				return nil, nil, err
			}

			// return the file
			return file, info, nil
		}
		err = nil
	}

	// load the image
	err := i.Load()
	if err != nil {
		return nil, nil, err
	}
	// get the thumbnail
	_, err = i.Thumbnail(params.ThumbnailParams)
	if err != nil {
		return nil, nil, err
	}

	// export the thumbnail
	b, _, err := i.Export(params.ImageExportParams)
	if err != nil {
		return nil, nil, err
	}

	// save the thumbnail to cache
	cacheStorage := i.GetCacheStorage(params.StorageName)
	fileName := i.GetThumbnailFileName(params)

	cacheStorage.MkdirAll(params.Path, 0755)

	cacheFile, err := cacheStorage.Create(fileName)
	if err != nil {
		return nil, nil, err
	}

	_, err = cacheFile.Write(b)
	if err != nil {
		return nil, nil, err
	}

	// return the reader
	_, err = cacheFile.Seek(0, 0)
	if err != nil {
		return nil, nil, err
	}

	// stats
	info, err := cacheFile.Stat()
	if err != nil {
		return nil, nil, err
	}

	return cacheFile, info, nil

}

func (i *storageImageImpl) GetThumnailCacheStats(params ThumbnailWithCacheParams) (fs.FileInfo, error) {
	return GetThumnailCacheStats(params)
}
