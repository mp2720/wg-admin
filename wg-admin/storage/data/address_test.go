package data_test

import (
	"mp2720/wg-admin/wg-admin/storage/data"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewAddressInNetwork(t *testing.T) {
	netV4 := netip.MustParsePrefix("1.2.3.0/24")
	netV6 := netip.MustParsePrefix("1:2::0/120")

	addr := data.Address{
		HostID:  10,
		Version: data.IPv6,
		Name:    "some free ipv6 address",
		Owner:   nil,
	}

	addrInNet, err := data.NewAddressInNetwork(addr, netV4, netV6)
	require.NoError(t, err)
	require.Equal(t, data.AddressInNetwork{
		Address:    addr,
		NetAddress: netip.MustParseAddr("1:2::A"),
		Net:        netV6,
	}, addrInNet)

	addr.Version = data.IPv4
	addrInNet, err = data.NewAddressInNetwork(addr, netV4, netV6)
	require.NoError(t, err)
	require.Equal(t, data.AddressInNetwork{
		Address:    addr,
		NetAddress: netip.MustParseAddr("1.2.3.10"),
		Net:        netV4,
	}, addrInNet)

	addr.HostID = 256
	addrInNet, err = data.NewAddressInNetwork(addr, netV4, netV6)
	require.ErrorIs(t, err, data.ErrAddressIsOutOfNetworkRange)
}
