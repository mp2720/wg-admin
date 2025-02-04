package data_test

import (
	"mp2720/wg-admin/wg-admin/storage/data"
	"mp2720/wg-admin/wg-admin/utils/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func Test_NewUser(t *testing.T) {
	privateKey := testutils.MustGenerateWireguardPrivateKey()

	_, err := data.NewUser("", true, &privateKey, "300$", 2)
	require.ErrorIs(t, err, data.ErrInvalidUserName)

	_, err = data.NewUser("user", true, &privateKey, "300$", -1)
	require.ErrorIs(t, err, data.ErrUserExceedsAddressLimit)

	_, err = data.NewUser("fasfd    \n", true, &privateKey, "300$", 0)
	require.ErrorIs(t, err, data.ErrInvalidUserName)

	user, err := data.NewUser("user", true, &privateKey, "300$", 3)
	require.NoError(t, err)
	require.NotEqual(t, uuid.UUID{}, user.UUID)
	require.Equal(t, data.User{
		UUID:          user.UUID,
		Name:          "user",
		IsAdmin:       true,
		PrivateKey:    privateKey,
		Fare:          "300$",
		AddressCount:  0,
		MaxAddresses:  3,
		TokenIssuedAt: nil,
		LastSeenAt:    nil,
		PaidByTime:    nil,
	}, user)

	user2, err := data.NewUser("user", true, nil, "300$", 3)
	require.NoError(t, err)
	require.NotEqual(t, wgtypes.Key{}, user2.PrivateKey) // may fail :(
	user.PrivateKey = user2.PrivateKey
	user.UUID = user2.UUID
	require.Equal(t, user, user2)
}

func Test_Update(t *testing.T) {
	privateKey := testutils.MustGenerateWireguardPrivateKey()

	user, err := data.NewUser("user", false, &privateKey, "100$", 10)
	require.NoError(t, err)

	newName := "new_user_name"
	newIsAdmin := true
	newIsBanned := true
	newPrivateKey := testutils.MustGenerateWireguardPrivateKey()
	newFare := "0.4$"
	newAddressCount := int64(2)
	newMaxAddresses := int64(4)
	newTokenIssuedAt := time.Now()
	newLastSeenAt := time.Now()
	newPaidByTime := time.Now()

	patch := data.UserPatch{
		Name:          &newName,
		IsAdmin:       &newIsAdmin,
		IsBanned:      &newIsBanned,
		PrivateKey:    &newPrivateKey,
		Fare:          &newFare,
		AddressCount:  &newAddressCount,
		MaxAddresses:  &newMaxAddresses,
		TokenIssuedAt: &newTokenIssuedAt,
		LastSeenAt:    &newLastSeenAt,
		PaidByTime:    &newPaidByTime,
	}

	// check all fields are updated
	err = user.Update(patch)
	require.NoError(t, err)
	require.Equal(t, data.User{
		UUID:          user.UUID,
		Name:          newName,
		IsAdmin:       newIsAdmin,
		IsBanned:      newIsBanned,
		PrivateKey:    newPrivateKey,
		Fare:          newFare,
		AddressCount:  newAddressCount,
		MaxAddresses:  newMaxAddresses,
		TokenIssuedAt: &newTokenIssuedAt,
		LastSeenAt:    &newLastSeenAt,
		PaidByTime:    &newPaidByTime,
	}, user)

	// check all fields stay the same
	userCopy := user
	err = userCopy.Update(patch)
	require.NoError(t, err)
	require.Equal(t, user, userCopy)

	newBadName := ""
	err = userCopy.Update(data.UserPatch{Name: &newBadName})
	require.ErrorIs(t, err, data.ErrInvalidUserName)
	require.Equal(t, user, userCopy)

	newBadMaxAddresses := int64(1)
	err = userCopy.Update(data.UserPatch{MaxAddresses: &newBadMaxAddresses})
	require.ErrorIs(t, err, data.ErrUserExceedsAddressLimit)
	require.Equal(t, user, userCopy)

	newBadAddressCount := int64(-1)
	err = userCopy.Update(data.UserPatch{AddressCount: &newBadAddressCount})
	require.ErrorIs(t, err, data.ErrInvalidAddressCount)
	require.Equal(t, user, userCopy)
}
