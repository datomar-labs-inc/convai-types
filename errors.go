package ctypes

import (
	"fmt"
	"net/http"

	"upper.io/db.v3"
)

var (
	InsufficientPermissions = &APIError{Code: ErrInsufficientPermissions, Message: "Insufficient permissions to perform this action", statusCode: http.StatusUnauthorized}
	RedisFailure            = &APIError{Code: ErrRedisFailure, Message: "Something went wrong", statusCode: http.StatusInternalServerError}
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
