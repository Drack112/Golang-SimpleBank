package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
    router := gin.Default()

    router.POST("/api/accounts", server.createAccount)
    router.GET("/api/accounts/:id", server.getAccount)
    router.GET("/api/accounts", server.listAccounts)

    router.POST("/api/transfers", server.CreateTransfer)

    router.POST("/api/users", server.createUser)
    router.POST("/api/users/login", server.loginUser)

    server.router = router
}
