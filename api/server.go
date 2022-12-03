package api

import (
    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/gin-gonic/gin"
)

type Server struct {
    db     db.Store
    router *gin.Engine
}

func NewServer(db db.Store) *Server {
    server := &Server{
        db: db,
    }

    router := gin.Default()
    router.POST("/api/accounts", server.CreateAccount)
    router.GET("/api/accounts/:id", server.GetAccountByID)
    router.GET("/api/accounts", server.ListAccounts)

    server.router = router

    return server
}

func (server *Server) Start(address string) error {
    return server.router.Run(address)
}
