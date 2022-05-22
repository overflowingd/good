package middleware

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/overflowingd/good/middleware/context"
)

func BindRequestBody(bind any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// todo: Handle all available bind initializations.
		//       Currently supported `&` initialization like `&struct{}`.
		//       Add support for `new(struct)`, `struct{}`
		body := reflect.New(reflect.ValueOf(bind).Elem().Type()).Interface()

		// todo: Add proper error handling
		if err := ctx.ShouldBindBodyWith(body, binding.JSON); err != nil {
			context.AbortWithError(ctx, ErrInvalid, err)
			return
		}

		context.SetRequestBody(ctx, body)
		ctx.Next()
	}
}

func BindQuery(bind any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// todo: Handle all available bind initializations.
		//       Currently supported `&` initialization like `&struct{}`.
		//       Add support for `new(struct)`, `struct{}`
		query := reflect.New(reflect.ValueOf(bind).Elem().Type()).Interface()

		// todo: Add proper error handling
		if err := ctx.ShouldBindQuery(query); err != nil {
			context.AbortWithError(ctx, ErrInvalid, err)
			return
		}

		context.SetQuery(ctx, query)
		ctx.Next()
	}
}

func BindRequestUriParams(bind any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// todo: Handle all available bind initializations.
		//       Currently supported `&` initialization like `&struct{}`.
		//       Add support for `new(struct)`, `struct{}`
		params := reflect.New(reflect.ValueOf(bind).Elem().Type()).Interface()

		// todo: Add proper error handling
		if err := ctx.ShouldBindUri(params); err != nil {
			context.AbortWithError(ctx, ErrInvalid, err)
			return
		}

		context.SetUriParams(ctx, params)
		ctx.Next()
	}
}
