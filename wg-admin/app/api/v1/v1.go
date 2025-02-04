package v1

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed swagger-ui/*
//go:embed api.yaml
var docFS embed.FS

func AddHandlers(e *echo.Echo) {
	g := e.Group("v1/")

	g.StaticFS("doc", docFS)
}
