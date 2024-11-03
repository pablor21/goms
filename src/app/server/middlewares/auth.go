package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/server/session"
	"github.com/pablor21/goms/app/services"
	"github.com/pablor21/goms/pkg/auth"
	"github.com/pablor21/goms/pkg/errors"
)

// Middleware to require authentication
// Must be used after the AuthenticateRequest middleware
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !auth.IsAuthenticated(c.Request().Context()) {
				return errors.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func RequireAnonymous() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if auth.IsAuthenticated(c.Request().Context()) {
				return errors.ErrUnauthorized
			}
			return next(c)
		}
	}
}

// Puts the principal in the context if the session is valid
// Must be used after the session middleware
func AuthenticateRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			session := session.GetSession(c.Request().Context())
			if session != nil {
				rawUser, err := session.GetPrincipal()
				if err != nil {
					return next(c)
				}
				authCtx := auth.SetContextPrincipal(ctx, rawUser)
				c.SetRequest(c.Request().WithContext(authCtx))

			}
			return next(c)
		}
	}
}

// Puts the user in the context if the session is valid
// Must be used after the authenticate request middleware
func AddUserToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			session := session.GetSession(c.Request().Context())
			if session != nil {
				rawPrincipal := auth.GetContextPrincipal(ctx)
				if rawPrincipal != nil {
					_, ctx, err := services.GetAuthService().GetContextUser(ctx)
					if err == nil {
						c.SetRequest(c.Request().WithContext(ctx))
					}
				}

			}
			return next(c)
		}
	}
}
