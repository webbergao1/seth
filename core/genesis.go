package core

import (
	"encoding/json"
	"errors"
	"math/big"
	"seth/common"
	"seth/core/types"
	"seth/database"
)

const (
	// TagMainNetGenesis tag for main net genesis
	TagMainNetGenesis = "mainnet"
	// TagTestNetGenesis tag for test net genesis
	TagTestNetGenesis = "testnet"
	// TagDeveloperNetGenesis tag for developer net genesis
	TagDeveloperNetGenesis = "dev"
)

var (
	// ErrHashGenesisBlock error has genesis block in blockchain
	ErrHasGenesisBlock = errors.New("Found genesis block in blockchain")
)

// Genesis is genesis struct to
type Genesis struct {
	ChainID    *big.Int       `json:"chainId"`
	Nonce      uint64         `json:"nonce"`
	Timestamp  uint64         `json:"timestamp"`
	ExtraData  []byte         `json:"extraData"`
	Difficulty *big.Int       `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash    `json:"mixHash"`
	Coinbase   common.Address `json:"coinbase"`

	Number     uint64      `json:"number"`
	ParentHash common.Hash `json:"parentHash"`
}

// DefaultGenesis is default main net genesis block info
func DefaultGenesis() *Genesis {
	return &Genesis{
		ChainID:    big.NewInt(1),
		Nonce:      66,
		ExtraData:  []byte("Mainnet Ethereum Genesis Block"),
		Difficulty: big.NewInt(17179869184),
	}
}

// TestnetGenesis is test net genesis block info
func TestnetGenesis() *Genesis {
	return &Genesis{
		ChainID:    big.NewInt(0),
		Nonce:      66,
		ExtraData:  []byte("Testnet Ethereum Genesis Block"),
		Difficulty: big.NewInt(1048576),
	}
}

// DevelopernetGenesis is test net genesis block info
func DevelopernetGenesis() *Genesis {
	return &Genesis{
		ChainID:    big.NewInt(-1),
		Nonce:      66,
		ExtraData:  []byte("Devnet Ethereum Genesis Block"),
		Difficulty: big.NewInt(1048576),
	}
}

// SetupGensisBlock setup genesis block
func (g Genesis) SetupGensisBlock(db database.Database) (common.Hash, error) {
	stored := GetCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		block, err := g.Commit(db)
		return block.Hash(), err
	}
	return stored, ErrHasGenesisBlock
}

// Commit commit genesis block to blockchain
func (g Genesis) Commit(db database.Database) (*types.Block, error) {
	block := g.ToBlock(db)
	if block.Header.Number.Sign() != 0 {
		//return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	batch := db.NewBatch()
	if err := WriteTd(batch, block.Hash(), block.NumberU64(), g.Difficulty); err != nil {
		batch.Rollback()
		return nil, err
	}
	if err := WriteBlock(batch, block); err != nil {
		batch.Rollback()
		return nil, err
	}

	WriteCanonicalHash(batch, block.Hash(), block.NumberU64())
	WriteHeadBlockHash(batch, block.Hash())

	json, err := json.Marshal(g)
	if err != nil {
		batch.Rollback()
		return nil, err
	}

	WriteChainConfig(batch, block.Hash(), json)
	err = batch.Commit()
	return block, err
}

// ToBlock genesis to block
func (g *Genesis) ToBlock(db database.Database) *types.Block {

	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		//Root:       root,
	}

	return types.NewBlock(head, nil)
}
