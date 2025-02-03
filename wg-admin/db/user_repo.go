package db

import (
	"context"
	"mp2720/wg-admin/wg-admin/db/sqlgen"
	"mp2720/wg-admin/wg-admin/storage"
	"mp2720/wg-admin/wg-admin/storage/data"
	"mp2720/wg-admin/wg-admin/transaction"
	"mp2720/wg-admin/wg-admin/utils"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type userRepo struct {
	db DB
	tm transaction.Manager
}

func NewUserRepo(db DB, tm transaction.Manager) storage.UserRepo {
	return userRepo{db, tm}
}

func mapUser(user sqlgen.User) (data.User, error) {
	privateKey, err := wgtypes.ParseKey(user.PrivateKey)
	if err != nil {
		return data.User{}, err
	}

	return data.User{
		ID:            user.ID,
		Name:          user.Name,
		IsAdmin:       user.IsAdmin,
		IsBanned:      user.IsBanned,
		PrivateKey:    privateKey,
		Fare:          user.Fare,
		AddressCount:  user.AddressCount,
		MaxAddresses:  user.MaxAddresses,
		TokenIssuedAt: user.TokenIssuedAt,
		LastSeenAt:    user.LastSeenAt,
		PaidByTime:    user.PaidByTime,
	}, nil
}

func (ur userRepo) Create(
	ctx context.Context,
	user *data.User,
) error {
	id, err := ur.db.With(ctx).CreateUser(
		ctx,
		sqlgen.CreateUserParams{
			Name:           user.Name,
			IsAdmin:        user.IsAdmin,
			IsBanned:       user.IsBanned,
			Fare:           user.Fare,
			PrivateKey:     user.PrivateKey.String(),
			PublicKey:      user.PrivateKey.PublicKey().String(),
			AddressesCount: user.AddressCount,
			MaxAddresses:   user.MaxAddresses,
			PaidByTime:     user.PaidByTime,
			TokenIssuedAt:  user.TokenIssuedAt,
			LastSeenAt:     user.LastSeenAt,
		},
	)
	if err != nil {
		return HandleSQLError("user", err)
	}

	user.ID = id
	return nil
}

func (ur userRepo) GetByName(ctx context.Context, name string) (data.User, error) {
	user, err := ur.db.With(ctx).GetUserByName(ctx, name)
	if err != nil {
		return data.User{}, HandleSQLError("user", err)
	}

	return mapUser(user)
}

func (ur userRepo) GetByNameLocked(ctx context.Context, name string) (data.User, error) {
    // NOTE: there's no SELECT FOR UPDATE in SQLite, but if txlock = immediate this will fork fine
	return ur.GetByName(ctx, name)
}

func (ur userRepo) GetAll(ctx context.Context) ([]data.User, error) {
	users, err := ur.db.With(ctx).GetAllUsers(ctx)
	if err != nil {
		return nil, HandleSQLError("user", err)
	}

	var mappedUsers []data.User
	for _, user := range users {
		mappedUser, err := mapUser(user)
		if err != nil {
			return nil, HandleSQLError("user", err)
		}

		mappedUsers = append(mappedUsers, mappedUser)
	}

	return mappedUsers, nil
}

func (ur userRepo) Save(ctx context.Context, user data.User) error {
	rowsAffected, err := ur.db.With(ctx).UpdateUser(ctx,
		sqlgen.UpdateUserParams{
			ID:            user.ID,
			Name:          user.Name,
			IsAdmin:       user.IsAdmin,
			IsBanned:      user.IsBanned,
			Fare:          user.Fare,
			PrivateKey:    user.PrivateKey.String(),
			PublicKey:     user.PrivateKey.PublicKey().String(),
			AddressCount:  user.AddressCount,
			MaxAddresses:  user.MaxAddresses,
			PaidByTime:    user.PaidByTime,
			TokenIssuedAt: user.TokenIssuedAt,
			LastSeenAt:    user.LastSeenAt,
		},
	)
	if err != nil {
		return HandleSQLError("user", err)
	}
	if rowsAffected == 0 {
		return utils.ErrNotFound{What: "user"}
	}

	return nil
}

func (ur userRepo) Delete(ctx context.Context, name string) error {
	rowsAffected, err := ur.db.With(ctx).DeleteUser(ctx, name)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return utils.ErrNotFound{What: "user"}
	}

	return nil
}
