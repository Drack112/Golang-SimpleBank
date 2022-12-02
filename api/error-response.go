package api

import (
    "github.com/gin-gonic/gin"
)

func ErrorResponse(err error) gin.H {
    return gin.H{
        "message": err.Error(),
    }
}
