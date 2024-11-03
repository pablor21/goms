package middlewares

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/server/session"
	"github.com/pablor21/goms/pkg/logger"
)

func Session(name string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			sessionManager := session.GetSessionManager().GetStore(name)

			var token string
			cookie, err := c.Cookie(sessionManager.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err = sessionManager.Load(ctx, token)
			if err != nil {
				return err
			}

			// save session in context
			ctx = session.SetSession(ctx, sessionManager)

			c.SetRequest(c.Request().WithContext(ctx))

			c.Response().Before(func() {
				// logger.Debug().Msg("Saving session")
				responseCookie := &http.Cookie{
					Name:     sessionManager.Cookie.Name,
					Path:     sessionManager.Cookie.Path,
					Domain:   sessionManager.Cookie.Domain,
					Secure:   sessionManager.Cookie.Secure,
					HttpOnly: sessionManager.Cookie.HttpOnly,
					SameSite: sessionManager.Cookie.SameSite,
					// Partitioned: true,
					// MaxAge:   int(sessionManager.Lifetime / time.Second),
					MaxAge: int(sessionManager.Lifetime.Seconds()),
				}

				if sessionManager.Status(ctx) != scs.Unmodified {
					switch sessionManager.Status(ctx) {
					case scs.Modified:
						token, t, err := sessionManager.Commit(ctx)
						if err != nil {
							logger.Error().Err(err).Msg("Error committing session")
						}
						responseCookie.MaxAge = int(time.Until(t).Seconds())
						responseCookie.Value = token

					case scs.Destroyed:
						responseCookie.Expires = time.Unix(1, 0)
						responseCookie.MaxAge = -1
					}
					c.SetCookie(responseCookie)

				} else if token != "" {
					token, _, err := sessionManager.Commit(ctx)
					if err != nil {
						logger.Error().Err(err).Msg("Error committing session")
					}
					responseCookie.Value = token
					c.SetCookie(responseCookie)
				}

			})

			c.Response().Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
			c.Response().Header().Add("Vary", "Cookie")

			return next(c)
		}
	}
}
