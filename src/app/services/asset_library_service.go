package services

import (
	"context"

	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/pkg/database"
)

type AssetLibraryService interface {
	SaveLibrary(ctx context.Context, input *models.AssetLibrary) (*models.AssetLibrary, error)
	_SaveLibrary(ctx context.Context, input *models.AssetLibrary) (*models.AssetLibrary, error)
}

var _assetLibraryService AssetLibraryService

func GetAssetLibraryService() AssetLibraryService {
	if _assetLibraryService == nil {
		_assetLibraryService = &assetLibraryService{
			db: database.GetConnection("default").(*database.GormConnection),
		}
	}
	return _assetLibraryService
}

type assetLibraryService struct {
	db *database.GormConnection
}

func (s *assetLibraryService) SaveLibrary(ctx context.Context, input *models.AssetLibrary) (res *models.AssetLibrary, err error) {
	return
}

func (s *assetLibraryService) _SaveLibrary(ctx context.Context, input *models.AssetLibrary) (res *models.AssetLibrary, err error) {
	return
}
