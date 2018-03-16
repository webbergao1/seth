package state

import "math/big"

// Account is a balance model for blockchain
type Account struct {
	Nonce  uint64
	Amount *big.Int
}

// NodeObject is state object for statedb
type NodeObject struct {
	account Account
}
