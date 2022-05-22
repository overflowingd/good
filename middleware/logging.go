package middleware

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/overflowingd/good/middleware/context"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[SERVICE_ERROR]: ", log.LstdFlags)
}

func LogErrors(ctx *gin.Context) {
	ctx.Next()

	err, exists := context.GetError(ctx)
	if !exists {
		return
	}

	if err != nil {
		logger.Println(err)

		metas := context.GetErrorMeta(ctx)
		if len(metas) == 0 {
			return
		}

		logger.Println(metas...)
	}
}
