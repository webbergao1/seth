package trie

import (
	"fmt"
	"io/ioutil"
	"os"
	"seth/common"
	"seth/database"
	"seth/database/leveldb"
	"testing"
)

func newTestTrieDB() (database.Database, func()) {
	dir, err := ioutil.TempDir("", "trietest")
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

func Test_trie_Update(t *testing.T) {
	db, remove := newTestTrieDB()
	defer remove()
	trie, err := NewTrie(common.Hash{}, []byte("trietest"), db)
	if err != nil {
		panic(err)
	}
	trie.Update([]byte("12345678"), []byte("test"))
	trie.Update([]byte("12345678"), []byte("testnew"))
	trie.Update([]byte("12345557"), []byte("test1"))
	value := trie.Get([]byte("12345678"))
	fmt.Println(string(value))
	value = trie.Get([]byte("12345557"))
	fmt.Println(string(value))
}
