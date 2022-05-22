package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/overflowingd/good/middleware/context"
)

func WriteResponseBody(ctx *gin.Context) {
	ctx.Next()

	body, exists := context.GetResponseBody(ctx)
	if !exists {
		return
	}

	renderer := render.JSON{
		Data: body,
	}

	if err := renderer.Render(ctx.Writer); err != nil {
		panic(err)
	}
}
