package main

import (
    "database/sql"
    "log"

    "github.com/Drack112/simplebank/api"
    db "github.com/Drack112/simplebank/db/sqlc"
    _ "github.com/lib/pq"
)

const (
    dbDriver     = "postgres"
    dbSource     = "postgresql://drack:123@db:5432/simplebank?sslmode=disable"
    serverAdress = "0.0.0.0:8080"
)

func main() {
    log.Printf("Starting server...")

    conn, err := sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect do pg db: ", err)
    }

    store := db.NewStore(conn)
    server := api.NewServer(store)

    err = server.Start(serverAdress)
    if err != nil {
        log.Fatal("cannot connect to server: ", err)
    }

}
