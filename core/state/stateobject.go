package state

import "math/big"

// Account is a balance model for blockchain
type Account struct {
	Nonce  uint64
	Amount *big.Int
}

// StateObject is state object for statedb
type StateObject struct {
	account Account
}
