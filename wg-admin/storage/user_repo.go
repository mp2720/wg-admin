package storage

import (
	"context"
	"mp2720/wg-admin/wg-admin/storage/data"

	"github.com/google/uuid"
)

type UserRepo interface {
	// Creates new user. User data given by pointer will be updated.
	//
	// Errors: ErrAlreadyExists.
	Create(ctx context.Context, user *data.User) error

	// Get user by name.
	//
	// Errors: ErrNotFound.
	GetByUUID(ctx context.Context, uuid uuid.UUID) (data.User, error)

	// Get user by name and lock within the transaction.
	//
	// Errors: ErrNotFound.
	GetByUUIDLocked(ctx context.Context, uuid uuid.UUID) (data.User, error)

	// Get all users.
	GetAll(ctx context.Context) ([]data.User, error)

	// Save existing user.
	//
	// Errors: ErrNotFound.
	Save(ctx context.Context, user data.User) error

	// Delete the user by name.
	//
	// Errors: ErrNotFound.
	Delete(ctx context.Context, uuid uuid.UUID) error
}
