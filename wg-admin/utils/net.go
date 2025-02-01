package utils

import (
	"encoding/binary"
	"net/netip"
)

// ok = false only if hostID is out of range of the given network.
func AddrForHostInNet(hostID uint64, netPrefix netip.Prefix) (addr netip.Addr, ok bool) {
	if !netPrefix.IsValid() {
		panic("invalid network prefix")
	}

	netPrefix = netPrefix.Masked()

	// check that there's enough space in given network for this host.

	var bitsInAddress int
	if netPrefix.Addr().Is4() {
		bitsInAddress = 32
	} else {
		bitsInAddress = 128
	}
	bitsInHost := bitsInAddress - netPrefix.Bits()

	if (1<<bitsInHost)-1 < hostID {
		return netip.Addr{}, false
	}

	// Apply bitwise OR

	if netPrefix.Addr().Is4() {
		netPrefixBytes := netPrefix.Addr().As4()

		addr := binary.BigEndian.Uint32(netPrefixBytes[:])
		addr = addr | uint32(hostID)

		var addrBytes [4]byte
		binary.BigEndian.PutUint32(addrBytes[:], addr)

		return netip.AddrFrom4(addrBytes), true
	} else {
		netPrefixBytes := netPrefix.Addr().As16()

		addrLo := binary.BigEndian.Uint64(netPrefixBytes[8:])
		addrHi := binary.BigEndian.Uint64(netPrefixBytes[:8])
		addrLo = addrLo | hostID

		var addrBytes [16]byte
		binary.BigEndian.PutUint64(addrBytes[:], addrHi)
		binary.BigEndian.PutUint64(addrBytes[8:], addrLo)

		return netip.AddrFrom16(addrBytes), true
	}
}

// Get wildcard for netPrefix.
//
// Panics if prefix is invalid.
func GetNetWildcard(netPrefix netip.Prefix) (hi uint64, lo uint64) {
	if !netPrefix.IsValid() {
		panic("invalid prefix")
	}

	netPrefix = netPrefix.Masked()

	if netPrefix.Addr().Is4() {
		wildcard := uint64(1)<<(32-netPrefix.Bits()) - 1
		return 0, wildcard
	} else {
		wildcardBits := 128 - netPrefix.Bits()

		var wildcardLo uint64
		var wildcardHi uint64
		if wildcardBits == 128 {
			wildcardLo = 0xffff_ffff_ffff_ffff
			wildcardHi = 0xffff_ffff_ffff_ffff
		} else if wildcardBits >= 64 {
			wildcardLo = 0xffff_ffff_ffff_ffff
			wildcardHi = uint64(1)<<(wildcardBits-64) - 1
		} else {
			wildcardLo = uint64(1)<<(wildcardBits) - 1
			wildcardHi = 0
		}

		return wildcardHi, wildcardLo
	}
}

// Get host ID of address in given network. Returns ok flag and 128-bit uint as high and low parts.
// ok is only set to false if the net does not contain address.
//
// Panic if address or prefix is empty, and if versions of address and netPrefix are different.
func GetHostIDOfAddr(address netip.Addr, netPrefix netip.Prefix) (hi uint64, lo uint64, ok bool) {
	// Will panic if prefix is invalid
	wildcardHi, wildcardLo := GetNetWildcard(netPrefix)

	if address == (netip.Addr{}) {
		panic("empty address")
	}
	if address.Is4() != netPrefix.Addr().Is4() {
		panic("IP version mismatch")
	}

	netPrefix = netPrefix.Masked()

	if !netPrefix.Contains(address) {
		return 0, 0, false
	}

	if address.Is4() {
		addressBytes := address.As4()
		addressUint32 := binary.BigEndian.Uint32(addressBytes[:])

		return 0, uint64(addressUint32) & wildcardLo, true
	} else {
		addressBytes := address.As16()
		addressUint64Lo := binary.BigEndian.Uint64(addressBytes[8:])
		addressUint64Hi := binary.BigEndian.Uint64(addressBytes[:8])

		return wildcardHi & addressUint64Hi, wildcardLo & addressUint64Lo, true
	}
}

// Check that address is broadcast in given network.
// Panics if address or netPrefix is empty, and if versions of address and netPrefix are different.
func AddrIsNetBroadcast(address netip.Addr, netPrefix netip.Prefix) bool {
	// [GetHostIDOfAddr] panics on invalid input.
	_, hostIDLo, ok := GetHostIDOfAddr(address, netPrefix)
	if !ok {
		return false
	}

	netPrefix = netPrefix.Masked()

	if address.Is6() {
		// IPv6 does not define broadcast address
		return false
	}

	// all bits in host part must be 1
	wildcard := uint64(1)<<(32-netPrefix.Bits()) - 1

	return hostIDLo == wildcard
}
