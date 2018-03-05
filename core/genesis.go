package core

import (
	"math/big"
	"seth/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	// TagMainNetGenesis tag for main net genesis
	TagMainNetGenesis = "mainnet"
	// TagTestNetGenesis tag for test net genesis
	TagTestNetGenesis = "testnet"
	// TagDeveloperNetGenesis tag for developer net genesis
	TagDeveloperNetGenesis = "dev"
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
}

// DefaultGenesis is default main net genesis block info
func DefaultGenesis() *Genesis {
	return &Genesis{
		ChainID:    big.NewInt(1),
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		Difficulty: big.NewInt(17179869184),
	}
}

// TestnetGenesis is test net genesis block info
func TestnetGenesis() *Genesis {
	return &Genesis{
		ChainID:    big.NewInt(0),
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353535"),
		Difficulty: big.NewInt(1048576),
	}
}

// DevelopernetGenesis is test net genesis block info
func DevelopernetGenesis() *Genesis {
	return &Genesis{
		ChainID:    big.NewInt(-1),
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353535"),
		Difficulty: big.NewInt(1048576),
	}
}
