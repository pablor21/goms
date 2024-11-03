package services

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"os"
	"path"
	"strings"

	"github.com/pablor21/goms/app/config"
	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/mappers"
	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/errors"
	"github.com/pablor21/goms/pkg/storage"
	"github.com/pablor21/goms/pkg/storage/images"
	"github.com/pablor21/goms/pkg/storage/videos"
	"github.com/pablor21/goms/pkg/uuid"
	"gorm.io/gorm"
)

type ListByOwnerResult struct {
	Assets []*models.Asset
	First  int
	Last   int
	Total  int64
}

type AssetService interface {
	_SaveAsset(ctx context.Context, input dtos.AssetInput) (*models.Asset, error)
	DeleteAssetByID(ctx context.Context, id int64) error
	_DeleteAsset(ctx context.Context, model *models.Asset) error
	_GetAssetByID(ctx context.Context, id int64) (*models.Asset, error)
	_ListAssetsByOwnerID(ctx context.Context, input dtos.AssetListByOwnerInput) (ListByOwnerResult, error)
	// Gets the asset file reader and details
	GetAssetDetails(ctx context.Context, input dtos.AssetDownloadInput) (dtos.AssetDetailsResultData, error)
	// Gets the asset file reader (without querying the db)
	GetAssetFileReader(ctx context.Context, input dtos.AssetDownloadInput) (dtos.AssetDetailsResultData, error)
	PurgeAssetCache(ctx context.Context, id int64) error
}

var _assetService AssetService

func GetAssetService() AssetService {
	if _assetService == nil {
		_assetService = &assetService{
			db: database.GetConnection("default").(*database.GormConnection),
		}
	}
	return _assetService
}

type assetService struct {
	db *database.GormConnection
}

func (s *assetService) _SaveAsset(ctx context.Context, input dtos.AssetInput) (res *models.Asset, err error) {
	storageName := input.StorageName
	bucket := storage.GetStorage().GetBucket(storageName)
	dstFilename := path.Join(storage.UniqueFileName(input.File.Filename))

	model := mappers.GetAssetMapper().MapAssetInputToModel(ctx, input)
	model.Uri = dstFilename
	model.UniqueID = uuid.GenerateUUIDWithoutHyphens()

	var oldModel *models.Asset
	isNew := model.ID == 0
	if model.ID > 0 {
		err = s.db.Conn().Where("id = ?", model.ID).First(&oldModel).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return
		}
		if oldModel != nil {
			isNew = false
			// model.Uri = oldModel.Uri
			model.UniqueID = oldModel.UniqueID
		}
	}

	// Open the file
	org, err := input.File.Open()
	if err != nil {
		return
	}
	defer org.Close()

	if model.IsImage() {
		// Get the image info
		image := storage.NewStorageImageFromReader(org, model.GetFileName(), storageName)
		defer image.Close()
		info, err := image.GetInfo()
		if err != nil {
			return res, err
		}
		model.Metadata.Set("width", info.Width)
		model.Metadata.Set("height", info.Height)
		model.Metadata.Set("format", info.Format)
		org.Seek(0, 0)
	} else if model.IsVideo() {
		v := videos.NewFFmpegVideo()
		v.LoadFromReader(org)
		metadata, err := v.GetMetadata()
		if err != nil {
			return res, err
		}

		model.Metadata.Set("width", metadata.Width)
		model.Metadata.Set("height", metadata.Height)
		model.Metadata.Set("format", metadata.Format)
		model.Metadata.Set("codec", metadata.Codec)
		model.Metadata.Set("bitrate", metadata.Bitrate)
		model.Metadata.Set("duration", metadata.Duration)

	}

	// save to database
	tx, _, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()
	err = tx.Save(&model).Error
	if err != nil {
		return
	}

	dstPath := path.Dir(model.GetFileName())

	// if the file is not new, delete the old file
	if !isNew {
		err = bucket.RemoveAll(dstPath)
		if err != nil {
			return
		}
		// delete the cache
		storage.DeleteCache(storage.StorageName(oldModel.StorageName), dstPath)
	}

	// Create directory
	err = bucket.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return
	}

	// Open file
	dst, err := bucket.OpenFile(model.GetFileName(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, org)
	if err != nil {
		return
	}

	if model.IsVideo() {
		video := videos.NewFFmpegVideo()
		org.Seek(0, 0)
		video.LoadFromReader(org)
		poster, _ := video.MakePoster(videos.PosterParams{
			Time:   0,
			Frames: 3,
			Rate:   3,
			Format: "gif",
			PixFmt: "rgb24",
		})

		posterFilename := model.GetFileName() + "_poster.gif"
		model.Metadata.Set("poster_uri", model.Uri+"_poster.gif")

		// create a file from the bytes
		// Open file
		dst2, err := bucket.OpenFile(posterFilename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return res, err
		}
		defer dst2.Close()

		_, err = io.Copy(dst2, poster)
		if err != nil {
			return res, err
		}

		// Get the image info
		// image := storage.NewStorageImage(posterFilename, storageName)
		// defer image.Close()
		// info, err := image.GetInfo()
		// if err != nil {
		// 	return res, err
		// }

		// model.Metadata.Set("poster", map[string]interface{}{
		// 	"width":  info.Width,
		// 	"height": info.Height,
		// 	"format": info.Format,
		// 	"uri":    posterFilename,
		// })

		err = tx.Save(&model).Error
		if err != nil {
			return res, err
		}

		// save metadata

	}

	// save metadata
	metadatFile, err := model.OpenMetadataFile(os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer metadatFile.Close()

	metadata := mappers.GetAssetMapper().MapAssetToDTO(ctx, model)
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return
	}

	_, err = io.Copy(metadatFile, strings.NewReader(string(metadataBytes)))
	if err != nil {
		return
	}

	err = tx.Commit()

	// res = mappers.MapAssetModelToAssetDTO(ctx, model)

	res = model

	return

}

