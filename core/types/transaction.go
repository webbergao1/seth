package types

import (
	"math/big"
	"seth/common"
	"seth/crypto"
	"sync/atomic"
)

// Transaction transaction struct
type Transaction struct {
	Data *txData `json:"data"`

	// cache
	hash atomic.Value
}

type txData struct {
	From         common.Address    `json:"from"    gencodec:"required"`
	To           *common.Address   `json:"to"       rlp:"nil"` // nil means contract creation
	AccountNonce uint64            `json:"nonce"    gencodec:"required"`
	Amount       *big.Int          `json:"value"    gencodec:"required"`
	Signature    *crypto.Signature `json:"signature"    gencodec:"required"`
}

// NewTransaction creates a new transaction to transfer asset.
func NewTransaction(from, to common.Address, amount *big.Int, nonce uint64) *Transaction {
	txdata := &txData{
		From:         from,
		To:           &to,
		Amount:       amount,
		AccountNonce: nonce,
	}

	return &Transaction{Data: txdata}
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := crypto.RlpHash(tx)
	tx.hash.Store(v)
	return v
}
