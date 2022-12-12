package api

import (
    "database/sql"
    "errors"
    "fmt"
    "net/http"

    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/Drack112/simplebank/token"
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

    fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
    if !valid {
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    if fromAccount.Owner != authPayload.Username {
        err := errors.New("from account doesn't belong to authenticated user")
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    _, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
    if !valid {
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

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
    account, err := server.db.GetAccount(ctx, accountID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return account, false
        }
        ctx.JSON(http.StatusInternalServerError, err)
        return account, false
    }

    if account.Currency != currency {
        err := fmt.Errorf("account [%d] currency mismatch: [%s] vs [%s]", accountID, account.Currency, currency)
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return account, false
    }

    return account, true
}
