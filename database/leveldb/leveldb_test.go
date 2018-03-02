package leveldb

import (
	"io/ioutil"
	"os"
	"seth/database"
	"testing"

	"github.com/magiconair/properties/assert"
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

func Test_Put(t *testing.T) {
	// Init levelDB
	dir := prepareDbFolder("", "leveldbtest")
	defer os.RemoveAll(dir)
	db := newDbInstance(dir)
	defer db.Close()

	// check insert and get
	err := db.PutString("1", "2")
	assert.Equal(t, err, nil)

	value, err := db.GetString("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, value, "2")
}
