package dtos

import (
	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/pkg/interactions/request"
	"github.com/pablor21/goms/pkg/interactions/response"
)

type TagDTO struct {
	models.Tag
} // @Name TagDTO

type TagCreateInput struct {
	request.Input
	Name      string `json:"name" form:"name"`
	Slug      string `json:"slug" form:"slug"`
	OwnerType string `json:"ownerType" form:"ownerType"`
} // @Name TagCreateInput

type TagUpdateInput struct {
	TagCreateInput
} // @Name TagUpdateInput

type TagListInput struct {
	request.PaginatedInput
	OwnerType string `json:"ownerType" param:"ownerType"`
} // @Name TagListInput

type TagGetInput struct {
	request.Input
	OwnerType string `json:"ownerType" param:"ownerType"`
} // @Name TagGetInput

type TagListResponse struct {
	response.PaginatedResponse[TagDTO]
} // @Name TagListResponse
