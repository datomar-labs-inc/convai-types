package ctypes

import (
	"fmt"
	"net/http"

	"upper.io/db.v3"
)

var (
	InsufficientPermissions = &APIError{Code: ErrDatabaseIssue, Message: "Insufficient permissions to perform this action", statusCode: http.StatusUnauthorized}
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
