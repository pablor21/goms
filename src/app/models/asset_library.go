package models

import "github.com/pablor21/goms/pkg/models"

type AssetLibrary struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	RootFolderID *int64       `json:"rootFolderId" gorm:"column:root_folder_id"`
	RootFolder   *AssetFolder `json:"rootFolder" gorm:"foreignKey:RootFolderID;references:ID"`
}

type AssetFolder struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	LibraryID   int64         `json:"libraryId" gorm:"column:library_id"`
	Library     *AssetLibrary `json:"library" gorm:"foreignKey:LibraryID;references:ID"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ParentID    int64         `json:"parentId"`
	Parent      *AssetFolder  `json:"parent" gorm:"foreignKey:ParentID;references:ID"`
	Path        string        `json:"path" gorm:"column:path"`
	// Assets      []*Asset       `json:"assets" gorm:"gorm:polymorphicType:OwnerType;polymorphicId:OwnerID;polymorphicValue:media_folders"`
	Children []*AssetFolder `json:"children" gorm:"foreignKey:ParentID;references:ID"`
}
