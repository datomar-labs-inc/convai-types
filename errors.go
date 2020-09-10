package ctypes

import (
	"fmt"
	"net/http"

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
