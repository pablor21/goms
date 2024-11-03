package api

import (
	"io/fs"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablor21/goms/app/server/apps/api/docs"
	"github.com/pablor21/goms/app/server/middlewares"
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
// @BasePath /api
// @BaseURL http://localhost:8080/api
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
func MapApiRoutes(e *echo.Group) {

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowCredentials: true,
	}))

	// custom middlewares
	g := e.Group("")
	g.Use(middlewares.Session("default"))
	g.Use(middlewares.AuthenticateRequest())

	// routes
	mapApiDocs(e)
	g.GET("/ping", pingHanlder)
	mapAuthApi(g)
	mapTagsApi(g)
}

func mapApiDocs(e *echo.Group) {
	spec, err := fs.ReadFile(docs.ApiDocs, "swagger.json")
	if err != nil {
		panic(err)
	}

	e.GET("/docs*", openapidocs.SwaggerUIDocumentsHandler(openapidocs.SwaggerUIConfig{
		DeepLinking: true,
		Spec:        string(spec),
	}))
}

// @Id			HealthCheck
// @Title		pingHanlder
// @Description	Health check
// @Tags		health
// @Accept		json
// @Produce	json
// @Success	200	{string}	string	"ok"
// @Router		/ping [get]
func pingHanlder(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "ok",
		// "session": s.Values["foo"],
	})
}
