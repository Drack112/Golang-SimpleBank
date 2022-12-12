package api

import (
    "database/sql"
    "errors"
    "net/http"

    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/Drack112/simplebank/token"
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
)

type createAccountRequest struct {
    Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
    var req createAccountRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    arg := db.CreateAccountParams{
        Owner:    authPayload.Username,
        Currency: req.Currency,
        Balance:  0,
    }

    account, err := server.db.CreateAccount(ctx, arg)
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code.Name() {
            case "foreign_key_violation", "unique_violation":
                ctx.JSON(http.StatusForbidden, ErrorResponse(err))
                return
            }
        }
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
    var req getAccountRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    account, err := server.db.GetAccount(ctx, req.ID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return
        }

        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

    if account.Owner != authPayload.Username {
        err := errors.New("account doesn't belong to the authenticated user")
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
    PageID   int32 `form:"page_id" binding:"required,min=1"`
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
    var req listAccountRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    arg := db.ListAccountsParams{
        Owner:  authPayload.Username,
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
