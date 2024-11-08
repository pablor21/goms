package models

import (
	"time"

	"github.com/pablor21/goms/pkg/models/datatypes"
)

type BaseModel[IDType any] struct {
	ID IDType `json:"id"`
}

type BaseTimestampedModel[IDType any] struct {
	BaseModel[IDType]
	CreatedAt time.Time  `json:"createdAt,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
}

type SoftDeleteModel[IDType any] struct {
	BaseTimestampedModel[IDType]
	DeletedAt string `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

type MetadataModel struct {
	// Metadata *datatypes.JSON `json:"metadata" gorm:"metadata"`
	Metadata *datatypes.Metadata `json:"metadata,omitempty" gorm:"column:metadata"`
}

func (m *MetadataModel) MergeMetadata(metadata map[string]interface{}) {
	if m.Metadata == nil {
		return
	}
	for k, v := range metadata {
		m.Metadata.Set(k, v)
	}
}

type DisplayOrderModel struct {
	DisplayOrder int `json:"displayOrder,omitempty" gorm:"column:d_order,default:0,index"`
}
