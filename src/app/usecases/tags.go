package usecases

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/mappers"
	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/app/repositories"
	"github.com/pablor21/goms/app/usecases/base"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/interactions/response"
	"gorm.io/gorm"
)

type TagUseCases interface {
	base.CrudUseCases[*models.Tag, int64, dtos.TagDTO, dtos.TagCreateInput, dtos.TagUpdateInput]
	FindByOwner(ctx context.Context, input dtos.TagListInput) (res dtos.TagListResponse, err error)
}

var _tagUseCases TagUseCases

func GetTagUseCases() TagUseCases {
	if _tagUseCases == nil {
		_tagUseCases = &tagUseCasesImpl{
			db: database.GetConnection("default").(*database.GormConnection).Conn(),
		}
	}
	return _tagUseCases
}

// TagUseCasesImpl ...
type tagUseCasesImpl struct {
	db *gorm.DB
}

func (u *tagUseCasesImpl) Create(ctx context.Context, input dtos.TagCreateInput) (res response.TypedResponse[dtos.TagDTO], err error) {
	r := repositories.Use(u.db)
	t := r.Tag
	m := mappers.GetTagMapper()
	model := m.MapTagCreateInputToModel(ctx, input)
	err = t.WithContext(ctx).Create(model)
	if err != nil {
		return
	}
	res = response.NewTypedResponse(m.MapTagToDTO(ctx, model))
	return
}

func (u *tagUseCasesImpl) Update(ctx context.Context, id int64, input dtos.TagUpdateInput) (res response.TypedResponse[dtos.TagDTO], err error) {
	r := repositories.Use(u.db)
	t := r.Tag
	m := mappers.GetTagMapper()
	model := m.MapTagUpdateInputToModel(ctx, input)
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Updates(model)
	if err != nil {
		return
	}
	res = response.NewTypedResponse(m.MapTagToDTO(ctx, model))
	return
}

func (u *tagUseCasesImpl) DeleteById(ctx context.Context, id int64) (res response.Response, err error) {
	r := repositories.Use(u.db)
	t := r.Tag
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return
	}
	res = response.NewResponse()
	return
}

func (u *tagUseCasesImpl) FindById(ctx context.Context, id int64) (res response.TypedResponse[dtos.TagDTO], err error) {
	r := repositories.Use(u.db)
	t := r.Tag
	result, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return
	}
	res = response.NewTypedResponse(dtos.TagDTO{
		Tag: *result,
	})
	return
}

func (u *tagUseCasesImpl) FindAll(ctx context.Context) (res response.TypedResponse[[]dtos.TagDTO], err error) {
	r := repositories.Use(u.db)
	t := r.Tag
	result, err := t.WithContext(ctx).Find()
	if err != nil {
		return
	}
	res = response.NewTypedResponse([]dtos.TagDTO{})
	for _, v := range result {
		res.Result = append(res.Result, dtos.TagDTO{
			Tag: *v,
		})
	}
	return
}

func (u *tagUseCasesImpl) FindByOwner(ctx context.Context, input dtos.TagListInput) (res dtos.TagListResponse, err error) {
	r := repositories.Use(u.db)
	t := r.Tag
	result, err := t.WithContext(ctx).Where(t.OwnerType.Eq(input.OwnerType)).Find()
	if err != nil {
		return
	}

	res.Result.Items = make([]dtos.TagDTO, 0)
	for _, v := range result {
		res.Result.Items = append(res.Result.Items, dtos.TagDTO{
			Tag: *v,
		})
	}

	return
}
