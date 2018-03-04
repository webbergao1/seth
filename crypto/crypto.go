package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"seth/common"
	"seth/common/math"
	"seth/crypto/secp256k1"
	"seth/crypto/sha3"
	"seth/rlp"
)

const (
	// SignatureSize represents the signature length
	SignatureSize = 65
)

var (
	emptySignature = Signature{}
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

type (
	// PrivateKey represents the ecdsa privatekey
	PrivateKey ecdsa.PrivateKey
	// PublicKey represents the ecdsa publickey
	PublicKey ecdsa.PublicKey
	// Signature represents the ecdsa_signcompact signature
	// data format [r - s - v]
	Signature [SignatureSize]byte
)

// Keccak256 Calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// ToECDSAPub Acquire ecdsa publickey from bytes
func ToECDSAPub(pub []byte) *ecdsa.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(S256(), pub)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}
}

// FromECDSAPub Acquire bytes from ECDSA publickey
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

// FromECDSA Exports a private key into a binary dump.
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// S256 Returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}

// PubkeyToAddress Get address from public key
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := FromECDSAPub(&p)
	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

// RlpHash rlp encode content & sum hash
func RlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if s.Cmp(secp256k1halfN) > 0 {
		return false
	}

	return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
}

// GenerateKey returns a random PrivateKey
func GenerateKey() (*PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(priv), err
}

// GetBytes get the actual bytes of ecdsa privatekey
func (priv *PrivateKey) GetBytes() []byte {
	return FromECDSA((*ecdsa.PrivateKey)(priv))
}

// GetPublicKey return the public key
func (priv *PrivateKey) GetPublicKey() *PublicKey {
	return (*PublicKey)(&priv.PublicKey)
}

// Sign signs the hash and returns the signature
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func (priv *PrivateKey) Sign(hash []byte) (sig *Signature, err error) {
	if len(hash) != 32 {
		return &emptySignature, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	secretKey := priv.GetBytes()
	defer zeroBytes(secretKey)

	rawSig, err := secp256k1.Sign(hash, secretKey)

	sig = new(Signature)
	sig.SetBytes(rawSig, false)
	return sig, err
}

// SetBytes returns a format signature according the raw signature
func (sig *Signature) SetBytes(data []byte, compressed bool) {
	if len(data) == 65 {
		copy(sig[:], data[:])
	}
}

// RSV returns the r s v values
func (sig *Signature) RSV() (v, r, s *big.Int) {
	return new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:64]), big.NewInt(int64(sig[64]))
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
