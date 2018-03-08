package core

import (
	"errors"
	"seth/common"
	"seth/config"
	"seth/core/types"
	"sync"
)

var (
	errTxHashExists = errors.New("transaction hash already exists")
	errTxPoolFull   = errors.New("transaction pool is full")
)

// TxPool transaction pool
type TxPool struct {
	mutex   sync.RWMutex
	pending map[common.Hash]*types.Transaction
	signer  types.Signer
}

// NewTxPool new Tx Pool
func NewTxPool() *TxPool {
	pool := &TxPool{
		signer:  types.NewSethSigner(config.Config.ChainID),
		pending: make(map[common.Hash]*types.Transaction),
	}

	return pool
}

// AddTx add transaction to pool
func (pool *TxPool) AddTx(tx *types.Transaction) error {
	if tx == nil {
		return nil
	}
	err := pool.validateTx(tx)
	if err != nil {
		return err
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	hash := tx.Hash()
	if pool.pending[hash] != nil {
		return errTxHashExists
	}
	pool.pending[hash] = tx
	return nil
}

func (pool *TxPool) validateTx(tx *types.Transaction) error {
	return nil
}
