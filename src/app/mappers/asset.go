package mappers

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/models"
	base_dtos "github.com/pablor21/goms/pkg/dtos"
	"github.com/pablor21/goms/pkg/models/datatypes"
)

var _assetMapper *AssetMapper

type AssetMapper struct{}

func GetAssetMapper() *AssetMapper {
	if _assetMapper == nil {
		_assetMapper = &AssetMapper{}
	}
	return _assetMapper
}

func (m *AssetMapper) MapAssetInputToModel(ctx context.Context, input dtos.AssetInput) *models.Asset {
	model := models.NewAsset()
	model.ID = input.ID
	model.AssetType = input.AssetType
	// m.Uri = input.Uri
	model.StorageName = input.StorageName
	model.Description = input.Description
	model.Title = input.Title
	model.OwnerID = input.OwnerId
	model.OwnerType = input.OwnerType
	model.Section = input.Section
	model.MimeType = input.File.Header.Get("Content-Type")
	model.Metadata = &datatypes.Metadata{
		"size":            input.File.Size,
		"mimeType":        input.File.Header.Get("Content-Type"),
		"displayFileName": input.File.Filename,
	}
	return model
}

func (m *AssetMapper) MapAssetToDTO(ctx context.Context, input *models.Asset) dtos.AssetDTO {
	if input == nil {
		return dtos.AssetDTO{}
	}
	return dtos.AssetDTO{
		BaseTimestampDTO: base_dtos.BaseTimestampDTO[int64]{
			BaseDTO:   base_dtos.BaseDTO[int64]{ID: input.ID},
			CreatedAt: input.CreatedAt,
			UpdatedAt: input.UpdatedAt,
		},
		Title:       input.Title,
		Description: input.Description,
		Uri:         input.Uri,
		Section:     input.Section,
		StorageName: input.StorageName,
		MimeType:    input.MimeType,
		AssetType:   input.AssetType,
		MetadataDTO: base_dtos.MetadataDTO{
			Metadata: input.Metadata,
		},
		Permalink:   input.GetPermalink(ctx),
		DownloadUrl: input.GetDownloadUrl(ctx),
		// DownloadUrl: urlGenerator("asset.download", *m.StorageName, m.Section, m.ID, m.Uri),
		// DetailsUrl:  urlGenerator("api.asset.get", m.ID),
	}
}
