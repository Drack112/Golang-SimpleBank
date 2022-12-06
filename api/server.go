package api

import (
    "fmt"

    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/Drack112/simplebank/token"
    "github.com/Drack112/simplebank/util"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
)

type Server struct {
    config     util.Config
    db         db.Store
    tokenMaker token.Maker
    router     *gin.Engine
}

func NewServer(config util.Config, db db.Store) (*Server, error) {

    tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
    if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }

    server := &Server{
        config:     config,
        db:         db,
        tokenMaker: tokenMaker,
    }

    router := gin.Default()

    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("currency", validCurrency)
    }

    router.POST("/api/accounts", server.createAccount)
    router.GET("/api/accounts/:id", server.getAccount)
    router.GET("/api/accounts", server.listAccounts)

    router.POST("/api/transfers", server.CreateTransfer)

    router.POST("/api/users", server.createUser)
    router.POST("/api/users/login", server.loginUser)

    server.router = router

    return server, nil
}

func (server *Server) Start(address string) error {
    return server.router.Run(address)
}
