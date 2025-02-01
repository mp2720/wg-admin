package repos

import (
	"context"
	"mp2720/wg-admin/wg-admin/storage/data"
)

type UserRepo interface {
	// Creates new user. User data given by pointer will be updated.
	//
	// Errors: ErrAlreadyExists.
	Create(ctx context.Context, user *data.User) error

	// Get user by name.
	//
	// Errors: ErrNotFound.
	GetByName(ctx context.Context, name string) (data.User, error)

	// Get all users.
	GetAll(ctx context.Context) ([]data.User, error)

	// Update the user by name.
	//
	// Errors: ErrNotFound.
	Update(
		ctx context.Context,
		name string,
		updateFn func(ctx context.Context, user *data.User) error,
	) (data.User, error)

	// Delete the user by name.
	//
	// Errors: ErrNotFound.
	Delete(ctx context.Context, name string) error
}
