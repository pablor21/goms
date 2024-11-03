package context

import (
	"context"

	"github.com/labstack/echo/v4"
)

type UrlGenerator func(string, ...interface{}) string

type UrlGeneratorKey string

var UrlGeneratorKeyValue UrlGeneratorKey = "urlGenerator"

type ServerContextKey string

var ServerContextKeyValue ServerContextKey = "server"

func GetUrlGenerator(ctx context.Context) UrlGenerator {
	generator, ok := ctx.Value(UrlGeneratorKeyValue).(UrlGenerator)
	if !ok {
		return nil
	}
	return generator
}

func SetUrlGenerator(ctx context.Context, generator UrlGenerator) context.Context {
	return context.WithValue(ctx, UrlGeneratorKeyValue, generator)
}

func SetServerContext(ctx context.Context, server echo.Context) context.Context {
	return context.WithValue(ctx, ServerContextKeyValue, server)
}

func GetServerContext(ctx context.Context) echo.Context {
	server, ok := ctx.Value(ServerContextKeyValue).(echo.Context)
	if !ok {
		return nil
	}
	return server
}
