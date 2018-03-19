package state

import (
	"math/big"
	"seth/common"
	"seth/database"
)

// Statedb use to store accout with the merkle trie
type Statedb struct {
}

// NewStatedb new a statedb
func NewStatedb(root common.Hash, db database.Database) *Statedb {
	return &Statedb{}
}

// GetAmount get amount of account
func (s *Statedb) GetAmount(addr common.Address) *big.Int {
	return common.Big0
}

// SetAmount set amount of account
func (s *Statedb) SetAmount(addr common.Address, amount *big.Int) {

}

// AddAmount add amount for account
func (s *Statedb) AddAmount(addr common.Address, amount *big.Int) {

}

// SubAmount sub amount for account
func (s *Statedb) SubAmount(addr common.Address, amount *big.Int) {

}

// GetNonce get nonce of account
func (s *Statedb) GetNonce(addr common.Address) uint64 {
	return 0
}

// SetNonce set nonce of account
func (s *Statedb) SetNonce(addr common.Address, nonce uint64) {

}

// Commit commit
func (s *Statedb) Commit() (root common.Hash, err error) {
	return common.Hash{}, nil
}