package leveldb

import (
	"seth/database"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	levelDBName                string = "leveldb"
	defaultFilterBits          int    = 10
	defaultMincache            int    = 16
	defaultMinOpenFilesHandles int    = 16
)

type levelDB struct {
	db *leveldb.DB
}

func init() {
	database.Register(&levelDB{})
}

// NewLevelDB new level database
func NewLevelDB(path string, cache int, handles int) (database.Database, error) {

	result := levelDB{}
	err := result.Open(path, cache, handles)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Name retuen regitser name
func (db *levelDB) Name() string {
	return levelDBName
}

// Open open level database
func (db *levelDB) Open(path string, cache int, handles int) error {

	// Ensure we have some minimal caching and file guarantees
	if cache < defaultMincache {
		cache = defaultMincache
	}
	if handles < defaultMinOpenFilesHandles {
		handles = defaultMinOpenFilesHandles
	}

	// ref https://godoc.org/github.com/btcsuitereleases/goleveldb/leveldb/opt

	ldb, err := leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: handles, // The default value is 500.
		// Two of these are used internally
		BlockCacheCapacity: cache / 2 * opt.MiB,                      // The default value is 8MiB.
		WriteBuffer:        cache / 4 * opt.MiB,                      // The default value is 4MiB.
		Filter:             filter.NewBloomFilter(defaultFilterBits), // default is nil
	})

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		ldb, err = leveldb.RecoverFile(path, nil)
	}

	if err != nil {
		return err
	}

	db.db = ldb

	return nil
}

// Close don't forget close db when not use
func (db *levelDB) Close() {
	if db.db != nil {
		db.db.Close()
	}
}

// Get gets the value for the given key
func (db *levelDB) GetString(key string) (string, error) {
	value, err := db.Get([]byte(key))

	return string(value), err
}

// Get gets the value for the given key
func (db *levelDB) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

// Put sets the value for the given key
func (db *levelDB) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// Put sets the value for the given key
func (db *levelDB) PutString(key string, value string) error {
	return db.Put([]byte(key), []byte(value))
}

// Has returns true if the DB does contains the given key.
func (db *levelDB) Has(key []byte) (ret bool, err error) {
	return db.db.Has(key, nil)
}

// Has returns true if the DB does contains the given key.
func (db *levelDB) HasString(key string) (ret bool, err error) {
	return db.Has([]byte(key))
}

// Delete deletes the value for the given key.
func (db *levelDB) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

// Delete deletes the value for the given key.
func (db *levelDB) DeleteSring(key string) error {
	return db.Delete([]byte(key))
}

// NewBatch new a batch operator
func (db *levelDB) NewBatch() database.Batch {
	batch := &Batch{
		db:      db,
		leveldb: db.db,
		batch:   new(leveldb.Batch),
	}
	return batch
}
