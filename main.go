package main

import (
	"database/sql"
	"log"

	"github.com/Drack112/simplebank/api"
	db "github.com/Drack112/simplebank/db/sqlc"
	"github.com/Drack112/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {

    config, err := util.LoadConfig("./")
    if err != nil {
        log.Fatal("cannot load config: ", err)
    }

    log.Printf("Starting server...")

    conn, err := sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal("cannot connect do pg db: ", err)
    }

    store := db.NewStore(conn)
    server := api.NewServer(store)

    err = server.Start(config.ServerAddress)
    if err != nil {
        log.Fatal("cannot connect to server: ", err)
    }

}
