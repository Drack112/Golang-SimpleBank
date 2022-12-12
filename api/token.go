package api

import (
    "database/sql"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
    AccessToken          string    `json:"access_token"`
    AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
    var req renewAccessTokenRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    session, err := server.db.GetSession(ctx, refreshPayload.ID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    if session.IsBlocked {
        err := fmt.Errorf("blocked session")
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    if session.Username != refreshPayload.Username {
        err := fmt.Errorf("incorrect session user")
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    if session.RefreshToken != req.RefreshToken {
        err := fmt.Errorf("mismatched session token")
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    if time.Now().After(session.ExpiresAt) {
        err := fmt.Errorf("expired session")
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    accessToken, accessPayload, err := server.tokenMaker.CreateToken(
        refreshPayload.Username,
        server.config.AccessTokenDuration,
    )
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    rsp := renewAccessTokenResponse{
        AccessToken:          accessToken,
        AccessTokenExpiresAt: accessPayload.ExpiredAt,
    }
    ctx.JSON(http.StatusOK, rsp)
}
