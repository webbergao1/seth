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
	hash   atomic.Value
	sender atomic.Value
}

type txData struct {
	To           *common.Address   `json:"to"       rlp:"nil"` // nil means contract creation
	AccountNonce uint64            `json:"nonce"    gencodec:"required"`
	Amount       *big.Int          `json:"value"    gencodec:"required"`
	Signature    *crypto.Signature `json:"signature"    gencodec:"required"`
}

// NewTransaction creates a new transaction to transfer asset.
func NewTransaction(to common.Address, amount *big.Int, nonce uint64) *Transaction {
	txdata := &txData{
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
	v := crypto.RlpHash([]interface{}{
		tx.Data.AccountNonce,
		tx.Data.To,
		tx.Data.Amount,
	})
	tx.hash.Store(v)
	return v
}

// SignTx sign transaction
func (tx *Transaction) SignTx(signer Signer, privatekey *crypto.PrivateKey) error {
	h := signer.Hash(tx)
	sig, err := privatekey.Sign(h[:])
	if err != nil {
		return err
	}
	tx.Data.Signature = signer.SignatureValues(sig)
	return nil
}

// Sender return tx sender
func (tx *Transaction) Sender(signer Signer) (common.Address, error) {
	if sender := tx.sender.Load(); sender != nil {
		return sender.(common.Address), nil
	}
	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.sender.Store(addr)
	return addr, nil
}

// ChainID returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainID() *big.Int {
	_, _, V := tx.Data.Signature.RSV()
	return deriveChainID(V)
}
