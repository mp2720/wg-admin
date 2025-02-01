package repos

import (
	"context"
	"mp2720/wg-admin/wg-admin/storage/data"
	"time"
)

type AddressRepo interface {
	// Get addresses count (including the reserved ones) for given version.
	GetCount(ctx context.Context, version data.IPVersion) (int64, error)

	// Find unassigned address in repo.
	//
	// Returns nil if there's no one.
	FindFree(ctx context.Context, version data.IPVersion) (*data.Address, error)

	// Insert new address.
	// Owner must be set.
	// HostID will be updated on success.
	//
	// Returns ErrNotFound if owner doesn't exist.
	Create(ctx context.Context, address *data.Address) error

	// Update an address owned by user.
	//
	// Returns ErrNotFound if address doesn't exist.
	Update(
		ctx context.Context,
		hostID int64,
		updateFn func(ctx context.Context, address *data.Address) error,
	) (data.Address, error)

	// Get one address by host id.
	//
	// Returns ErrNotFound if address doesn't exist.
	Get(ctx *context.Context, hostID int64) (data.Address, error)

	// Get all addresses owned by user.
	//
	// Returns empty slice if owner was not found.
	GetAllOwnedBy(ctx context.Context, owner *data.User) ([]data.Address, error)

	// Get all addresses of the given version desynchronized not before provided time with host ID
	// not greater than given value.
	GetDesynced(
		ctx context.Context,
		notBefore time.Time,
		version data.IPVersion,
		maxHostID int64,
	) ([]data.Address, error)
}
