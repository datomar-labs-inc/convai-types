package ctypes

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"upper.io/db.v3"
)

var (
	InvalidHeaderFormat         = &APIError{Code: ErrInvalidHeaderFormat, Message: "A header was supplied with an invalid format", statusCode: http.StatusBadRequest}
	MissingEnvironmentIDHeader  = &APIError{Code: ErrMissingEnvHeader, Message: "The X-Environment-ID header must be present", statusCode: http.StatusUnauthorized}
	MissingBotIDHeader          = &APIError{Code: ErrMissingBotHeader, Message: "The X-Bot-ID header must be present", statusCode: http.StatusUnauthorized}
	MissingOrganizationIDHeader = &APIError{Code: ErrMissingOrgHeader, Message: "The X-Organization-ID header must be present", statusCode: http.StatusUnauthorized}
	InsufficientPermissions     = &APIError{Code: ErrInsufficientPermissions, Message: "Insufficient permissions to perform this action", statusCode: http.StatusUnauthorized}
	RedisFailure                = &APIError{Code: ErrRedisFailure, Message: "Something went wrong", statusCode: http.StatusInternalServerError}
)

type FieldError struct {
	Tag       string      `json:"tag"`
	ActualTag string      `json:"actual_tag"`
	Namespace string      `json:"namespace"`
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
}

// InputValidationError should be called with errors from the `validator` package only
func InputValidationError(err error) *APIError {
	errs := err.(validator.ValidationErrors)

	var errList []FieldError

	for _, e := range errs {
		errList = append(errList, FieldError{
			Tag:       e.Tag(),
			ActualTag: e.ActualTag(),
			Namespace: e.Namespace(),
			Field:     e.Field(),
			Value:     e.Value(),
		})
	}

	return &APIError{
		Code:    ErrInputValidation,
		Message: err.Error(),
		Data:    errList,
	}
}

func DatabaseError(err error) *APIError {
	apiErr := &APIError{}

	switch err {
	case db.ErrNoMoreRows:
		apiErr = apiErr.
			WithCode(ErrResourceNotFound).
			WithHTTPCode(http.StatusNotFound).
			WithMessage("resource not found")
	default:
		apiErr = apiErr.
			WithCode(ErrDatabaseIssue).
			WithHTTPCode(http.StatusInternalServerError).
			WithMessage(fmt.Sprintf("db err: %s", err.Error()))
	}

	return apiErr
}
