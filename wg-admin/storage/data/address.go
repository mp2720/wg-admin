package data

import (
	"errors"
	"net/netip"
)

type IPVersion int

const (
	IPv4 IPVersion = iota
	IPv6
)

// Assigned address belongs to some user.
// Do not edit fields directly, call the Update method instead.
type Address struct {
	HostID  int64 // hostId = address & wildcard; unique for all addresses with the same version
	Version IPVersion
	Name    string
	Owner   *User
}

func (a *Address) Update(newName string) {
	a.Name = newName
}

type AssignedAddressInNetwork struct {
	Address
	NetAddress netip.Addr // should be unique for all addresses within the same net.
	Net        netip.Prefix
}

var (
	ErrAddressIsOutOfNetworkRange = errors.New("Address is out of the network range")
)

func NewAssignedAddressInNetwork(
	a Address,
	net netip.Prefix,
) (AssignedAddressInNetwork, error) {

}
