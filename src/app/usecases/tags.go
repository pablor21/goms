package usecases

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/mappers"
	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/app/services"
	"github.com/pablor21/goms/app/usecases/base"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/interactions/response"
)

type TagUseCases interface {
	base.CrudUseCases[*models.Tag, int64, dtos.TagDTO, dtos.TagCreateInput, dtos.TagUpdateInput]
	FindByOwner(ctx context.Context, input dtos.TagListInput) (res dtos.TagListResponse, err error)
}

var _tagUseCases TagUseCases

func GetTagUseCases() TagUseCases {
	if _tagUseCases == nil {
		_tagUseCases = &tagUseCasesImpl{
			db: database.GetConnection("default").(*database.GormConnection),
		}
	}
	return _tagUseCases
}

// TagUseCasesImpl ...
type tagUseCasesImpl struct {
	db *database.GormConnection
}

func (u *tagUseCasesImpl) Create(ctx context.Context, input dtos.TagCreateInput) (res response.TypedResponse[dtos.TagDTO], err error) {
	m := mappers.GetTagMapper()
	model, err := services.GetTagService().Create(ctx, input)
	res = response.NewTypedResponse(m.MapTagToDTO(ctx, model))
	return
}

func (u *tagUseCasesImpl) Update(ctx context.Context, id int64, input dtos.TagUpdateInput) (res response.TypedResponse[dtos.TagDTO], err error) {
	m := mappers.GetTagMapper()
	model, err := services.GetTagService().Update(ctx, id, input)
	res = response.NewTypedResponse(m.MapTagToDTO(ctx, model))
	return
}

func (u *tagUseCasesImpl) DeleteById(ctx context.Context, id int64) (res response.Response, err error) {
	_, err = services.GetTagService().Delete(ctx, id)
	if err != nil {
		return
	}
	res = response.NewResponse().SetCode(204)
	return
}

func (u *tagUseCasesImpl) FindById(ctx context.Context, id int64) (res response.TypedResponse[dtos.TagDTO], err error) {
	m := mappers.GetTagMapper()
	model, err := services.GetTagService().GetByID(ctx, id)
	res = response.NewTypedResponse(m.MapTagToDTO(ctx, model))
	return
}

func (u *tagUseCasesImpl) FindByParentId(ctx context.Context, parentId int64) (res response.TypedResponse[dtos.TagDTO], err error) {
	m := mappers.GetTagMapper()
	model, err := services.GetTagService().GetByID(ctx, parentId)
	res = response.NewTypedResponse(m.MapTagToDTO(ctx, model))
	return
}

func (u *tagUseCasesImpl) FindAll(ctx context.Context) (res response.TypedResponse[[]dtos.TagDTO], err error) {
	m := mappers.GetTagMapper()
	models, err := services.GetTagService().FindAll(ctx, dtos.TagListInput{})
	res = response.NewTypedResponse(m.MapTagListToDTO(ctx, models.Result.Items))
	return
}

func (u *tagUseCasesImpl) FindByOwner(ctx context.Context, input dtos.TagListInput) (res dtos.TagListResponse, err error) {
	m := mappers.GetTagMapper()
	models, err := services.GetTagService().FindByOwner(ctx, input)
	res.Result.Items = m.MapTagListToDTO(ctx, models.Result.Items)
	res.Result.Total = models.Result.Total
	return
}
