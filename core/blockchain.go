package core

import (
	"errors"
	"seth/common"
	"seth/core/types"
	"seth/database"

	"github.com/hashicorp/golang-lru"
)

const (
	blockCacheLimit = 256
)

var (
	ErrNoGenesis = errors.New("Genesis not found in chain")
)

// BlockChain block chain
type BlockChain struct {
	db database.Database

	genesisBlock *types.Block

	blockCache *lru.Cache // Cache for the most recent entire blocks
}

// NewBlockChain new block chain
func NewBlockChain(db database.Database) (*BlockChain, error) {
	bc := &BlockChain{
		db: db,
	}

	bc.blockCache, _ = lru.New(blockCacheLimit)
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		return nil, ErrNoGenesis
	}
	return bc, nil
}

// GetBlockByNumber get block by number
func (bc *BlockChain) GetBlockByNumber(number uint64) *types.Block {
	hash := GetCanonicalHash(bc.db, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash, number)
}

// GetBlock get block by hash & number
func (bc *BlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	if block, ok := bc.blockCache.Get(hash); ok {
		return block.(*types.Block)
	}
	block := GetBlock(bc.db, hash, number)
	if block == nil {
		return nil
	}
	// Cache the found block for next time and return
	bc.blockCache.Add(block.Hash(), block)
	return block
}
