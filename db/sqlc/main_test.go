package db

import (
    "database/sql"
    "log"
    "os"
    "testing"

    _ "github.com/lib/pq"
)

const (
    dbDriver = "postgres"
    dbSource = "postgres://db_test_user:db_test_password@localhost:5432/db_test_database?sslmode=disable"
)

var (
    testQueries *Queries
    testDB      *sql.DB
)

func TestMain(m *testing.M) {
    var err error
    testDB, err = sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    testQueries = New(testDB)

    os.Exit(m.Run())
}
