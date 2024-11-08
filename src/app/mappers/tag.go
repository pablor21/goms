package mappers

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/models"
)

var _tagMapper *TagMapper

type TagMapper struct{}

func GetTagMapper() *TagMapper {
	if _tagMapper == nil {
		_tagMapper = &TagMapper{}
	}
	return _tagMapper
}

func (m *TagMapper) MapTagToDTO(ctx context.Context, input *models.Tag) dtos.TagDTO {
	return dtos.TagDTO{
		Tag: *input,
	}
}

func (m *TagMapper) MapTagListToDTO(ctx context.Context, input []*models.Tag) []dtos.TagDTO {
	dtos := make([]dtos.TagDTO, 0)
	for _, tag := range input {
		dtos = append(dtos, m.MapTagToDTO(ctx, tag))
	}
	return dtos
}

func (m *TagMapper) MapTagCreateInputToModel(ctx context.Context, input dtos.TagCreateInput) *models.Tag {
	t := &models.Tag{
		Name:      input.Name,
		OwnerType: input.OwnerType,
		Slug:      input.Slug,
	}

	if input.ParentID != nil && *input.ParentID > 0 {
		t.ParentID = input.ParentID
	}

	return t
}

func (m *TagMapper) MapTagUpdateInputToModel(ctx context.Context, input dtos.TagUpdateInput) *models.Tag {
	return m.MapTagCreateInputToModel(ctx, input.TagCreateInput)
}
