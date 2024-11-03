package models

import "github.com/pablor21/goms/pkg/models"

type Tag struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	Name         string `json:"name" gorm:"column:name"`
	Slug         string `json:"slug" gorm:"column:slug"`
	CompleteSlug string `json:"complete_slug" gorm:"column:complete_slug"`
	OwnerType    string `json:"owner_type" gorm:"column:owner_type"`
	ParentID     int64  `json:"parent_id" gorm:"column:parent_id"`
	Parent       *Tag   `json:"parent" gorm:"foreignKey:ParentID;references:ID"`
	Children     []*Tag `json:"children" gorm:"foreignKey:ParentID;references:ID"`
}

func NewTag() *Tag {
	return &Tag{}
}

func (t *Tag) TableName() string {
	return "tags"
}

type TagEntry struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	TagID     int64  `json:"tag_id" gorm:"column:tag_id"`
	Tag       *Tag   `json:"tag" gorm:"foreignKey:TagID;references:ID"`
	EntryID   int64  `json:"entry_id" gorm:"column:entry_id"`
	EntryType string `json:"entry_type" gorm:"column:entry_type"`
}

func NewTagEntry() *TagEntry {
	return &TagEntry{}
}

func (te *TagEntry) TableName() string {
	return "tag_entries"
}
