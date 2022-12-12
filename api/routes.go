package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
    router := gin.Default()

    authRoutes := router.Group("/api").Use(authMiddleware(server.tokenMaker))

    authRoutes.POST("accounts", server.createAccount)
    authRoutes.GET("accounts/:id", server.getAccount)
    authRoutes.GET("accounts", server.listAccounts)
    authRoutes.POST("transfers", server.CreateTransfer)

    router.POST("/api/users", server.createUser)
    router.POST("/api/users/login", server.loginUser)

    server.router = router
}