func (s *assetService) _ListAssetsByOwnerID(ctx context.Context, input dtos.AssetListByOwnerInput) (res ListByOwnerResult, err error) {
	var models []*models.Asset
	tx, _, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()

	total := int64(0)
	shouldCount := false

	query := tx.Where("owner_id = ? AND owner_type = ?", input.OwnerId, input.OwnerType)
	if input.Skip > 0 {
		query = query.Offset(input.Skip)
		shouldCount = true
	}

	if input.Take > 0 {
		query = query.Limit(input.Take)
		shouldCount = true
	}

	if input.SearchTerm != "" {
		query = query.Where("name LIKE ?", "%"+input.SearchTerm+"%")
	}

	if shouldCount {
		err = query.Model(&models).Count(&total).Error
		if err != nil {
			return
		}
	}

	err = query.Find(&models).Error
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	res.Assets = models
	res.Total = total
	res.First = input.Skip
	res.Last = input.Skip + input.Take
	if res.Last > int(res.Total) {
		res.Last = int(res.Total)
	}

	return
}

func (s *assetService) PurgeAssetCache(ctx context.Context, id int64) (err error) {
	var model *models.Asset
	err = s.db.Conn().Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}

	return model.DeleteCache()
}

func (s *assetService) _GetAssetByID(ctx context.Context, id int64) (res *models.Asset, err error) {
	var model *models.Asset
	err = s.db.Conn().Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}

	return model, nil
}

func (s *assetService) _DeleteAsset(ctx context.Context, model *models.Asset) (err error) {
	tx, _, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()

	err = tx.Delete(&model).Error
	if err != nil {
		return
	}

	go func() {
		model.DeleteFile()
		model.DeleteCache()
	}()

	err = tx.Commit()

	return
}

func (s *assetService) DeleteAssetByID(ctx context.Context, id int64) (err error) {
	var model *models.Asset
	err = s.db.Conn().Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}
	return s._DeleteAsset(ctx, model)
}

