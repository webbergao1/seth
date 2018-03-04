package types

import (
	"errors"
	"math/big"
	"seth/common"
	"seth/crypto"
)

const (
	// magicNumberForV magic number for V
	magicNumberForV = 35
)

var (
	// ErrInvalidChainID a error for invalid chain id
	ErrInvalidChainID = errors.New("invalid chain id for signer")
	// ErrInvalidSig a error for invalid signature value
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	// ErrInvalidPublicKey error for invalid public key
	ErrInvalidPublicKey = errors.New("invalid public key")

	bigMagicNumberForV = big.NewInt(magicNumberForV)
)

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)
	// SignatureValues returns the *crypto.Signature values corresponding to the
	// given signature.
	SignatureValues(sign *crypto.Signature) *crypto.Signature
	// Hash returns the hash to be signed.
	Hash(tx *Transaction) common.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

// SethSigner implements Signer using the eth EIP155 rules.
type SethSigner struct {
	chainID, chainIDMul *big.Int
}

// NewSethSigner new seth signer
func NewSethSigner(chainID *big.Int) SethSigner {
	if chainID == nil {
		chainID = new(big.Int)
	}
	return SethSigner{
		chainID:    chainID,
		chainIDMul: new(big.Int).Mul(chainID, big.NewInt(2)),
	}
}

// Equal returns true if the given signer is the same as the receiver.
func (ss SethSigner) Equal(s2 Signer) bool {
	seth, ok := s2.(SethSigner)
	return ok && seth.chainID.Cmp(ss.chainID) == 0
}

// Sender returns the sender address of the transaction.
func (ss SethSigner) Sender(tx *Transaction) (common.Address, error) {

	if tx.ChainID().Cmp(ss.chainID) != 0 {
		return common.Address{}, ErrInvalidChainID
	}
	R, S, V := tx.Data.Signature.RSV()
	if ss.chainID.Sign() != 0 {
		V.Sub(V, ss.chainIDMul)
		V.Sub(V, bigMagicNumberForV)
	}
	return recoverPlain(ss.Hash(tx), R, S, V)
}

// Hash returns the hash to be signed.
func (ss SethSigner) Hash(tx *Transaction) common.Hash {
	return crypto.RlpHash([]interface{}{
		tx.Data.AccountNonce,
		tx.Data.To,
		tx.Data.Amount,
		ss.chainID,
	})
}

func recoverPlain(sighash common.Hash, R, S, V *big.Int) (common.Address, error) {
	if V.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	v := byte(V.Uint64())
	if !crypto.ValidateSignatureValues(v, R, S) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = v
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, ErrInvalidPublicKey
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// SignatureValues returns the crypto.Signature values corresponding to the
// given signature.
func (ss SethSigner) SignatureValues(sign *crypto.Signature) *crypto.Signature {
	R, S, V := sign.RSV()
	if ss.chainID.Sign() != 0 {
		V = big.NewInt(int64(sign[64] + magicNumberForV))
		V.Add(V, ss.chainIDMul)
	}
	r, s := R.Bytes(), S.Bytes()
	sig := new(crypto.Signature)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = byte(V.Uint64())
	return sig
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		return new(big.Int).SetUint64((v - magicNumberForV) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(magicNumberForV))
	return v.Div(v, big.NewInt(2))
}
