package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/usecases"
)

var _tagsHandler *ApiTagsHandler

func TagsHandler() *ApiTagsHandler {
	if _tagsHandler == nil {
		_tagsHandler = &ApiTagsHandler{}
	}
	return _tagsHandler
}

type ApiTagsHandler struct{}

func mapTagsApi(g *echo.Group) {
	tagsGroup := g.Group("/tags")
	tagsGroup.GET("/by-owner/:ownerType", TagsHandler().List).Name = "api.tags.list"
	tagsGroup.POST("", TagsHandler().Create).Name = "api.tags.create"
	tagsGroup.GET("/:id", TagsHandler().Get).Name = "api.tags.get"
	tagsGroup.PUT("/:id", TagsHandler().Update).Name = "api.tags.update"
	tagsGroup.DELETE("/:id", TagsHandler().Delete).Name = "api.tags.delete"
}

// @Title ListTags
// @Id ListTags
// @Summary ListTags
// @Description ListTags
// @Tags Tags
// @Accept json
// @Produce json
// @Param ownerType path string true "Owner Type"
// @Success 200 {object}  dtos.TagListResponse
// @Router /tags/by-owner/{ownerType} [get]
func (h *ApiTagsHandler) List(c echo.Context) (err error) {
	var input dtos.TagListInput
	err = c.Bind(&input)
	if err != nil {
		return
	}

	res, err := usecases.GetTagUseCases().FindByOwner(c.Request().Context(), input)
	if err != nil {
		return
	}
	return c.JSON(200, res)
}

// @Title CreateTag
// @Id CreateTag
// @Summary CreateTag
// @Description CreateTag
// @Tags Tags
// @Accept json
// @Produce json
// @Param input body dtos.TagCreateInput true "Tag data"
// @Success 200 {object}  dtos.TagDTO
// @Router /tags [post]
func (h *ApiTagsHandler) Create(c echo.Context) (err error) {
	var input dtos.TagCreateInput
	err = c.Bind(&input)
	if err != nil {
		return
	}
	res, err := usecases.GetTagUseCases().Create(c.Request().Context(), input)
	if err != nil {
		return
	}
	return c.JSON(200, res)
}

// @Title GetTag
// @Id GetTag
// @Summary GetTag
// @Description GetTag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object}  dtos.TagDTO
// @Router /tags/{id} [get]
func (h *ApiTagsHandler) Get(c echo.Context) (err error) {
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return
	}
	res, err := usecases.GetTagUseCases().FindById(c.Request().Context(), id)
	if err != nil {
		return
	}
	return c.JSON(200, res)
}

// @Title UpdateTag
// @Id UpdateTag
// @Summary UpdateTag
// @Description UpdateTag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param input body dtos.TagUpdateInput true "Tag data"
// @Success 200 {object}  dtos.TagDTO
// @Router /tags/{id} [put]
func (h *ApiTagsHandler) Update(c echo.Context) (err error) {
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return
	}
	var input dtos.TagUpdateInput
	err = c.Bind(&input)
	if err != nil {
		return
	}
	res, err := usecases.GetTagUseCases().Update(c.Request().Context(), id, input)
	if err != nil {
		return
	}
	return c.JSON(200, res)
}

// @Title DeleteTag
// @Id DeleteTag
// @Summary DeleteTag
// @Description DeleteTag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object}  response.Response
// @Router /tags/{id} [delete]
func (h *ApiTagsHandler) Delete(c echo.Context) (err error) {
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return
	}
	res, err := usecases.GetTagUseCases().DeleteById(c.Request().Context(), id)
	if err != nil {
		return
	}
	return c.JSON(200, res)
}
