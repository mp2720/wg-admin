package storage

import (
	"context"
	"errors"
	"fmt"
	"mp2720/wg-admin/wg-admin/storage/data"
	"net/netip"
	"time"
)

type ErrNoFreeAddressInNetwork struct {
	net netip.Prefix
}

func (e ErrNoFreeAddressInNetwork) Error() string {
	return fmt.Sprintf("No free addresses in network %s", e.net.String())
}

var ErrAddressNetsNotSet = errors.New("Address nets are required to be set")

type ErrNetIsTooSmall struct {
	net netip.Prefix
}

func (e ErrNetIsTooSmall) Error() string {
	return fmt.Sprintf("Network %s is too small to be set", e.net.String())
}

type AddressRepoConfigurator interface {
	// Errors: [ErrNetIsTooSmall]
	SetNets(ctx context.Context, v4Net netip.Prefix, v6Net netip.Prefix) error

	// Returns repo with configured nets.
	// If no nets were set, then [ErrAddressNetsNotSet] is returned.
	GetRepo(ctx context.Context) (AddressRepo, error)
}

type AddressRepo interface {
	GetNet(version data.IPVersion) netip.Prefix

	// Errors: [ErrNoFreeAddressInNetwork].
	Allocate(ctx context.Context, version data.IPVersion) (data.Address, error)

	// Returns nil if no address was found.
	Get(ctx context.Context, netAddress netip.Addr) (*data.Address, error)

	// Get address with locking within the current transaction.
	// Returns nil if no address was found.
	GetLocked(ctx context.Context, netAddress netip.Addr) (*data.Address, error)

	GetOwnedBy(ctx context.Context, owner *data.User) ([]data.Address, error)

	// Use within the transaction with GetLocked
	// Errors: [ErrNotFound]
	Save(ctx context.Context, address data.Address) error

	DesyncOwnedBy(ctx context.Context, owner *data.User, desyncTime time.Time) error

	GetDesynced(ctx context.Context, notBefore time.Time) ([]data.Address, error)
}
