package core

import (
	"encoding/binary"
	"math/big"
	"seth/common"
	"seth/core/types"
	"seth/database"
	"seth/rlp"
)

var (
	headBlockKey = []byte("LastBlock")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`).
	headerPrefix    = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	tdSuffix        = []byte("t") // headerPrefix + num (uint64 big endian) + hash + tdSuffix -> td
	numSuffix       = []byte("n") // headerPrefix + num (uint64 big endian) + numSuffix -> hash
	blockHashPrefix = []byte("H") // blockHashPrefix + hash -> num (uint64 big endian)
	bodyPrefix      = []byte("b") // bodyPrefix + num (uint64 big endian) + hash -> block body

	configPrefix = []byte("seth-config-") // config prefix for the db
)

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// GetCanonicalHash retrieves a hash assigned to a canonical block number.
func GetCanonicalHash(db database.Database, number uint64) common.Hash {
	data, _ := db.Get(append(append(headerPrefix, encodeBlockNumber(number)...), numSuffix...))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteCanonicalHash write a hash with canonical block number
func WriteCanonicalHash(batch database.Batch, hash common.Hash, number uint64) {
	key := append(append(headerPrefix, encodeBlockNumber(number)...), numSuffix...)
	batch.Put(key, hash.Bytes())
}

// WriteTd write total difficulty of block
func WriteTd(batch database.Batch, hash common.Hash, number uint64, td *big.Int) error {
	data, err := rlp.EncodeToBytes(td)
	if err != nil {
		return err
	}
	key := append(append(append(headerPrefix, encodeBlockNumber(number)...), hash.Bytes()...), tdSuffix...)
	batch.Put(key, data)
	return nil
}

// WriteBlock write block
func WriteBlock(batch database.Batch, block *types.Block) error {
	if err := WriteBody(batch, block.Hash(), block.NumberU64(), block.Body()); err != nil {
		return err
	}
	return WriteHeader(batch, block.Header)

}

// WriteBody write block body
func WriteBody(batch database.Batch, hash common.Hash, number uint64, body *types.Body) error {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		return err
	}

	key := append(append(bodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
	batch.Put(key, data)
	return nil
}

// WriteHeader write block header
func WriteHeader(batch database.Batch, header *types.Header) error {
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	hash := header.Hash().Bytes()
	num := header.Number.Uint64()
	encNum := encodeBlockNumber(num)
	key := append(blockHashPrefix, hash...)
	batch.Put(key, encNum)

	key = append(append(headerPrefix, encNum...), hash...)
	batch.Put(key, data)

	return nil
}

// WriteHeadBlockHash write last block hash
func WriteHeadBlockHash(batch database.Batch, hash common.Hash) {
	batch.Put(headBlockKey, hash.Bytes())
}

// WriteChainConfig write chain config to db
func WriteChainConfig(batch database.Batch, hash common.Hash, jsoncfg []byte) {
	batch.Put(append(configPrefix, hash[:]...), jsoncfg)
}
