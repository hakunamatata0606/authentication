package db

import (
	"authentication/config"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testConn *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	c := config.GetConfig()
	conn, err := sql.Open(c.Db.Driver, c.Db.Addr)
	if err != nil {
		log.Fatal("failed to connect to datbase")
	}
	testConn, err = sql.Open(c.Db.Driver, c.Db.Addr)
	if err != nil {
		log.Fatal("failed to connect to datbase 1")
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
