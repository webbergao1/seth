package database_test

import (
	"fmt"
	"seth/database"
	_ "seth/database/leveldb"
	"testing"
)

func Test_database(t *testing.T) {
	dbs := database.ListDatabases()
	for _, name := range dbs {
		fmt.Println(name)
	}
	db, err := database.GetDatabase("leveldb")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(db.Name())
	}
	db, err = database.GetDatabase("mysql")

}
