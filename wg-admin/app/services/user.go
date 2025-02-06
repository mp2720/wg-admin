package services

import (
	"context"
	"mp2720/wg-admin/wg-admin/storage"
	"mp2720/wg-admin/wg-admin/storage/data"
	"mp2720/wg-admin/wg-admin/transaction"
	"mp2720/wg-admin/wg-admin/utils"
	"time"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type RegisterRequest struct {
	Name         string
	IsAdmin      bool
	PrivateKey   *wgtypes.Key
	Fare         string
	MaxAddresses int64
}

type UpdateRequest struct {
	Name         *string
	IsAdmin      *bool
	IsBanned     *bool
	PrivateKey   *wgtypes.Key
	Fare         *string
	MaxAddresses *int64
}

type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (data.User, error)

	Get(ctx context.Context, uuid uuid.UUID) (data.User, error)

	GetAll(ctx context.Context) ([]data.User, error)

	Update(ctx context.Context, uuid uuid.UUID, req UpdateRequest) (data.User, error)

	SetTokenIssuedNow(ctx context.Context, uuid uuid.UUID) error

	MakePayment(ctx context.Context, uuid uuid.UUID, byTime time.Time) error

	DesyncUnpaid(ctx context.Context) error
}

type userService struct {
	userRepo    storage.UserRepo
	addressRepo storage.AddressRepo
	tm          transaction.Manager
	clock       utils.Clock
	// desyncNotifier desync.Notifier
}

func NewUserService(
	userRepo storage.UserRepo,
	addressRepo storage.AddressRepo,
	tm transaction.Manager,
	clock utils.Clock,
) UserService {
	return &userService{userRepo, addressRepo, tm, clock}
}

func (us *userService) Register(ctx context.Context, req RegisterRequest) (data.User, error) {
	user, err := data.NewUser(
		req.Name,
		req.IsAdmin,
		req.PrivateKey,
		req.Fare,
		req.MaxAddresses,
	)
	if err != nil {
		return data.User{}, err
	}

	if err = us.userRepo.Create(ctx, &user); err != nil {
		return data.User{}, err
	}

	return user, nil
}

func (us *userService) Get(ctx context.Context, uuid uuid.UUID) (data.User, error) {
	return us.userRepo.GetByUUID(ctx, uuid)
}

func (us *userService) GetAll(ctx context.Context) ([]data.User, error) {
	return us.userRepo.GetAll(ctx)
}

func (us *userService) Update(ctx context.Context, uuid uuid.UUID, req UpdateRequest) (data.User, error) {
	var user data.User
	err := us.tm.InTransaction(ctx, func(ctx context.Context) error {
		var err error

		user, err = us.userRepo.GetByUUIDLocked(ctx, uuid)
		if err != nil {
			return err
		}

		err = user.Update(data.UserPatch{
			Name:         req.Name,
			IsAdmin:      req.IsAdmin,
			IsBanned:     req.IsBanned,
			PrivateKey:   req.PrivateKey,
			Fare:         req.Fare,
			MaxAddresses: req.MaxAddresses,
		})
		if err != nil {
			return err
		}

		if err = us.userRepo.Save(ctx, user); err != nil {
			return err
		}

		if req.IsBanned != nil || req.PrivateKey != nil {
			if err := us.addressRepo.DesyncOwnedBy(ctx, &user, us.clock.Now()); err != nil {
				return err
			}
			panic("TODO: implement")
			// us.desyncNotifier.Notify()
		}

		return nil
	})
	return user, err
}

func (us *userService) update(ctx context.Context, uuid uuid.UUID, patch data.UserPatch) error {
	return us.tm.InTransaction(ctx, func(ctx context.Context) error {
		user, err := us.userRepo.GetByUUIDLocked(ctx, uuid)
		if err != nil {
			return err
		}

		if err = user.Update(patch); err != nil {
			return err
		}

		return us.userRepo.Save(ctx, user)
	})
}

func (us *userService) SetTokenIssuedNow(ctx context.Context, uuid uuid.UUID) error {
	now := us.clock.Now()
	return us.update(ctx, uuid, data.UserPatch{TokenIssuedAt: &now})
}

func (us *userService) MakePayment(ctx context.Context, uuid uuid.UUID, byTime time.Time) error {
	return us.update(ctx, uuid, data.UserPatch{PaidByTime: &byTime})
}

func (us *userService) DesyncUnpaid(ctx context.Context) error {
	panic("TODO: implement")
}
