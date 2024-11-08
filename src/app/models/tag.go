package models

import "github.com/pablor21/goms/pkg/models"

type Tag struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	Name         string `json:"name,omitempty" gorm:"column:name"`
	Slug         string `json:"slug,omitempty" gorm:"column:slug"`
	CompleteSlug string `json:"completeSlug,omitempty" gorm:"column:complete_slug"`
	OwnerType    string `json:"ownerType,omitempty" gorm:"column:owner_type"`
	ParentID     *int64 `json:"parentId,omitempty" gorm:"column:parent_id"`
	Parent       *Tag   `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children     []*Tag `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
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
	TagID     int64  `json:"tagId,omitempty" gorm:"column:tag_id"`
	Tag       *Tag   `json:"tag,omitempty" gorm:"foreignKey:TagID;references:ID"`
	EntryID   int64  `json:"entryId,omitempty" gorm:"column:entry_id"`
	EntryType string `json:"entryType,omitempty" gorm:"column:entry_type"`
}

func NewTagEntry() *TagEntry {
	return &TagEntry{}
}

func (te *TagEntry) TableName() string {
	return "tag_entries"
}
