package types

import (
	"encoding/binary"
	"math/big"
	"seth/common"
	"seth/crypto"
	"sync/atomic"
	"time"
)

// BlockNonce Block nonce
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// Header Block Header
type Header struct {
	ParentHash common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase   common.Address `json:"miner"            gencodec:"required"`
	Root       common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash     common.Hash    `json:"transactionsRoot" gencodec:"required"`
	Difficulty *big.Int       `json:"difficulty"       gencodec:"required"`
	Number     *big.Int       `json:"number"           gencodec:"required"`
	Time       *big.Int       `json:"timestamp"        gencodec:"required"`
	Extra      []byte         `json:"extraData"        gencodec:"required"`
	MixDigest  common.Hash    `json:"mixHash"          gencodec:"required"`
	Nonce      BlockNonce     `json:"nonce"            gencodec:"required"`
}

// Clone  clone block header
func (h Header) Clone() *Header {
	clone := h
	if clone.Time = new(big.Int); h.Time != nil {
		clone.Time.Set(h.Time)
	}
	if clone.Difficulty = new(big.Int); h.Difficulty != nil {
		clone.Difficulty.Set(h.Difficulty)
	}
	if clone.Number = new(big.Int); h.Number != nil {
		clone.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		clone.Extra = make([]byte, len(h.Extra))
		copy(clone.Extra, h.Extra)
	}
	return &clone
}

// Hash return the block hash of the heade
func (h *Header) Hash() common.Hash {
	return crypto.RlpHash(h)
}

// Body block body struct
type Body struct {
	Transactions []*Transaction
}

// Block block struct define
type Block struct {
	Header       *Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// NewBlock new block
func NewBlock(header *Header, txs []*Transaction) *Block {
	block := &Block{Header: header.Clone(), td: new(big.Int)}

	if len(txs) == 0 {
		// TODO add calc TxHash by trie
		//block.header.TxHash = EmptyRootHash
	} else {
		// TODO add calc TxHash by trie
		//block.header.TxHash = DeriveSha(Transactions(txs))
		block.transactions = make(Transactions, len(txs))
		copy(block.transactions, txs)
	}

	return block
}

// Hash returns the keccak256 hash of block's header.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
}

// FindTransaction find transaction with hash value
func (b *Block) FindTransaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

// NumberU64 return the number of block
func (b *Block) NumberU64() uint64 { return b.Header.Number.Uint64() }

// Body return the body of block
func (b *Block) Body() *Body { return &Body{b.transactions} }
