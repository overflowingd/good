package context

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	ctxError        = "got_ctx_error"
	ctxErrorMeta    = "got_ctx_error_meta"
	ctxQuery        = "got_ctx_query"
	ctxRequestBody  = "got_ctx_request_body"
	ctxUriParams    = "got_ctx_uri_params"
	ctxResult       = "got_ctx_result"
	ctxResponseBody = "got_ctx_response_body"
)

func AbortWithError(ctx *gin.Context, err error, metas ...any) {
	ctx.Set(ctxError, err)
	ctx.Set(ctxErrorMeta, metas)
	ctx.Abort()
}

func GetError(ctx *gin.Context) (err error, exists bool) {
	e, exists := ctx.Get(ctxError)

	if nil == e || !exists {
		return nil, false
	}

	err, isErr := e.(error)
	if !isErr {
		panic(fmt.Sprintf("Unexpected context item at key: %s. Item: %v", ctxError, e))
	}

	return err, true
}

func GetErrorMeta(ctx *gin.Context) []any {
	meta, exists := ctx.Get(ctxErrorMeta)
	if !exists {
		return nil
	}

	_, converted := meta.([]any)
	if converted {
		return meta.([]any)
	}

	return nil
}

func GetQuery(ctx *gin.Context) (query any, exists bool) {
	return ctx.Get(ctxQuery)
}

func MustGetQuery(ctx *gin.Context) any {
	return ctx.MustGet(ctxQuery)
}

func SetQuery(ctx *gin.Context, query any) {
	ctx.Set(ctxQuery, query)
}

func GetRequestBody(ctx *gin.Context) (body any, exists bool) {
	return ctx.Get(ctxRequestBody)
}

func MustGetRequestBody(ctx *gin.Context) any {
	return ctx.MustGet(ctxRequestBody)
}

func SetRequestBody(ctx *gin.Context, body any) {
	ctx.Set(ctxRequestBody, body)
}

func GetUriParams(ctx *gin.Context) (params any, exists bool) {
	return ctx.Get(ctxUriParams)
}

func MustGetUriParams(ctx *gin.Context) any {
	return ctx.MustGet(ctxUriParams)
}

func SetUriParams(ctx *gin.Context, params any) {
	ctx.Set(ctxUriParams, params)
}

func GetResult(ctx *gin.Context) (r any, exists bool) {
	return ctx.Get(ctxResult)
}

func SetResult(ctx *gin.Context, r any) {
	ctx.Set(ctxResult, r)
}

func SetResultWithAbort(ctx *gin.Context, r any) {
	SetResult(ctx, r)
	ctx.Abort()
}

func GetResponseBody(ctx *gin.Context) (b any, exists bool) {
	return ctx.Get(ctxResponseBody)
}

func SetResponseBody(ctx *gin.Context, b any) {
	ctx.Set(ctxResponseBody, b)
}
