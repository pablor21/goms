package api

import (
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/server/middlewares"
	"github.com/pablor21/goms/app/server/session"
	"github.com/pablor21/goms/app/services"
	"github.com/pablor21/goms/pkg/auth"
	"github.com/pablor21/goms/pkg/interactions/request"
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
	authGroup.Any("/logout", AuthHandler().Logout)
	authGroup.POST("/request-otp", AuthHandler().RequestOtp).Name = "api.auth.request-otp"
	authGroup.POST("/otp-login", AuthHandler().LoginWithOTP).Name = "api.auth.otp-login"
	authGroup.GET("/otp-login", AuthHandler().LoginWithOTP).Name = "api.auth.otp-login-token"

	// Authenticated routes
	authenticatedGroup := authGroup.Group("", middlewares.AuthenticateRequest())
	authenticatedGroup.GET("/me", AuthHandler().Me).Name = "api.auth.me"
	authenticatedGroup.PUT("/update-profile", AuthHandler().UpdateProfile).Name = "api.auth.update-profile"

}

type ApiAuthHandler struct{}

// @Title Login
// @Id Login
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dtos.UserLoginInput true "Login data"
// @Success 200 {object}  dtos.UserLoginResponse
// @Router /auth/login [post]
func (h *ApiAuthHandler) Login(c echo.Context) (err error) {
	var input dtos.UserLoginInput
	if err = c.Bind(&input); err != nil {
		return
	}
	s := session.GetSession(c.Request().Context())

	// Renew token to avoid session fixation
	s.RenewToken()
	service := services.GetAuthService()
	data, ctx, err := service.Login(c.Request().Context(), input)
	if err != nil {
		s.SetPrincipal(nil)
		return err
	}

	res := dtos.UserLoginResponse{}.SetCode(200).SetResult(data)

	// Set user in session
	s.SetPrincipal(auth.GetContextPrincipal(ctx))

	// set request context
	c.SetRequest(c.Request().WithContext(ctx))

	// Do something with input
	return c.JSON(200, res)

}

// @Title Logout
// @Id Logout
// @Summary Logout
// @Description Logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /auth/logout [post]
func (h *ApiAuthHandler) Logout(c echo.Context) (err error) {
	s := session.GetSession(c.Request().Context())
	s.Destroy()
	return c.NoContent(204)
}

// @Title Me
// @Id Me
// @Summary Me
// @Description Me
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object}  dtos.UserResponse
// @Router /auth/me [get]
func (h *ApiAuthHandler) Me(c echo.Context) (err error) {
	data, _, err := services.GetAuthService().GetContextUser(c.Request().Context())
	if err != nil {
		return err
	}

	res := dtos.UserResponse{}.SetCode(200).SetResult(data)

	return c.JSON(200, res)
}

// @Title UpdateProfile
// @Id UpdateProfile
// @Summary UpdateProfile
// @Description UpdateProfile
// @Tags Auth
// @Accept json, multipart/form-data
// @Produce json
// // @Param input body dtos.UpdateProfileInput true "the request body, only available for application/json content type"
// @Param avatar formData file false "avatar"
// @Param firstName formData string false "firstName"
// @Param lastName formData string false "lastName"
// @Param email formData string false "email"
// @Param phoneNumber formData string false "phoneNumber"
// @Param lang formData string false "lang"
// @Success 200 {object}  dtos.UserResponse
// @Router /auth/update-profile [put]
func (h *ApiAuthHandler) UpdateProfile(c echo.Context) (err error) {
	var input dtos.UpdateProfileInput
	if err = c.Bind(&input); err != nil {
		return
	}

	// check if request is multipart/form-data
	if request.IsMultiPart(c.Request()) {
		// get avatar file
		file, _ := c.FormFile("avatar")
		input.Avatar = *request.NewFileInputFromMultipart(file)
	}

	data, err := services.GetAuthService().UpdateProfile(c.Request().Context(), input)
	if err != nil {
		return err
	}

	res := dtos.UserResponse{}.SetCode(200).SetResult(data)

	return c.JSON(200, res)
}

// @Title RequestOtp
// @Id RequestOtp
// @Summary RequestOtp
// @Description RequestOtp
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dtos.RequestOTPInput true "Request OTP data"
// @Success 200 {object}  dtos.RequestOTPResponse
// @Router /auth/request-otp [post]
func (h *ApiAuthHandler) RequestOtp(c echo.Context) (err error) {
	var input dtos.RequestOTPInput
	if err = c.Bind(&input); err != nil {
		return
	}

	data, err := services.GetAuthService().RequestOTP(c.Request().Context(), input)
	if err != nil {
		return err
	}

	res := dtos.RequestOTPResponse{}.SetCode(200).SetResult(data)

	return c.JSON(200, res)
}

// @Title LoginWithOTP
// @Id LoginWithOTP
// @Summary LoginWithOTP
// @Description LoginWithOTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dtos.LoginWithOTPInput true "Login with OTP data"
// @Success 200 {object}  dtos.UserLoginResponse
// @Router /auth/otp-login [post]
func (h *ApiAuthHandler) LoginWithOTP(c echo.Context) (err error) {
	var input dtos.LoginWithOTPInput
	if err = c.Bind(&input); err != nil {
		return
	}

	s := session.GetSession(c.Request().Context())

	// Renew token to avoid session fixation
	s.RenewToken()
	service := services.GetAuthService()
	data, ctx, err := service.LoginWithOTP(c.Request().Context(), input)
	if err != nil {
		s.SetPrincipal(nil)
		return err
	}

	res := dtos.UserLoginResponse{}.SetCode(200).SetResult(data)

	// Set user in session
	s.SetPrincipal(auth.GetContextPrincipal(ctx))

	// set request context
	c.SetRequest(c.Request().WithContext(ctx))

	// Do something with input
	return c.JSON(200, res)
}
