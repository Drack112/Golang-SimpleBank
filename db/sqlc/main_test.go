package db

import (
    "database/sql"
    "log"
    "os"
    "testing"

    _ "github.com/lib/pq"
)

var (
    testQueries *Queries
    testDB      *sql.DB
)

func TestMain(m *testing.M) {

    var err error
    testDB, err = sql.Open("postgres", "postgresql://db_test_user:db_test_password@localhost:5432/db_test_database?sslmode=disable")
    if err != nil {
        log.Fatal("Cannot connect to PostgreSQL Database: ", err)
    }

    testQueries = New(testDB)

    os.Exit(m.Run())
}
