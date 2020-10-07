package ctypes

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"upper.io/db.v3"
)

var (
	RateLimitExceeded           = &APIError{Code: ErrRateLimitExceeded, Message: "Rate limit exceeded", statusCode: http.StatusTooManyRequests}
	NotFoundError               = &APIError{Code: ErrResourceNotFound, Message: "Resource not found", statusCode: http.StatusNotFound}
	InvalidTokenError           = &APIError{Code: ErrInvalidToken, Message: "Invalid authentication token", statusCode: http.StatusUnauthorized}
	NotAuthenticatedError       = &APIError{Code: ErrNotAuthenticated, Message: "Cannot perform action without being authenticated", statusCode: http.StatusUnauthorized}
	InvalidHeaderFormat         = &APIError{Code: ErrInvalidHeaderFormat, Message: "A header was supplied with an invalid format", statusCode: http.StatusBadRequest}
	InvalidParameterFormat      = &APIError{Code: ErrInvalidParameter, Message: "A parameter was supplied in an invalid format", statusCode: http.StatusBadRequest}
	MissingEnvironmentIDHeader  = &APIError{Code: ErrMissingEnvHeader, Message: "The X-Environment-ID header must be present", statusCode: http.StatusForbidden}
	MissingBotIDHeader          = &APIError{Code: ErrMissingBotHeader, Message: "The X-Bot-ID header must be present", statusCode: http.StatusForbidden}
	MissingOrganizationIDHeader = &APIError{Code: ErrMissingOrgHeader, Message: "The X-Organization-ID header must be present", statusCode: http.StatusForbidden}
	InsufficientPermissions     = &APIError{Code: ErrInsufficientPermissions, Message: "Insufficient permissions to perform this action", statusCode: http.StatusForbidden}
	RedisFailure                = &APIError{Code: ErrRedisFailure, Message: "Something went wrong", statusCode: http.StatusInternalServerError}
	GrooveFailure               = &APIError{Code: ErrGrooveFailure, Message: "Something went wrong", statusCode: http.StatusInternalServerError}
	MongoFailure                = &APIError{Code: ErrMongoFailure, Message: "Something went wrong", statusCode: http.StatusInternalServerError}
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
	var errList []FieldError

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return &APIError{
			statusCode: http.StatusBadRequest,
			Code:       ErrInputValidation,
			Message:    err.Error(),
		}
	}

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

func GenericError(err error) *APIError {
	msg := "Something has gone wrong (generic)"

	if err != nil {
		msg = err.Error()
	}

	return &APIError{
		Code:    ErrGenericError,
		Message: msg,
	}
}

func InternalError(err error) *APIError {
	msg := "Convai encountered an internal error"
	if err != nil {
		msg = err.Error()
	}
	return &APIError{
		Code:    ErrGenericError,
		Message: msg,
		statusCode: http.StatusInternalServerError,
	}
}

