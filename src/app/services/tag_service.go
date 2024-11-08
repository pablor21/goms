package services

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/mappers"
	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/app/repositories"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/errors"
	"github.com/pablor21/goms/pkg/interactions/response"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type TagService interface {
	Create(ctx context.Context, input dtos.TagCreateInput) (res *models.Tag, err error)
	Update(ctx context.Context, id int64, input dtos.TagUpdateInput) (res *models.Tag, err error)
	Delete(ctx context.Context, id int64) (res int, err error)
	GetByID(ctx context.Context, id int64) (res *models.Tag, err error)
	GetBySlug(ctx context.Context, slug string) (res *models.Tag, err error)
	FindAll(ctx context.Context, input dtos.TagListInput) (res response.PaginatedResponse[*models.Tag], err error)
	FindByOwner(ctx context.Context, input dtos.TagListInput) (res response.PaginatedResponse[*models.Tag], err error)
}

var _tagService TagService

func GetTagService() TagService {
	if _tagService == nil {
		_tagService = &tagServiceImpl{
			db: database.GetConnection("default").(*database.GormConnection),
		}
	}
	return _tagService
}

// TagServiceImpl ...
type tagServiceImpl struct {
	db *database.GormConnection
}

func (s *tagServiceImpl) Create(ctx context.Context, input dtos.TagCreateInput) (res *models.Tag, err error) {
	tx, ctx, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()
	r := repositories.Use(tx.DB)
	t := r.Tag
	m := mappers.GetTagMapper()
	model := m.MapTagCreateInputToModel(ctx, input)

	err = validation.ValidateStructWithContext(ctx, model,
		validation.Field(&model.Name, validation.Required),
		validation.Field(&model.Slug, validation.Required),
		validation.Field(&model.OwnerType, validation.Required),
	)
	if err != nil {
		err = errors.NewValidationError(err)
		return
	}

	if model.ParentID != nil && *model.ParentID > 0 {
		parent, err := t.WithContext(ctx).Where(t.ID.Eq(*model.ParentID)).Where(t.OwnerType.Eq(model.OwnerType)).First()
		if err != nil {
			return res, err
		}
		// set the complete slug
		model.CompleteSlug = parent.CompleteSlug + "/" + model.Slug
	} else {
		model.CompleteSlug = model.Slug
	}

	err = t.WithContext(ctx).Create(model)
	if err != nil {
		return
	}
	res = model
	err = tx.Commit()
	return
}

func (s *tagServiceImpl) Update(ctx context.Context, id int64, input dtos.TagUpdateInput) (res *models.Tag, err error) {
	tx, ctx, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()
	r := repositories.Use(tx.DB)
	t := r.Tag

	m := mappers.GetTagMapper()
	model := m.MapTagUpdateInputToModel(ctx, input)

	err = validation.ValidateStructWithContext(ctx, model,
		validation.Field(&model.Name, validation.Required),
		validation.Field(&model.Slug, validation.Required),
		validation.Field(&model.OwnerType, validation.Required),
		validation.Field(&model.ParentID, validation.When(model.ParentID != nil && *model.ParentID > 0, validation.Required, validation.By(func(value interface{}) error {
			if value == nil {
				return nil
			}
			parent, err := t.WithContext(ctx).Where(t.ID.Eq(*model.ParentID)).Where(t.OwnerType.Eq(model.OwnerType)).First()
			if err != nil {
				return errors.ErrRecordNotFound
			}
			if parent != nil {
				model.CompleteSlug = parent.CompleteSlug + "/" + model.Slug
			}
			return nil
		}))),
	)
	if err != nil {
		err = errors.NewValidationError(err)
		return
	}

	if model.ParentID != nil && *model.ParentID > 0 {
		parent, err := t.WithContext(ctx).Where(t.ID.Eq(*model.ParentID)).Where(t.OwnerType.Eq(model.OwnerType)).First()
		if err != nil {
			return res, err
		}
		// set the complete slug
		model.CompleteSlug = parent.Slug + "/" + model.Slug
	} else {
		model.CompleteSlug = model.Slug
	}

	// check if the tag exists
	eModel, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Preload(t.Children).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}

	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Omit(field.AssociationFields).Updates(model)
	if err != nil {
		return
	}

	err = s.updateChildSlug(ctx, eModel)
	if err != nil {
		return
	}

	eModel, err = t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}

	res = eModel
	err = tx.Commit()
	return
}

func (s *tagServiceImpl) updateChildSlug(ctx context.Context, parent *models.Tag) (err error) {
	if len(parent.Children) == 0 {
		return
	}
	tx, ctx, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()
	r := repositories.Use(tx.DB)
	t := r.Tag
	for _, v := range parent.Children {
		v.CompleteSlug = parent.CompleteSlug + "/" + v.Slug
		_, err = t.WithContext(ctx).Where(t.ID.Eq(v.ID)).Updates(v)
		if err != nil {
			return
		}
		err = s.updateChildSlug(ctx, v)
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}

func (s *tagServiceImpl) Delete(ctx context.Context, id int64) (res int, err error) {
	tx, ctx, err := s.db.GetContextTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()
	r := repositories.Use(tx.DB)
	t := r.Tag

	// check if the tag exists
	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
		return
	}

	_, err = t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (s *tagServiceImpl) GetByID(ctx context.Context, id int64) (res *models.Tag, err error) {
	db, ctx, err := s.db.GetContextDb(ctx)
	if err != nil {
		return
	}

	t := repositories.Use(db).Tag
	res, err = t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
	}
	return
}

func (s *tagServiceImpl) GetBySlug(ctx context.Context, slug string) (res *models.Tag, err error) {
	db, ctx, err := s.db.GetContextDb(ctx)
	if err != nil {
		return
	}
	t := repositories.Use(db).Tag
	res, err = t.WithContext(ctx).Where(t.Slug.Eq(slug)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.ErrRecordNotFound
		}
	}
	return
}

func (s *tagServiceImpl) FindAll(ctx context.Context, input dtos.TagListInput) (res response.PaginatedResponse[*models.Tag], err error) {
	db, ctx, err := s.db.GetContextDb(ctx)
	if err != nil {
		return
	}
	r := repositories.Use(db)
	t := r.Tag
	result, err := t.WithContext(ctx).Where(t.OwnerType.Eq(input.OwnerType)).Find()
	if err != nil {
		return
	}

	res.Result.Items = result

	return
}

func (s *tagServiceImpl) FindByOwner(ctx context.Context, input dtos.TagListInput) (res response.PaginatedResponse[*models.Tag], err error) {
	db, ctx, err := s.db.GetContextDb(ctx)
	if err != nil {
		return
	}
	r := repositories.Use(db)
	t := r.Tag
	result, err := t.WithContext(ctx).Where(t.OwnerType.Eq(input.OwnerType)).Find()
	if err != nil {
		return
	}

	res.Result.Items = result

	return
}
