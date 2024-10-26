package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablor21/goms/app/config"
	"github.com/pablor21/goms/app/server/apps/api"
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

	e.Use(middleware.Recover())
	e.Use(lecho.Middleware(lecho.Config{
		Logger: eLogger,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowCredentials: true,
	}))
	api.MapApiRoutes(e.Group("/api"))

	return e.Start(fmt.Sprintf("%s:%d", s.Config.Server.Host, s.Config.Server.Port))
}
