package common

const (
	// HashLength hash length 32 byte
	HashLength = 32
)

// Hash common type for hash
type Hash [HashLength]byte

// Equal judge two hash type is equal
func (h Hash) Equal(b Hash) bool {
	for i := range h {
		if h[i] != b[i] {
			return false
		}
	}
	return true
}
