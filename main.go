package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/Drack112/simplebank/api"
	db "github.com/Drack112/simplebank/db/sqlc"
	"github.com/Drack112/simplebank/gapi"
	"github.com/Drack112/simplebank/pb"
	"github.com/Drack112/simplebank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	log.Printf("Starting server...")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect do PostgreSQL db: ", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	gRPCServer := grpc.NewServer()

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	pb.RegisterSimpleBankServer(gRPCServer, server)
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener")
	}

	log.Printf("start gRPC server at: %s", listener.Addr().String())
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create the server: ", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("Cannot connect to the Gin Server: ", err)
	}
}
