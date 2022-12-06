package api

import (
    "database/sql"
    "net/http"
    "time"

    db "github.com/Drack112/simplebank/db/sqlc"
    "github.com/Drack112/simplebank/util"
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
)

type createUserRequest struct {
    Username string `json:"username" binding:"required,alphanum"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
    Username          string    `json:"username"`
    FullName          string    `json:"full_name"`
    Email             string    `json:"email"`
    PasswordChangedAt time.Time `json:"password_changed_at"`
    CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) createUserResponse {
    return createUserResponse{
        Username:          user.Username,
        FullName:          user.FullName,
        Email:             user.Email,
        PasswordChangedAt: user.PasswordChangedAt,
        CreatedAt:         user.CreatedAt,
    }
}

func (server *Server) createUser(ctx *gin.Context) {
    var req createUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    hashedPassword, err := util.HashPassword(req.Password)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    arg := db.CreateUserParams{
        Username:       req.Username,
        FullName:       req.FullName,
        HashedPassword: hashedPassword,
        Email:          req.Email,
    }

    user, err := server.db.CreateUser(ctx, arg)
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code.Name() {
            case "unique_violation":
                ctx.JSON(http.StatusForbidden, ErrorResponse(err))
                return
            }
        }
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    rsp := newUserResponse(user)

    ctx.JSON(http.StatusOK, rsp)
    return
}

type loginRequest struct {
    Username string `json:"username" binding:"required,alphanum"`
    Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
    AccessToken string             `json:"access_token"`
    User        createUserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
    var req loginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    user, err := server.db.GetUser(ctx, req.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
    }

    err = util.CheckPassword(req.Password, user.HashedPassword)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
        return
    }

    accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    rsp := loginUserResponse{
        AccessToken: accessToken,
        User:        newUserResponse(user),
    }

    ctx.JSON(http.StatusOK, rsp)
    return
}
