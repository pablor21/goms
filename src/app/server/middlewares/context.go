package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/context"
)

func Context() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.SetServerContext(c.Request().Context(), c)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
