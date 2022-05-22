package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gobeam/stringy"
	dberrors "github.com/overflowingd/good/db/errors"
	"github.com/overflowingd/good/middleware/context"
	"gorm.io/gorm"
)

const (
	CodeUnknown        = "unknown"
	CodeNotImplemented = "not_implemented"
	CodeDataInvalid    = "data_invalid"
	CodeNotFound       = "not_found"
	CodeAlreadyExists  = "already_exists"
)

var (
	ErrUnknown = NewRestfulError(
		CodeUnknown,
		"The unknown internal error occured. Please contact our support",
	)
	ErrInvalid = NewRestfulError(
		CodeDataInvalid,
		"Data you send is invalid",
	)
	ErrNotImplemented = NewRestfulError(
		CodeNotImplemented,
		"Method on the requested resource is not implemented",
	)
	ErrNotFound = NewRestfulError(
		CodeNotFound,
		"Resource you requested does not exist",
	)
	ErrAlreadyExists = NewRestfulError(
		CodeAlreadyExists,
		"Resource you want to create already exists",
	)
)

var debug bool

func init() {
	_, debug = os.LookupEnv("DEBUG")
}

func NewRestfulError(code string, message string) *RestfulError {
	return &RestfulError{
		Code:    code,
		Message: message,
	}
}

type RestfulError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *RestfulError) Error() string {
	return fmt.Sprintf("%v: %s", time.Now().UTC(), e.Message)
}

type ErrorResponseBody struct {
	Error *RestfulError `json:"error"`
	Meta  any           `json:"_meta"`
	Debug any           `json:"_debug"`
}

type SuccessResponseBody struct {
	Payload any `json:"payload"`
	Meta    any `json:"_meta"`
}

func BuildRestfulResponse(ctx *gin.Context) {
	ctx.Next()

	err, exists := context.GetError(ctx)
	if exists {
		restfulFail(ctx, err)
		return
	}

	restfulSuccess(ctx)
}

func resolveErr(err error) *RestfulError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	myerr, isMysql := err.(*mysql.MySQLError)
	if isMysql {
		if dberrors.Is(dberrors.EDuplicateEntry, int(myerr.Number)) {
			return ErrAlreadyExists
		}
	}

	return ErrUnknown
}

func resolveErrMetas(metas []any) []any {
	var resolved []any

	for _, meta := range metas {
		switch meta := meta.(type) {
		case *json.UnmarshalTypeError, validator.ValidationErrors:
			resolved = append(resolved, representValidationErrs(meta))
			continue
		}

		resolved = append(resolved, meta)
	}

	return resolved
}

func representValidationErrs(errs any) []map[string]any {
	var results []map[string]any

	validationErrs, isValidation := errs.(validator.ValidationErrors)
	log.Println("Validation", isValidation)
	if isValidation {
		for _, err := range validationErrs {
			field := stringy.New(err.Field()).SnakeCase().ToLower()

			results = append(
				results,
				map[string]any{
					"rule":   err.Tag(),
					"field":  field,
					"actual": err.Value(),
					"message": fmt.Sprintf(
						"The value of field %q you passed does not satisfy validation rule %q",
						field,
						err.Tag(),
					),
				},
			)
		}

		return results
	}

	unmarshalTypeErr, isUnmarshalType := errs.(*json.UnmarshalTypeError)
	if isUnmarshalType {
		field := stringy.New(unmarshalTypeErr.Field).SnakeCase().ToLower()
		results = append(
			results,
			map[string]any{
				"rule":   "expect_type",
				"field":  field,
				"actual": unmarshalTypeErr.Value,
				"message": fmt.Sprintf(
					"The field %q you passed can not be of type %q",
					field,
					unmarshalTypeErr.Value,
				),
			},
		)
	}

	return results
}

func restfulFail(ctx *gin.Context, err error) {
	metas := context.GetErrorMeta(ctx)

	body := new(ErrorResponseBody)

	_, restful := err.(*RestfulError)
	if restful {
		body.Error = err.(*RestfulError)
	} else {
		body.Error = resolveErr(err)

		if debug {
			body.Debug = map[string]any{
				"previous": err,
			}
		}
	}

	body.Meta = resolveErrMetas(metas)

	switch body.Error.Code {
	case CodeUnknown:
		ctx.Status(http.StatusInternalServerError)
	case CodeDataInvalid:
		ctx.Status(http.StatusBadRequest)
	case CodeNotImplemented:
		ctx.Status(http.StatusNotImplemented)
	case CodeNotFound:
		ctx.Status(http.StatusNotFound)
	case CodeAlreadyExists:
		ctx.Status(http.StatusUnprocessableEntity)
	}

	context.SetResponseBody(ctx, body)
}

func restfulSuccess(ctx *gin.Context) {
	var payload any

	result, exists := context.GetResult(ctx)
	if !exists {
		ctx.Status(http.StatusNoContent)
		return
	}

	switch result := result.(type) {
	case *Processing:
		ctx.Status(http.StatusAccepted)
		payload = result.Result
	case *Created:
		ctx.Status(http.StatusCreated)
		payload = result.Result
	default:
		ctx.Status(http.StatusOK)
		payload = result
	}

	body := new(SuccessResponseBody)
	body.Payload = payload

	context.SetResponseBody(ctx, body)
}
