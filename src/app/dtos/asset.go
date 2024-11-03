package dtos

import (
	"io"

	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/pkg/dtos"
	"github.com/pablor21/goms/pkg/interactions/request"
	"github.com/pablor21/goms/pkg/interactions/response"
	"github.com/pablor21/goms/pkg/storage"
	"github.com/pablor21/goms/pkg/storage/images"
)

type AssetDTO struct {
	dtos.BaseTimestampDTO[int64]
	dtos.MetadataDTO
	UniqueID    string              `json:"uniqueId,omitempty"`
	Uri         string              `json:"uri,omitempty"`
	Title       string              `json:"title,omitempty"`
	Section     string              `json:"section,omitempty"`
	Description string              `json:"description,omitempty"`
	StorageName storage.StorageName `json:"storageName,omitempty"`
	MimeType    string              `json:"mimeType,omitempty"`
	AssetType   models.AssetType    `json:"assetType,omitempty"`
	DownloadUrl string              `json:"downloadUrl,omitempty"`
	Permalink   string              `json:"permalink,omitempty"`
} // @name Asset

type AssetInput struct {
	request.Input
	dtos.MetadataDTO
	ID          int64               `json:"id" form:"id"`
	OwnerId     int64               `json:"ownerId" form:"ownerId"`
	OwnerType   string              `json:"ownerType" form:"ownerType"`
	File        *request.FileInput  `json:"file" form:"file"`
	Title       string              `json:"title" form:"title"`
	Description string              `json:"description" form:"description"`
	Section     string              `json:"section" form:"section"`
	AssetType   models.AssetType    `json:"assetType" form:"assetType"`
	StorageName storage.StorageName `json:"storageName" form:"storageName"`
} // @name AssetInput

type AssetShowParams struct {
	Width       int            `json:"width" form:"width" query:"w"`
	Height      int            `json:"height" form:"height" query:"h"`
	Quality     int            `json:"quality" form:"quality" query:"q"`
	Fit         images.FitType `json:"fit" form:"fit" query:"fit"` // cover (the image will be resized to cover the width and height specified, cropping the image if necessary), contain (the image will be resized to fit within the width and height specified, maintaining the aspect ratio)
	Format      string         `json:"format" form:"format" query:"f"`
	IgnoreCache bool           `json:"ignoreCache" form:"ignoreCache" query:"ignoreCache"`
} // @name AssetShowParams

type AssetDownloadInput struct {
	request.Input
	Id          int64               `json:"id" form:"id" query:"id" param:"id"`
	StorageName storage.StorageName `json:"storageName" form:"storageName" query:"storageName" param:"storageName"`
	UniqueId    string              `json:"uniqueId" form:"uniqueId" query:"uniqueId" param:"uniqueId"`
	FileName    string              `json:"fileName" form:"fileName" query:"fileName" param:"fileName"`
	Section     string              `json:"section" form:"section" query:"section" param:"section"`
	DisplayName string              `json:"displayName" form:"displayName" query:"displayName" param:"displayName"`
	AssetShowParams
} // @name AssetDownloadInput

type AssetDetailsResultData struct {
	AssetDTO
	Reader io.ReadSeekCloser `json:"-"`
} // @name AssetDetailsResultData

type AssetDetailsResult struct {
	response.TypedResponse[AssetDetailsResultData]
} // @name AssetDetailsResult

type AssetListByOwnerInput struct {
	request.PaginatedInput
	OwnerId   int64  `json:"ownerId" form:"ownerId" query:"ownerId" param:"ownerId"`
	OwnerType string `json:"ownerType" form:"ownerType" query:"ownerType" param:"ownerType"`
} // @name AssetListByOwnerInput

type AssetListByOwnerResultData struct {
	response.PaginatedResponseData[AssetDTO]
} // @name AssetListByOwnerResultData
