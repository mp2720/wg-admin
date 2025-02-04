package user

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

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (data.User, error)

	Get(ctx context.Context, uuid uuid.UUID) (data.User, error)

	GetAll(ctx context.Context) ([]data.User, error)

	Update(ctx context.Context, uuid uuid.UUID, req UpdateRequest) (data.User, error)

	SetLastSeenNow(ctx context.Context, uuid uuid.UUID) error

	SetTokenIssuedNow(ctx context.Context, uuid uuid.UUID) error

	MakePayment(ctx context.Context, uuid uuid.UUID, byTime time.Time) error

	DesyncUnpaid(ctx context.Context) error
}

type service struct {
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
) Service {
	return &service{userRepo, addressRepo, tm, clock}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (data.User, error) {
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

	if err = s.userRepo.Create(ctx, &user); err != nil {
		return data.User{}, err
	}

	return user, nil
}

func (s *service) Get(ctx context.Context, uuid uuid.UUID) (data.User, error) {
	return s.userRepo.GetByUUID(ctx, uuid)
}

func (s *service) GetAll(ctx context.Context) ([]data.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *service) Update(ctx context.Context, uuid uuid.UUID, req UpdateRequest) (data.User, error) {
	var user data.User
	err := s.tm.InTransaction(ctx, func(ctx context.Context) error {
		var err error

		user, err = s.userRepo.GetByUUIDLocked(ctx, uuid)
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

		if err = s.userRepo.Save(ctx, user); err != nil {
			return err
		}

		if req.IsBanned != nil || req.PrivateKey != nil {
			panic("TODO: implement")
			// s.addressRepo.MarkOwnedByDesynced()
			// s.desyncNotifier.Notify()
		}

		return nil
	})
	return user, err
}

func (s *service) update(ctx context.Context, uuid uuid.UUID, patch data.UserPatch) error {
	return s.tm.InTransaction(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.GetByUUIDLocked(ctx, uuid)
		if err != nil {
			return err
		}

		if err = user.Update(patch); err != nil {
			return err
		}

		return s.userRepo.Save(ctx, user)
	})
}

func (s *service) SetLastSeenNow(ctx context.Context, uuid uuid.UUID) error {
	now := s.clock.Now()
	return s.update(ctx, uuid, data.UserPatch{LastSeenAt: &now})
}

func (s *service) SetTokenIssuedNow(ctx context.Context, uuid uuid.UUID) error {
	now := s.clock.Now()
	return s.update(ctx, uuid, data.UserPatch{TokenIssuedAt: &now})
}

func (s *service) MakePayment(ctx context.Context, uuid uuid.UUID, byTime time.Time) error {
	return s.update(ctx, uuid, data.UserPatch{PaidByTime: &byTime})
}

func (s *service) DesyncUnpaid(ctx context.Context) error {
	panic("TODO: implement")
}
