package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testConn *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("mysql", "bao:123@tcp(172.17.0.2:3306)/test")
	if err != nil {
		log.Fatal("failed to connect to datbase")
	}
	testConn, err = sql.Open("mysql", "bao:123@tcp(172.17.0.2:3306)/test")
	if err != nil {
		log.Fatal("failed to connect to datbase 1")
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
