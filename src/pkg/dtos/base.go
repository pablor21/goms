package dtos

import "time"

type BaseDTO[IDType any] struct {
	ID IDType `json:"id"`
}

type BaseTimestampDTO[IDType any] struct {
	BaseDTO[IDType]
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
