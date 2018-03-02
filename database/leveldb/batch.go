package leveldb

import "github.com/syndtr/goleveldb/leveldb"

// Batch batch implenent for leveldb
type Batch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

// Put sets the value for the given key
func (b *Batch) Put(key []byte, value []byte) {

	b.batch.Put(key, value)

}

// Delete deletes the value for the given key.
func (b *Batch) Delete(key []byte) {

	b.batch.Delete(key)

}

// Commit commit batch operator.
func (b *Batch) Commit() error {

	return b.db.Write(b.batch, nil)
}

// Rollback rollback batch operator.
func (b *Batch) Rollback() {

	b.batch.Reset()
}

// Close close the batch operator,just rollback anything not commit.
func (b *Batch) Close() {

	b.batch.Reset()
}
