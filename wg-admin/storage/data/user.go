package data

import (
	"errors"
	"regexp"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type User struct {
	ID           int64
	Name         string
	IsAdmin      bool
	IsBanned     bool
	PrivateKey   wgtypes.Key
	Fare         string
	AddressCount int
	MaxAddresses int

	TokenIssuedAt *time.Time
	LastSeenAt    *time.Time
	PaidByTime    *time.Time
}

var (
	ErrInvalidUserName          = errors.New("Invalid user name, it should contain only alphanums and underscore, at least one")
	ErrZeroPrivateKey           = errors.New("Zero private key was provided")
	ErrUserViolatesAddressLimit = errors.New("Address limit is violated")
)

func (u *User) validate() error {
	nameValid, err := regexp.Match("[A-Za-z0-9_]+", []byte(u.Name))
	if err != nil {
		return err
	}
	if !nameValid {
		return ErrInvalidUserName
	}

	if u.PrivateKey == (wgtypes.Key{}) {
		return ErrZeroPrivateKey
	}

	return nil
}

func NewUser(
	name string,
	isAdmin bool,
	privateKey *wgtypes.Key,
	fare string,
	maxAddresses int,
) (User, error) {
	if privateKey == nil {
		randomKey, err := wgtypes.GenerateKey()
		if err != nil {
			return User{}, err
		}
		privateKey = &randomKey
	}

	u := User{
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

// Nil means no update, TokenIssuedAt, LastSeenAt, PaidByTime could not be reset to nil.
type UserPatch struct {
	Name          *string
	IsAdmin       *bool
	PrivateKey    *wgtypes.Key
	Fare          *string
	AddressCount  *int
	MaxAddresses  *int
	TokenIssuedAt *time.Time
	LastSeenAt    *time.Time
	PaidByTime    *time.Time
}

func (u *User) Update(patch UserPatch) error {
	if patch.Name != nil {
		u.Name = *patch.Name
	}
	if patch.IsAdmin != nil {
		u.IsAdmin = *patch.IsAdmin
	}
	if patch.PrivateKey != nil {
		u.PrivateKey = *patch.PrivateKey
	}
	if patch.Fare != nil {
		u.Fare = *patch.Fare
	}
	if patch.AddressCount != nil {
		u.AddressCount = *patch.AddressCount
	}
	if patch.MaxAddresses != nil {
		u.MaxAddresses = *patch.MaxAddresses
	}
	if patch.TokenIssuedAt != nil {
		u.TokenIssuedAt = patch.TokenIssuedAt
	}
	if patch.LastSeenAt != nil {
		u.LastSeenAt = patch.LastSeenAt
	}
	if patch.PaidByTime != nil {
		u.PaidByTime = patch.PaidByTime
	}

	return u.validate()
}
