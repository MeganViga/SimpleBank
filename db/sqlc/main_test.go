package db

import (
	"database/sql"
	"log"
	"testing"
	"os"
	_ "github.com/lib/pq"
)
var testQueries *Queries
var testDB *sql.DB
const (
	dbDriver = "postgres"
	dbSource ="postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)
func TestMain(m *testing.M){
	var err error
	testDB, err = sql.Open(dbDriver,dbSource)
	if err != nil{
		log.Fatal(err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}