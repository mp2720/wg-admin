package v1

import (
	"errors"
	"fmt"
	"mp2720/wg-admin/wg-admin/app/services"
	"mp2720/wg-admin/wg-admin/storage/data"
	"mp2720/wg-admin/wg-admin/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type userAPI struct {
	baseUrl     string
	userService services.UserService
}

func registerUserHandlers(g *echo.Group, v1BaseUrl string, userService services.UserService) {
	userAPI := userAPI{v1BaseUrl, userService}

	g.POST("users", userAPI.registerUser)
	g.GET("users", userAPI.getAllUsers)
	g.GET("users/:uuid", userAPI.getUserByUUID)
}

type UserResponseLinks struct {
	Addresses string `json:"addresses"`
}

type UserResponse struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	IsAdmin      bool   `json:"is_admin"`
	IsBanned     bool   `json:"is_banned"`
	PublicKey    string `json:"public_key"`
	Fare         string `json:"fare"`
	AddressCount int64  `json:"address_count"`
	MaxAddresses int64  `json:"max_addresses"`

	LastSeenAt *time.Time `json:"last_seen_at"`
	PaidByTime *time.Time `json:"paid_by_time"`

	Links UserResponseLinks `json:"links"`
}

type RegisterUserRequest struct {
	Name         string  `json:"name"`
	IsAdmin      bool    `json:"is_admin"`
	PrivateKey   *string `json:"private_key"`
	Fare         string  `json:"fare"`
	MaxAddresses int64   `json:"max_addresses"`
}

func (api userAPI) makeUserResponse(user data.User) UserResponse {
	addressesPath, err := url.JoinPath(
		api.baseUrl,
		fmt.Sprintf("users/%s/addresses", url.PathEscape(user.UUID.String())),
	)
	if err != nil {
		panic(err)
	}

	return UserResponse{
		UUID:         user.UUID.String(),
		Name:         user.Name,
		IsAdmin:      user.IsAdmin,
		IsBanned:     user.IsBanned,
		PublicKey:    user.PrivateKey.PublicKey().String(),
		Fare:         user.Fare,
		AddressCount: user.AddressCount,
		MaxAddresses: user.MaxAddresses,
		LastSeenAt:   user.LastSeenAt,
		PaidByTime:   user.PaidByTime,
		Links: UserResponseLinks{
			Addresses: addressesPath,
		},
	}
}

// registerUser godoc
//
//	@Summary		Register a new user.
//	@Description	Register a new user. Admin only.
//	@Accept			json
//	@Produce		json
//	@Param			user	body		RegisterUserRequest	true	"user data"
//	@Success		201		{array}		UserResponse		"Ok"
//	@Failure		401		{object}	APIError			"Unauthorized"
//	@Failure		403		{object}	APIError			"Forbidden"
//	@Failure		409		{object}	APIError			"User with given name or key already exists"
//	@Security		ApiKeyAuth
//	@Router			/users [post]
func (api userAPI) registerUser(ctx echo.Context) error {
	var req RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, APIError{"Invalid request body"})
	}

	var optionalPrivateKey *wgtypes.Key
	if req.PrivateKey != nil {
		privateKey, err := wgtypes.ParseKey(*req.PrivateKey)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, APIError{"Invalid private key"})
		}
		optionalPrivateKey = &privateKey
	}

	user, err := api.userService.Register(ctx.Request().Context(), services.RegisterRequest{
		Name:         req.Name,
		IsAdmin:      req.IsAdmin,
		PrivateKey:   optionalPrivateKey,
		Fare:         req.Fare,
		MaxAddresses: req.MaxAddresses,
	})
	if err != nil {
		if errors.As(err, &utils.ErrAlreadyExists{}) {
			return ctx.JSON(http.StatusConflict, APIError{"User with given name or key already exists"})
		}
		if errors.Is(err, data.ErrInvalidUserName) || errors.Is(err, data.ErrInvalidAddressCount) {
			return ctx.JSON(http.StatusBadRequest, APIError{err.Error()})
		}

		panic(err)
	}

	return ctx.JSON(http.StatusCreated, api.makeUserResponse(user))
}

// getAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Get all registered users. Admin only.
//	@Produce		json
//	@Success		200	{array}		UserResponse	"Ok"
//	@Failure		401	{object}	APIError		"Unauthorized"
//	@Failure		403	{object}	APIError		"Forbidden"
//	@Security		ApiKeyAuth
//	@Router			/users [get]
func (api userAPI) getAllUsers(ctx echo.Context) error {
	users, err := api.userService.GetAll(ctx.Request().Context())
	if err != nil {
		panic(err)
	}

	var usersResp []UserResponse
	for _, user := range users {
		usersResp = append(usersResp, api.makeUserResponse(user))
	}

	return ctx.JSON(http.StatusOK, utils.NilToEmptySlice(usersResp))
}

// getUserByUUID godoc
//
//	@Summary		Get user by UUID.
//	@Description	Get user by UUID. Admin only.
//	@Param			uuid	path	string	true	"user's UUID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	UserResponse	"ok"
//	@Failure		401	{object}	APIError		"Unauthorized"
//	@Failure		403	{object}	APIError		"Forbidden"
//	@Failure		404	{object}	APIError		"Not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{uuid} [get]
func (api userAPI) getUserByUUID(ctx echo.Context) error {
	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, APIError{"Invalid UUID"})
	}

	println(uuid.String())

	user, err := api.userService.Get(ctx.Request().Context(), uuid)
	if err != nil {
		if errors.As(err, &utils.ErrAlreadyExists{}) {
			return ctx.JSON(http.StatusNotFound, APIError{"User not found"})
		}

		panic(err)
	}

	return ctx.JSON(http.StatusOK, api.makeUserResponse(user))
}
