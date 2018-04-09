package state

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"seth/common"
	"seth/database"
	"seth/database/leveldb"
	"testing"
)

func newTestStateDB() (database.Database, func()) {
	dir, err := ioutil.TempDir("", "teststatedb")
	if err != nil {
		panic(err)
	}
	db, err := leveldb.NewLevelDB(dir, 0, 0)
	if err != nil {
		panic(err)
	}
	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

func Test_Statedb_Operate(t *testing.T) {
	db, remove := newTestStateDB()
	defer remove()

	statedb, err := NewStatedb(common.Hash{}, db)
	if err != nil {
		panic(err)
	}
	for i := byte(0); i < 255; i++ {
		statedb.AddAmount(common.BytesToAddress([]byte{i}), big.NewInt(2*int64(i)))
		statedb.SetNonce(common.BytesToAddress([]byte{i}), 1)
	}

	batch := db.NewBatch()
	hash, err := statedb.Commit(batch)
	if err != nil {
		panic(err)
	}
	batch.Commit()

	statedb, err = NewStatedb(hash, db)
	if err != nil {
		panic(err)
	}
	for i := byte(0); i < 255; i++ {
		amount := statedb.GetAmount(common.BytesToAddress([]byte{i}))
		nonce := statedb.GetNonce(common.BytesToAddress([]byte{i}))
		if amount.Cmp(big.NewInt(2*int64(i))) != 0 {
			panic(fmt.Errorf("error anount amount %d", i))
		}
		if nonce != 1 {
			panic(fmt.Errorf("error anount nonce %d", i))
		}
	}
}

func Test_Statedb_Reset(t *testing.T) {
	db, remove := newTestStateDB()
	defer remove()

	statedb, err := NewStatedb(common.Hash{}, db)
	if err != nil {
		panic(err)
	}
	for i := byte(0); i < 255; i++ {
		statedb.AddAmount(common.BytesToAddress([]byte{i}), big.NewInt(3*int64(i)))
		statedb.SetNonce(common.BytesToAddress([]byte{i}), uint64(i))
	}

	batch := db.NewBatch()
	hash, err := statedb.Commit(batch)
	if err != nil {
		panic(err)
	}
	batch.Commit()

	err = statedb.ResetStatedb(hash, db)
	if err != nil {
		panic(err)
	}

	for i := byte(0); i < 255; i++ {
		amount := statedb.GetAmount(common.BytesToAddress([]byte{i}))
		nonce := statedb.GetNonce(common.BytesToAddress([]byte{i}))
		if amount.Cmp(big.NewInt(3*int64(i))) != 0 {
			panic(fmt.Errorf("error anount amount %d", i))
		}
		if nonce != uint64(i) {
			panic(fmt.Errorf("error anount nonce %d", i))
		}
	}

}
