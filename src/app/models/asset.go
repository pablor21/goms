package models

import (
	"context"
	"io/fs"
	"mime"
	"path"
	"strings"

	"github.com/pablor21/goms/app/config"
	app_context "github.com/pablor21/goms/app/context"
	"github.com/pablor21/goms/pkg/errors"
	"github.com/pablor21/goms/pkg/models"
	"github.com/pablor21/goms/pkg/storage"
	"github.com/spf13/afero"
)

var ErrInvalidAssetEntryType = errors.NewAppError("Invalid asset entry type", 400)

type AssetType string //@Name AssetType

const (
	AssetTypeImage   AssetType = "IMAGE"         //@Name Image
	AssetTypeVideo   AssetType = "VIDEO"         //@Name Video
	AssetTypeYoutube AssetType = "YOUTUBE_VIDEO" //@Name YoutubeVideo
	AssetTypeAudio   AssetType = "AUDIO"         //@Name Audio
	AssetTypeFile    AssetType = "FILE"          //@Name File
)

func (t AssetType) String() string {
	return string(t)
}

func GetAssetTypes() []AssetType {
	return []AssetType{
		AssetTypeImage,
		AssetTypeVideo,
		AssetTypeYoutube,
		AssetTypeAudio,
		AssetTypeFile,
	}
}

func GetAssetType(assetType string) (AssetType, error) {
	assetType = strings.ToUpper(assetType)
	for _, t := range GetAssetTypes() {
		if t == AssetType(assetType) {
			return AssetType(strings.ToUpper(t.String())), nil
		}
	}
	return "", ErrInvalidAssetEntryType
}

type Asset struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	UniqueID    string              `json:"unique_id" gorm:"column:unique_id"`
	Title       string              `json:"title" gorm:"column:title"`
	Description string              `json:"description" gorm:"column:description"`
	Section     string              `json:"section" gorm:"column:section"`
	Uri         string              `json:"uri" gorm:"column:uri"`
	StorageName storage.StorageName `json:"storage_name" gorm:"column:storage_name"`
	MimeType    string              `json:"mime_type" gorm:"column:mime_type"`
	OwnerID     int64               `json:"owner_id" gorm:"column:owner_id"`
	OwnerType   string              `json:"owner_type" gorm:"column:owner_type"`
	AssetType   AssetType           `json:"asset_type" gorm:"column:asset_type"`
	AuthorId    int64               `json:"author_id" gorm:"column:author_id"`
	Author      *User               `json:"user" gorm:"foreignKey:AuthorId;references:ID"`
	TagEntries  []*TagEntry         `json:"tag_entries" gorm:"polymorphicType:EntryType;polymorphicID:EntryID;polymorphicValue:assets"`
}

func NewAsset() *Asset {
	return &Asset{}
}

func (a *Asset) TableName() string {
	return "assets"
}

func (m *Asset) GetStoragePath() string {
	return path.Join(config.GetConfig().Assets.BasePath, m.Section, m.UniqueID)
}

func (m *Asset) GetFileName() string {
	return path.Join(m.GetStoragePath(), m.Uri)
}

func (m *Asset) GetMetadataFileName() string {
	return path.Join(m.GetStoragePath(), "."+m.Uri+"_metadata.json")
}

func (m *Asset) OpenMetadataFile(flag int, mode fs.FileMode) (afero.File, error) {
	storageName := m.StorageName
	bucket := storage.GetStorage().GetBucket(storageName)
	reader, err := bucket.OpenFile(m.GetMetadataFileName(), flag, mode)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (m *Asset) DeleteCache() error {
	return storage.DeleteCache(m.StorageName, m.GetFileName())
}

func (m *Asset) DeleteFile() error {
	return storage.GetStorage().GetBucket(m.StorageName).RemoveAll(m.GetStoragePath())
}

func (m *Asset) OpenFile(flag int, mode fs.FileMode) (afero.File, error) {
	storageName := m.StorageName
	bucket := storage.GetStorage().GetBucket(storageName)
	reader, err := bucket.OpenFile(m.GetFileName(), flag, mode)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (m *Asset) GetFileStat() (fs.FileInfo, error) {
	storageName := m.StorageName
	bucket := storage.GetStorage().GetBucket(storageName)
	stat, err := bucket.Stat(m.GetFileName())
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (m *Asset) GetMimeType() string {
	if m.MimeType == "" {
		// try to get the mime type
		return mime.TypeByExtension(path.Ext(m.Uri))
	}
	return m.MimeType
}

func (m *Asset) IsImage() bool {
	switch m.MimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func (m *Asset) IsVideo() bool {
	switch m.MimeType {
	case "video/mp4", "video/webm", "video/ogg":
		return true
	default:
		return false
	}
}

func (m *Asset) IsAudio() bool {
	switch m.MimeType {
	case "audio/mpeg", "audio/ogg", "audio/wav":
		return true
	default:
		return false
	}
}

// GetPermalink returns the permalink of the asset, even if the uri changes the permalink will be the same
func (m *Asset) GetPermalink(ctx context.Context) string {
	urlGenerator := app_context.GetUrlGenerator(ctx)
	if urlGenerator == nil {
		return ""
	}
	return urlGenerator("assets.permalink", m.UniqueID)
}

func (m *Asset) GetDownloadUrl(ctx context.Context) string {
	urlGenerator := app_context.GetUrlGenerator(ctx)
	if urlGenerator == nil {
		return ""
	}
	return urlGenerator("assets.download", m.StorageName, m.Section, m.UniqueID, m.Uri)
}
