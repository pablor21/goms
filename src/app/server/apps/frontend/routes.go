package frontend

import "github.com/labstack/echo/v4"

func MapFrontendRoutes(g *echo.Group) {
	g.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
}
