package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQuery *Queries
var testDb *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:123456@localhost:5432/suoga?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cann't connect to db: ", err)
	}

	testQuery = New(testDb)
	os.Exit(m.Run())
}
