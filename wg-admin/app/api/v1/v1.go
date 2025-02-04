package v1

import (
	"mp2720/wg-admin/wg-admin/app/services"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "mp2720/wg-admin/wg-admin/app/api/v1/docs"
)

//	@title		Wireguard admin server
//	@version	1.0

//	@host		localhost:8080
//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-Token

// apiBaseUrl is used as base for the urls that api returns.
func RegisterHandlers(e *echo.Echo, apiBaseUrl string, userService services.UserService) {
	v1 := e.Group("v1/")
	v1.GET("swagger/*", echoSwagger.WrapHandler)

	v1BaseUrl, err := url.JoinPath(apiBaseUrl, "v1")
	if err != nil {
		panic(err)
	}

	registerUserHandlers(v1, v1BaseUrl, userService)
}
