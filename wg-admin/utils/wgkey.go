package utils

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

// Panics on error. Use in tests.
func MustGenerateWireguardPrivateKey() wgtypes.Key {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	return key
}
