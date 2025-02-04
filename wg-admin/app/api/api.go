package api

import (
	"encoding/json"
	v1 "mp2720/wg-admin/wg-admin/app/api/v1"
	"mp2720/wg-admin/wg-admin/app/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunHTTPServer(
	address string,
	apiBaseURL string,
	authService services.AuthService,
	userService services.UserService,
) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1.RegisterHandlers(e, apiBaseURL, authService, userService)

	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		return err
	}
	println(string(data))

	return e.Start(address)
}
