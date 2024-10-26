package api

import (
	"io/fs"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs"
	"github.com/labstack/echo/v4"
	"github.com/pablor21/goms/app/server/apps/api/docs"
)

// @Version 1.0.0
// @Title Backend API
// @Description GOMS API server
// @Contact.Name Pablo Ramirez
// @Contact.Email pablo@pramirez.dev
// @Contact.URL https://pramirez.dev
// @TermsOfService https://pramirez.dev
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License
// @Server /api
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
func MapApiRoutes(e *echo.Group) {
	spec, err := fs.ReadFile(docs.ApiDocs, "swagger.json")
	if err != nil {
		panic(err)
	}

	// e.GET("/docs*", openapidocs.ScalarDocumentsHandler(openapidocs.ScalarConfig{
	// 	Spec: string(spec),
	// }))
	e.GET("/docs*", openapidocs.SwaggerUIDocumentsHandler(openapidocs.SwaggerUIConfig{
		Spec: string(spec),
	}))
	// se.GET("/docs/*", echoSwagger.WrapHandler)

	mapAuthApi(e)

	e.GET("/ping", Ping)
}

// @Title		Ping
// @Description	Health check
// @Tags		health
// @Accept		json
// @Produce	json
// @Success	200	{string}	string	"pong"
// @Router		/ping [get]
func Ping(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "pong",
	})
}
