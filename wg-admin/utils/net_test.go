package utils_test

import (
	"mp2720/wg-admin/wg-admin/utils"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddrForHostInNet(t *testing.T) {
	require.Panics(t, func() {
		utils.AddrForHostInNet(14, netip.Prefix{})
	})

	// v4
	addr, err := utils.AddrForHostInNet(0, netip.MustParsePrefix("124.123.122.121/32"))
	require.NoError(t, err)
	require.Equal(t, netip.MustParseAddr("124.123.122.121"), addr)

	_, err = utils.AddrForHostInNet(2, netip.MustParsePrefix("124.123.122.121/31"))
	require.ErrorIs(t, err, utils.ErrNetworkIsTooSmall)

	addr, err = utils.AddrForHostInNet(1, netip.MustParsePrefix("124.123.122.126/31"))
	require.NoError(t, err)
	require.Equal(t, netip.MustParseAddr("124.123.122.127"), addr)

	// v6
	addr, err = utils.AddrForHostInNet(
		0,
		netip.MustParsePrefix("1234:5678:9ABC:DEF1:2345:6789:ABCD:EF13/128"),
	)
	require.NoError(t, err)
	require.Equal(t, netip.MustParseAddr("1234:5678:9ABC:DEF1:2345:6789:ABCD:EF13"), addr)

	_, err = utils.AddrForHostInNet(
		2,
		netip.MustParsePrefix("1234:5678:9ABC:DEF1:2345:6789:ABCD:EF13/127"),
	)
	require.ErrorIs(t, err, utils.ErrNetworkIsTooSmall)

	addr, err = utils.AddrForHostInNet(
		1,
		netip.MustParsePrefix("1234:5678:9ABC:DEF1:2345:6789:ABCD:EF14/126"),
	)
	require.NoError(t, err)
	require.Equal(t, netip.MustParseAddr("1234:5678:9ABC:DEF1:2345:6789:ABCD:EF15"), addr)
}

func Test_GetHostIDOfAddress(t *testing.T) {
	require.Panics(t, func() {
		utils.GetHostIDOfAddr(netip.Addr{}, netip.MustParsePrefix("0.0.0.0/0"))
	})

	require.Panics(t, func() {
		utils.GetHostIDOfAddr(netip.MustParseAddr("0.0.0.0"), netip.Prefix{})
	})

	require.Panics(t, func() {
		utils.GetHostIDOfAddr(netip.MustParseAddr("0.0.0.0"), netip.MustParsePrefix("::/0"))
	})

	// v4

	_, _, ok := utils.GetHostIDOfAddr(
		netip.MustParseAddr("0.0.0.255"),
		netip.MustParsePrefix("0.0.0.0/30"),
	)
	require.False(t, ok)

	hi, lo, ok := utils.GetHostIDOfAddr(
		netip.MustParseAddr("12.34.56.78"),
		netip.MustParsePrefix("12.34.48.0/20"),
	)
	require.True(t, ok)
	require.Zero(t, hi)
	require.Equal(t, uint64(2126), lo)

	hi, lo, ok = utils.GetHostIDOfAddr(
		netip.MustParseAddr("12.34.56.78"),
		netip.MustParsePrefix("0.0.0.0/0"),
	)
	require.True(t, ok)
	require.Zero(t, hi)
	require.Equal(t, uint64(203569230), lo)

	// v6

	_, _, ok = utils.GetHostIDOfAddr(
		netip.MustParseAddr("ffff::5"),
		netip.MustParsePrefix("ffff::2/127"),
	)
	require.False(t, ok)

	hi, lo, ok = utils.GetHostIDOfAddr(
		netip.MustParseAddr("ffff::3"),
		netip.MustParsePrefix("ffff::2/127"),
	)
	require.True(t, ok)
	require.Equal(t, uint64(0), hi)
	require.Equal(t, uint64(1), lo)

	hi, lo, ok = utils.GetHostIDOfAddr(
		netip.MustParseAddr("1234:5678:9ABC:DEF1:2345:5678:9ABC:DEF1"),
		netip.MustParsePrefix("1234:5678:9ABC:DEF1::/64"),
	)
	require.True(t, ok)
	require.Equal(t, uint64(0), hi)
	require.Equal(t, uint64(0x234556789ABCDEF1), lo)

	hi, lo, ok = utils.GetHostIDOfAddr(
		netip.MustParseAddr("1234:5678:9ABC:DEF1:2345:5678:9ABC:DEF1"),
		netip.MustParsePrefix("::/0"),
	)
	require.True(t, ok)
	require.Equal(t, uint64(0x123456789ABCDEF1), hi)
	require.Equal(t, uint64(0x234556789ABCDEF1), lo)
}

func Test_AddrIsNetBroadCast(t *testing.T) {
	require.Panics(t, func() {
		utils.AddrIsNetBroadcast(netip.Addr{}, netip.MustParsePrefix("0.0.0.0/0"))
	})

	require.Panics(t, func() {
		utils.AddrIsNetBroadcast(netip.MustParseAddr("0.0.0.0"), netip.Prefix{})
	})

	require.Panics(t, func() {
		utils.AddrIsNetBroadcast(netip.MustParseAddr("0.0.0.0"), netip.MustParsePrefix("::/0"))
	})

	require.False(t, utils.AddrIsNetBroadcast(
		netip.MustParseAddr("::ffff"),
		netip.MustParsePrefix("::/112"),
	))

	require.True(t, utils.AddrIsNetBroadcast(
		netip.MustParseAddr("0.0.0.1"),
		netip.MustParsePrefix("0.0.0.0/31"),
	))

	require.True(t, utils.AddrIsNetBroadcast(
		netip.MustParseAddr("10.8.0.255"),
		netip.MustParsePrefix("10.8.0.0/24"),
	))

	require.True(t, utils.AddrIsNetBroadcast(
		netip.MustParseAddr("255.255.255.255"),
		netip.MustParsePrefix("0.0.0.0/0"),
	))

	require.False(t, utils.AddrIsNetBroadcast(
		netip.MustParseAddr("10.8.1.0"),
		netip.MustParsePrefix("10.8.0.0/24"),
	))
}
