package mappers

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/models"
	base_dtos "github.com/pablor21/goms/pkg/dtos"
)

var _authMapper *AuthMapper

type AuthMapper struct{}

func GetAuthMapper() *AuthMapper {
	if _authMapper == nil {
		_authMapper = &AuthMapper{}
	}
	return _authMapper
}

func (m *AuthMapper) MapUserToDTO(ctx context.Context, input *models.User, includePrivateData bool) dtos.UserDTO {
	d := dtos.UserDTO{
		BaseTimestampDTO: base_dtos.BaseTimestampDTO[int64]{
			BaseDTO: base_dtos.BaseDTO[int64]{
				ID: input.ID,
			},
			CreatedAt: input.CreatedAt,
			UpdatedAt: input.UpdatedAt,
		},
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Lang:      input.Lang,
		Role:      input.Role,
		Status:    input.Status,
	}

	if includePrivateData {
		d.Email = input.Email
		d.PhoneNumber = input.PhoneNumber
	}

	if input.Avatar != nil && input.Avatar.ID > 0 {
		avatar := GetAssetMapper().MapAssetToDTO(ctx, input.Avatar)
		d.Avatar = &avatar
	}

	return d
}
