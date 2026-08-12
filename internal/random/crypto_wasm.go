//go:build js && wasm

package random

import (
	"syscall/js"
)

// CryptoSource uses the browser Web Crypto API for secure randomness.
type CryptoSource struct{}

// NewCryptoSource creates a production random source backed by crypto.getRandomValues.
func NewCryptoSource() (*CryptoSource, error) {
	crypto := js.Global().Get("crypto")
	if crypto.IsUndefined() || crypto.IsNull() {
		return nil, ErrRandomSourceFailure
	}
	fn := crypto.Get("getRandomValues")
	if fn.IsUndefined() || fn.IsNull() {
		return nil, ErrRandomSourceFailure
	}
	return &CryptoSource{}, nil
}

func (c *CryptoSource) RandomInt(max int) (int, error) {
	if max <= 0 {
		return 0, ErrRandomSourceFailure
	}

	crypto := js.Global().Get("crypto")
	buf := js.Global().Get("Uint32Array").New(1)
	limit := uint32(0xFFFFFFFF - (0xFFFFFFFF % uint32(max)))

	for {
		crypto.Call("getRandomValues", buf)
		val := uint32(buf.Index(0).Int())
		if val < limit {
			return int(val % uint32(max)), nil
		}
	}
}
