package base

import (
	"context"

	"github.com/pablor21/goms/pkg/interactions/response"
)

type CrudWriterUseCases[T any, IdT any, DtoT any, CreateT any, UpdateT any] interface {
	Create(ctx context.Context, input CreateT) (response.TypedResponse[DtoT], error)
	Update(ctx context.Context, id IdT, input UpdateT) (response.TypedResponse[DtoT], error)
	DeleteById(ctx context.Context, id IdT) (response.Response, error)
}

type CrudReaderUseCases[T any, IdT any, DtoT any] interface {
	FindById(ctx context.Context, id IdT) (response.TypedResponse[DtoT], error)
	FindAll(ctx context.Context) (response.TypedResponse[[]DtoT], error)
}

type CrudUseCases[T any, IdT any, DtoT any, CreateT any, UpdateT any] interface {
	CrudWriterUseCases[T, IdT, DtoT, CreateT, UpdateT]
	CrudReaderUseCases[T, IdT, DtoT]
}
