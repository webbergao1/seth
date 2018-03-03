package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
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

	sig[64] += 27
	if compressed {
		sig[64] += 4
	}
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
