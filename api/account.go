package api

import (
    "database/sql"
    "net/http"

    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/gin-gonic/gin"
)

type createAccountRequest struct {
    Owner    string `json:"owner" binding:"required"`
    Currency string `json:"currency" binding:"required,oneof=USD EUR""`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
    var req createAccountRequest
    if err := ctx.BindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    arg := db.CreateAccountParams{
        Owner:    req.Owner,
        Currency: req.Currency,
    }

    account, err := server.db.CreateAccount(ctx, arg)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccountByID(ctx *gin.Context) {
    var req getAccountRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    account, err := server.db.GetAccount(ctx, req.ID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
    PageID   int32 `form:"page_id" binding:"required,min=1"`
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccounts(ctx *gin.Context) {
    var req listAccountRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    arg := db.ListAccountsParams{
        Limit:  req.PageSize,
        Offset: (req.PageID - 1) * req.PageSize,
    }

    accounts, err := server.db.ListAccounts(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, accounts)

}
