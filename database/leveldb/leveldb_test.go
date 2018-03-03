package leveldb

import (
	"bytes"
	"io/ioutil"
	"os"
	"seth/database"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func prepareDbFolder(pathRoot string, subDir string) string {
	dir, err := ioutil.TempDir(pathRoot, subDir)
	if err != nil {
		panic(err)
	}

	return dir
}

func newDbInstance(dbPath string) database.Database {
	db, err := NewLevelDB(dbPath, 16, 16)
	if err != nil {
		panic(err)
	}

	return db
}

func newTestLevelDB() (database.Database, func()) {
	dir := prepareDbFolder("", "leveldbtest")
	db := newDbInstance(dir)
	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

var testvalues = []string{"", "a", "1251", "\x00123\x00"}

func testPutGet(db database.Database, t *testing.T) {

	t.Parallel()

	for _, v := range testvalues {
		err := db.Put([]byte(v), []byte(v))
		if err != nil {
			t.Fatalf("put failed: %v", err)
		}
	}

	for _, v := range testvalues {
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte(v)) {
			t.Fatalf("get returned wrong result, got %q expected %q", string(data), v)
		}
	}

	for _, v := range testvalues {
		err := db.Put([]byte(v), []byte("?"))
		if err != nil {
			t.Fatalf("put override failed: %v", err)
		}
	}

	for _, v := range testvalues {
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte("?")) {
			t.Fatalf("get returned wrong result, got %q expected ?", string(data))
		}
	}

	for _, v := range testvalues {
		orig, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		orig[0] = byte(0xff)
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte("?")) {
			t.Fatalf("get returned wrong result, got %q expected ?", string(data))
		}
	}

	for _, v := range testvalues {
		err := db.Delete([]byte(v))
		if err != nil {
			t.Fatalf("delete %q failed: %v", v, err)
		}
	}

	for _, v := range testvalues {
		_, err := db.Get([]byte(v))
		if err == nil {
			t.Fatalf("got deleted value %q", v)
		}
	}
}

func Test_levelDB_PutGet(t *testing.T) {
	// Init levelDB
	db, remove := newTestLevelDB()
	defer remove()
	testPutGet(db, t)
}

func Test_levelDB_Put(t *testing.T) {
	// Init levelDB
	db, remove := newTestLevelDB()
	defer remove()

	// check insert and get
	err := db.PutString("1", "2")
	assert.Equal(t, err, nil)

	value, err := db.GetString("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, value, "2")
}

func Test_levelDB_Has(t *testing.T) {
	// Init levelDB
	db, remove := newTestLevelDB()
	defer remove()

	// check whether key exists
	db.PutString("1", "2")
	exist, err := db.HasString("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, exist, true)
}

func Test_levelDB_Update(t *testing.T) {
	// Init levelDB
	db, remove := newTestLevelDB()
	defer remove()

	// check update and get
	db.PutString("1", "1")
	value, err := db.GetString("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, value, "1")

	db.PutString("1", "3")
	value, err = db.GetString("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, value, "3")
}

func Test_levelDB_Delete(t *testing.T) {
	// Init levelDB
	db, remove := newTestLevelDB()
	defer remove()

	// insert and then delete key
	db.PutString("1", "1")
	err := db.DeleteSring("1")
	assert.Equal(t, err, nil)

	// check not found
	value, err := db.GetString("3")
	assert.Equal(t, err, leveldb.ErrNotFound)
	assert.Equal(t, value, "")

	// empty set
	exist, err := db.HasString("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, exist, false)

	exist, err = db.HasString("3")
	assert.Equal(t, err, nil)
	assert.Equal(t, exist, false)
}

func Test_levelDB_Newbatch(t *testing.T) {
	// Init levelDB
	db, remove := newTestLevelDB()
	defer remove()

	batch := db.NewBatch()
	if batch == nil {
		t.Fatal("new level batch error")
	}
}
