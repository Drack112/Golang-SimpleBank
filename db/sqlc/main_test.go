package db

import (
    "database/sql"
    "log"
    "os"
    "testing"

    "github.com/Drack112/simplebank/util"
    _ "github.com/lib/pq"
)

var (
    testQueries *Queries
    testDB      *sql.DB
)

func TestMain(m *testing.M) {

    config, err := util.LoadConfig("../..")
    if err != nil {
        log.Panic("Cannot get app.env: ", err)
    }

    testDB, err = sql.Open(config.DBDriver, "postgresql://db_test_user:db_test_password@localhost:5432/db_test_database?sslmode=disable")
    if err != nil {
        log.Fatal("Cannot connect to PostgreSQL Database: ", err)
    }

    testQueries = New(testDB)

    os.Exit(m.Run())
}
