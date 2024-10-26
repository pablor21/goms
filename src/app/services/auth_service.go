package services

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
)

type AuthService interface {
	Login(ctx context.Context, req dtos.UserLoginInput) (dtos.UserLoginResponse, error)
	AuthenticateContext(ctx context.Context, token string) (dtos.UserDTO, context.Context, error)
	GetContextUser(ctx context.Context) (dtos.UserDTO, error)
}
