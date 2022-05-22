package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/overflowingd/good/middleware/context"
)

func NotImplemented(ctx *gin.Context) {
	context.AbortWithError(ctx, ErrNotImplemented)
}
