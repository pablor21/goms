package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/services"
)

type AssetsHandler struct {
}

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{}
}

func MapAssetRoutes(g *echo.Group) {
	h := NewAssetsHandler()
	g.GET("/assets/:storageName/:section/:uniqueId/:fileName", h.Download).Name = "assets.download"
	g.GET("/assets/:storageName/:section/:uniqueId/:fileName/:displayName", h.Download).Name = "assets.named-download"
	g.GET("/assets/:uniqueId", h.Details).Name = "assets.permalink"
	g.GET("/assets/:uniqueId/:displayName", h.Details).Name = "assets.named-permalink"
}

func (h *AssetsHandler) Download(c echo.Context) error {
	var input dtos.AssetDownloadInput
	if err := c.Bind(&input); err != nil {
		return err
	}
	service := services.GetAssetService()
	res, err := service.GetAssetFileReader(c.Request().Context(), input)
	if err != nil {
		return err
	}

	mimeType := res.MimeType

	c.Response().Header().Set(echo.HeaderContentType, mimeType)

	// immutable content
	c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	filename := input.DisplayName
	if filename == "" {
		filename = res.Uri
	}

	// display inline
	// c.Response().Header().Set("Content-Disposition", "inline; filename="+res.Metadata.GetStringOrDefault("displayFileName", res.Uri))

	http.ServeContent(c.Response(), c.Request(), filename, *res.UpdatedAt, res.Reader)
	return err
}

func (h *AssetsHandler) Details(c echo.Context) error {
	var input dtos.AssetDownloadInput
	if err := c.Bind(&input); err != nil {
		return err
	}
	service := services.GetAssetService()
	res, err := service.GetAssetDetails(c.Request().Context(), input)
	if err != nil {
		return err
	}

	headerData, err := json.Marshal(res.AssetDTO)
	if err != nil {
		return err
	}

	filename := input.DisplayName
	if filename == "" {
		filename = res.Metadata.GetStringOrDefault("displayFileName", res.Uri)
	}

	c.Response().Header().Set("X-Details", string(headerData))
	c.Response().Header().Set(echo.HeaderContentType, res.MimeType)
	c.Response().Header().Set("Content-Disposition", "inline; filename="+filename)
	http.ServeContent(c.Response(), c.Request(), filename, *res.UpdatedAt, res.Reader)
	return err

	// return c.JSON(http.StatusOK, res)
}
