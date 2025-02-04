package api

import (
	"mp2720/wg-admin/wg-admin/app/api/v1"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunHttpServer(address string) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1.AddHandlers(e)

	return e.Start(address)
}
