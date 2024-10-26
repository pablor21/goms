package api

import (
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/dtos"
)

var _authHandler *ApiAuthHandler

func AuthHandler() *ApiAuthHandler {
	if _authHandler == nil {
		_authHandler = &ApiAuthHandler{}
	}
	return _authHandler
}

func mapAuthApi(g *echo.Group) {
	authGroup := g.Group("/auth")
	authGroup.POST("/login", AuthHandler().Login)
}

type ApiAuthHandler struct{}

// @Title Login
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dtos.UserLoginInput true "dtos.UserLoginInput"
// @Success 200 {object}  dtos.UserLoginResponse
// @Router /auth/login [post]
func (h *ApiAuthHandler) Login(c echo.Context) error {
	var input dtos.UserLoginInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	// Do something with input
	return c.JSON(200, dtos.UserLoginResponse{})

}
