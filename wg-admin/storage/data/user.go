package data

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type User struct {
	ID           int64
	UUID         uuid.UUID
	Name         string
	IsAdmin      bool
	IsBanned     bool
	PrivateKey   wgtypes.Key
	Fare         string
	AddressCount int64
	MaxAddresses int64

	TokenIssuedAt *time.Time
	LastSeenAt    *time.Time
	PaidByTime    *time.Time
}

var (
	ErrInvalidUserName         = errors.New("Invalid user name, it should contain only alphanums and underscore, at least one")
	ErrInvalidAddressCount     = errors.New("Invalid address count, be non-negative integer is expected")
	ErrUserExceedsAddressLimit = errors.New("User's address limit is exceeded")
)

var userNameRegex = regexp.MustCompile("^[A-Za-z0-9_]+$")

func (u *User) validate() error {
	if ok := userNameRegex.Match([]byte(u.Name)); !ok {
		return ErrInvalidUserName
	}

	if u.AddressCount < 0 {
		return ErrInvalidAddressCount
	}

	if u.AddressCount > u.MaxAddresses {
		return ErrUserExceedsAddressLimit
	}

	return nil
}

func NewUser(
	name string,
	isAdmin bool,
	privateKey *wgtypes.Key,
	fare string,
	maxAddresses int64,
) (User, error) {
	if privateKey == nil {
		randomKey, err := wgtypes.GenerateKey()
		if err != nil {
			return User{}, err
		}
		privateKey = &randomKey
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return User{}, err
	}

	u := User{
		UUID:         uuid,
		Name:         name,
		IsAdmin:      isAdmin,
		IsBanned:     false,
		PrivateKey:   *privateKey,
		Fare:         fare,
		AddressCount: 0,
		MaxAddresses: maxAddresses,

		TokenIssuedAt: nil,
		LastSeenAt:    nil,
		PaidByTime:    nil,
	}

	return u, u.validate()
}

func (u User) CanHaveAuthenticationToken() bool {
	return !u.IsBanned
}

func (u User) CanBeAuthenticatedWithToken(issuedAt time.Time) bool {
	if u.IsBanned {
		return false
	}
	if u.TokenIssuedAt == nil {
		return false
	}
	return issuedAt.After(*u.TokenIssuedAt)
}

func (u User) PaidForTime(t time.Time) bool {
	if u.PaidByTime == nil {
		return false
	}
	return u.PaidByTime.After(t)
}

// Nil means no update, TokenIssuedAt, LastSeenAt, PaidByTime could not be reset to nil.
type UserPatch struct {
	Name          *string
	IsAdmin       *bool
	IsBanned      *bool
	PrivateKey    *wgtypes.Key
	Fare          *string
	AddressCount  *int64
	MaxAddresses  *int64
	TokenIssuedAt *time.Time
	LastSeenAt    *time.Time
	PaidByTime    *time.Time
}

func (u *User) Update(patch UserPatch) error {
	updated := *u

	if patch.Name != nil {
		updated.Name = *patch.Name
	}
	if patch.IsAdmin != nil {
		updated.IsAdmin = *patch.IsAdmin
	}
	if patch.IsBanned != nil {
		updated.IsBanned = *patch.IsBanned
	}
	if patch.PrivateKey != nil {
		updated.PrivateKey = *patch.PrivateKey
	}
	if patch.Fare != nil {
		updated.Fare = *patch.Fare
	}
	if patch.AddressCount != nil {
		updated.AddressCount = *patch.AddressCount
	}
	if patch.MaxAddresses != nil {
		updated.MaxAddresses = *patch.MaxAddresses
	}
	if patch.TokenIssuedAt != nil {
		updated.TokenIssuedAt = patch.TokenIssuedAt
	}
	if patch.LastSeenAt != nil {
		updated.LastSeenAt = patch.LastSeenAt
	}
	if patch.PaidByTime != nil {
		updated.PaidByTime = patch.PaidByTime
	}

	if err := updated.validate(); err != nil {
		return err
	}

	*u = updated
	return nil
}
