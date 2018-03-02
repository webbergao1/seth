package database

import (
	"fmt"
	"seth/log"
)

// Database interface of store
type Database interface {
	Name() string
	Open(path string, cache int, handles int) error
	Close()
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	GetString(key string) (string, error)
	PutString(key string, value string) error
	Has(key []byte) (ret bool, err error)
	HasString(key string) (ret bool, err error)
	Delete(key []byte) error
	DeleteSring(key string) error
	NewBatch() Batch
}

// Batch iterface of batch for database
type Batch interface {
	Put(key []byte, value []byte)
	Delete(key []byte)
	Commit() error
	Rollback()
	Close()
}

var dbs = map[string]Database{}

// Register register database implement
func Register(db Database) {
	name := db.Name()
	if _, ok := dbs[name]; ok {
		log.Error("database %s is registered", name)
		return
	}

	dbs[name] = db
}

// GetDatabase get database by name
func GetDatabase(name string) (Database, error) {
	db, ok := dbs[name]
	if !ok {
		log.Error("database %s is not registered", name)
		return nil, fmt.Errorf("database %s is not registered", name)
	}

	return db, nil
}

// ListDatabases list register databases
func ListDatabases() []string {
	names := []string{}
	for name := range dbs {
		names = append(names, name)
	}

	return names
}
