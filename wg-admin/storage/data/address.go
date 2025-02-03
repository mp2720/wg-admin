package data

import (
	"net/netip"
	"time"
)

type IPVersion int

const (
	IPv4 IPVersion = iota
	IPv6
)

// Do not edit fields directly, call the Update method instead.
type Address struct {
	ID          int64
	HostID      int64
	Address     netip.Addr
	Description *string
	Owner       *User
	DesyncedAt  time.Time
}

type AddressPatch struct {
	Description **string
	Owner       **User
	DesyncedAt  *time.Time
}

func (a *Address) Update(patch AddressPatch) {
	addrCopy := *a

	if patch.Description != nil {
		addrCopy.Description = *patch.Description
	}
	if patch.Owner != nil {
		addrCopy.Owner = *patch.Owner
	}
	if patch.DesyncedAt != nil {
		addrCopy.DesyncedAt = *patch.DesyncedAt
	}

	*a = addrCopy
}
