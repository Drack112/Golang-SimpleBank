package api

import (
    "database/sql"
    "fmt"
    "net/http"

    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/gin-gonic/gin"
)

type createTransferRequest struct {
    FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
    ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
    Amount        int64  `json:"amount" binding:"required,min=1"`
    Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
    var req createTransferRequest
    if err := ctx.BindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
        return
    }

    if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
        return
    }

    arg := db.TransferTxParams{
        FromAccountID: req.FromAccountID,
        ToAccountID:   req.ToAccountID,
        Amount:        req.Amount,
    }

    result, err := server.db.TransferTx(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
    account, err := server.db.GetAccount(ctx, accountID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return false
        }
        ctx.JSON(http.StatusInternalServerError, err)
        return false
    }

    if account.Currency != currency {
        err := fmt.Errorf("account [%d] currency mismatch: [%s] vs [%s]", accountID, account.Currency, currency)
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return false
    }

    return true
}