func (s *assetService) GetAssetDetails(ctx context.Context, input dtos.AssetDownloadInput) (res dtos.AssetDetailsResultData, err error) {
	var model *models.Asset
	err = s.db.Conn().Where("unique_id = ?", input.UniqueId).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}

	input.StorageName = model.StorageName
	input.FileName = model.Uri
	input.UniqueId = model.UniqueID
	input.Id = model.ID
	input.Section = model.Section

	readResult, err := s.GetAssetFileReader(ctx, input)
	if err != nil {
		return
	}

	res.Reader = readResult.Reader
	res.AssetDTO = mappers.GetAssetMapper().MapAssetToDTO(ctx, model)
	res.AssetDTO.MimeType = readResult.AssetDTO.MimeType

	return
}

func (s *assetService) GetAssetFileReader(ctx context.Context, input dtos.AssetDownloadInput) (res dtos.AssetDetailsResultData, err error) {

	if input.StorageName == "" || input.UniqueId == "" || input.FileName == "" || strings.HasPrefix(input.FileName, ".") {
		return res, errors.ErrBadRequest
	}

	model := models.NewAsset()
	model.ID = input.Id
	model.StorageName = input.StorageName
	model.Uri = input.FileName
	model.UniqueID = input.UniqueId
	model.MimeType = mime.TypeByExtension(path.Ext(input.FileName))
	model.Section = input.Section

	if model.IsImage() {
		filename := model.GetStoragePath()

		if input.Format == "" {
			input.Format = path.Ext(input.FileName)[1:] // remove the dot
		}

		if input.Quality == 0 {
			input.Quality = config.GetConfig().Assets.Image.DefaultQuality
		}

		if input.Quality < 0 || input.Quality > 100 {
			return res, errors.ErrBadRequest
		}

		if input.Fit == "" {
			input.Fit = "contain"
		}

		if input.Fit != "contain" && input.Fit != "cover" {
			return res, errors.ErrBadRequest
		}

		thumbParams := storage.ThumbnailWithCacheParams{
			ThumbnailParams: images.ThumbnailParams{
				Width:  input.Width,
				Height: input.Height,
				Fit:    input.Fit,
			},
			IgnoreCache: input.IgnoreCache,
			Path:        filename,
			ImageExportParams: images.ImageExportParams{
				Quality: input.Quality,
				Format:  input.Format,
			},
		}

		f, err := storage.GetThumbnailCacheFile(thumbParams)
		if err == nil && !input.IgnoreCache {
			stats, err := storage.GetThumnailCacheStats(thumbParams)
			if err != nil {
				if os.IsNotExist(err) {
					err = errors.ErrFileNotFound
				}
				return res, err
			}
			modTime := stats.ModTime()
			model.UpdatedAt = &modTime
			model.MimeType = mime.TypeByExtension(path.Ext(stats.Name()))

			res.AssetDTO = mappers.GetAssetMapper().MapAssetToDTO(ctx, model)
			res.Reader = f

		} else {

			reader, err := model.OpenFile(os.O_RDONLY, 0)
			if err != nil {
				if os.IsNotExist(err) {
					err = errors.ErrFileNotFound
				}
				return res, err
			}
			image := storage.NewStorageImageFromReader(reader, filename, model.StorageName)
			defer image.Close()
			//var thumb images.Image
			thumb, stats, err := image.ThumbnailWithCache(thumbParams)
			if err != nil {
				return res, err
			}
			modtime := stats.ModTime()
			model.UpdatedAt = &modtime
			model.MimeType = mime.TypeByExtension(path.Ext(stats.Name()))

			res.AssetDTO = mappers.GetAssetMapper().MapAssetToDTO(ctx, model)
			// check if the thumb cache exists

			// create a readseekcloser from the bytes
			res.Reader = thumb
		}

		return res, nil
	} else {

		reader, err := model.OpenFile(os.O_RDONLY, 0)
		if err != nil {
			return res, err
		}
		stat, err := model.GetFileStat()
		if err != nil {
			return res, err
		}

		modTime := stat.ModTime()
		model.UpdatedAt = &modTime
		res.AssetDTO = mappers.GetAssetMapper().MapAssetToDTO(ctx, model)
		res.Reader = reader
	}

	return
}
