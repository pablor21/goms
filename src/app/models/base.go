package models

type BaseModel[IDType any] struct {
	ID IDType `json:"id"`
}

type BaseTimestampModel[IDType any] struct {
	BaseModel[IDType]
	CreatedAt string `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type SoftDeleteModel[IDType any] struct {
	BaseTimestampModel[IDType]
	DeletedAt string `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}
