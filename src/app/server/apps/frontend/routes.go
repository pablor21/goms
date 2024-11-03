package frontend

import (
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/server/apps/frontend/handlers"
)

func MapFrontendRoutes(g *echo.Group) {
	g.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	handlers.MapAssetRoutes(g)

	// g.GET("/assets/:storageName/:uniqueId/:fileName", func(c echo.Context) error {
	// 	return c.String(200, "List of assets")
	// }).Name = "assets.download"
}
