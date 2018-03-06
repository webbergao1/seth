package common

import "encoding/hex"

const (
	// HashLength hash length 32 byte
	HashLength = 32
)

// Hash common type for hash
type Hash [HashLength]byte

// BytesToHash byte to hash
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// Equal judge two hash type is equal
func (h Hash) Equal(b Hash) bool {
	for i := range h {
		if h[i] != b[i] {
			return false
		}
	}
	return true
}

// SetBytes Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Bytes return hash bytes
func (h Hash) Bytes() []byte { return h[:] }

// Hex return hash hex string
func (h Hash) Hex() string {
	enc := make([]byte, len(h)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], h[:])
	return string(enc)
}
