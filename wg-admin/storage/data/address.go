package data

import (
	"errors"
	"mp2720/wg-admin/wg-admin/utils"
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
	// hostId = address & wildcard.
	// Unique for all addresses with the same version.
	// Non-negative.
	// 0, 1 and wildcard (only for v4) is reserved.
	HostID  int64
	Version IPVersion
	Name    string
	Owner   *User
}

func (a *Address) Update(newName string) {
	a.Name = newName
}

type AddressInNetwork struct {
	Address
	NetAddress netip.Addr // should be unique for all addresses within the same net.
	Net        netip.Prefix
}

var ErrAddressIsOutOfNetworkRange = errors.New("Address is out of the network range")

// Get address as a part of the network.
// netV4 or netV6 is chosen by looking at the given address version.
//
// Errors: [ErrAddressIsOutOfNetworkRange]
func NewAddressInNetwork(
	addr Address,
	netV4 netip.Prefix,
	netV6 netip.Prefix,
) (AddressInNetwork, error) {
	var net netip.Prefix
	if addr.Version == IPv4 {
		net = netV4
	} else {
		net = netV6
	}

	netAddress, ok := utils.AddrForHostInNet(uint64(addr.HostID), net)
	if !ok {
		return AddressInNetwork{}, ErrAddressIsOutOfNetworkRange
	}

	return AddressInNetwork{
		Address:    addr,
		NetAddress: netAddress,
		Net:        net,
	}, nil
}
