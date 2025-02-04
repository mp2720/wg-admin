package server

import (
	v1 "mp2720/wg-admin/wg-admin/app/api/v1"
	"mp2720/wg-admin/wg-admin/app/user"

	"github.com/labstack/echo/v4"
)

type userServer struct {
	service user.Service
}

var _ v1.ServerInterface = (*userServer)(nil)

// (DELETE /addresses/{address})
func (us *userServer) DeleteAddressesAddress(ctx echo.Context, address string) error {
    
}

// (PUT /addresses/{address})
// PutAddressesAddress(ctx echo.Context, address string) error
//
// // (POST /addresses/{version})
// PostAddressesVersion(ctx echo.Context, version PostAddressesVersionParamsVersion) error
//
// // (GET /me)
// GetMe(ctx echo.Context) error
//
// // (GET /users)
// GetUsers(ctx echo.Context) error
//
// // (POST /users)
// PostUsers(ctx echo.Context) error
//
// // (DELETE /users/{uuid})
// DeleteUsersUuid(ctx echo.Context, uuid string) error
//
// // (GET /users/{uuid})
// GetUsersUuid(ctx echo.Context, uuid string) error
//
// // (PUT /users/{uuid})
// PutUsersUuid(ctx echo.Context, uuid string) error
//
// // (GET /users/{uuid}/addresses)
// GetUsersUuidAddresses(ctx echo.Context, uuid string) error
//
// // (POST /users/{uuid}/key)
// PostUsersUuidKey(ctx echo.Context, uuid string) error
//
// // (POST /users/{uuid}/token)
// PostUsersUuidToken(ctx echo.Context, uuid string) error
