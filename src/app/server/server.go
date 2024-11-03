package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablor21/goms/app/config"
	"github.com/pablor21/goms/app/server/apps/api"
	"github.com/pablor21/goms/app/server/apps/frontend"
	"github.com/pablor21/goms/app/server/middlewares"
	"github.com/pablor21/goms/pkg/errors"
	"github.com/pablor21/goms/pkg/logger"
	"github.com/ziflex/lecho/v3"
)

type Server struct {
	Config *config.Config
	Echo   *echo.Echo
}

func NewServer() *Server {
	return &Server{
		Config: config.GetConfig(),
		Echo:   echo.New(),
	}
}

func (s *Server) Start() error {
	e := s.Echo
	e.HideBanner = true
	eLogger := lecho.From(logger.Logger())
	e.Logger = eLogger
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(lecho.Middleware(lecho.Config{
		Logger: eLogger,
	}))

	// Custom middlewares
	e.Use(middlewares.Context())
	e.Use(middlewares.UrlGenerator())

	// Map routes
	api.MapApiRoutes(e.Group("/api"))
	frontend.MapFrontendRoutes(e.Group(""))

	return e.Start(fmt.Sprintf("%s:%d", s.Config.Server.Host, s.Config.Server.Port))
}

func customHTTPErrorHandler(err error, c echo.Context) {

	code := http.StatusInternalServerError
	path := c.Request().URL.Path
	if c.Request().URL.RawQuery != "" {
		path += "?" + c.Request().URL.RawQuery
	}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		err = errors.NewAppError(http.StatusText(code), code)
	} else if appErr, ok := err.(*errors.AppError); ok {
		code = appErr.StatusCode
		err = appErr
	} else if valdiationErr, ok := err.(*errors.ValidationError); ok {
		code = valdiationErr.StatusCode
		err = valdiationErr
	} else {
		err = errors.NewAppErrorFromError(err)
	}

	// only log the error if it's not a 404
	if code != http.StatusNotFound {
		// custom_log.Error().Msg(err.Error() + ", url=" + path)
	}

	c.Response().Status = code
	// if is an API request
	if strings.Contains(c.Request().Header.Get("accept"), "text/html") && !strings.HasPrefix(c.Request().URL.Path, "/api") {
		// err2 := handlers.RenderError(c, err)
		// if err2 != nil {
		// 	c.JSON(code, err)
		// }
		c.JSON(code, err)
		return
	}

	c.JSON(code, err)
}
