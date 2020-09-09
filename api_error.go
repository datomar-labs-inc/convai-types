package ctypes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	statusCode int
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (e *APIError) HTTPStatusCode() int {
	return e.statusCode
}

func (e *APIError) JSONResponse() *APIError {
	return e
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error(%d): %s", e.Code, e.Message)
}

func (e *APIError) WithHTTPCode(status int) *APIError {
	e.statusCode = status
	return e
}

func (e *APIError) WithMessage(message string) *APIError {
	e.Message = message
	return e
}

func (e *APIError) WithCode(code int) *APIError {
	e.Code = code
	return e
}

func (e *APIError) WithData(data interface{}) *APIError {
	e.Data = data
	return e
}

func (e *APIError) RespondGin(c *gin.Context) {
	c.JSON(e.statusCode, e)
}

func (e *APIError) AbortGin(c *gin.Context) {
	c.AbortWithStatusJSON(e.statusCode, e)
}

type IAPIError interface {
	HTTPStatusCode() int
	JSONResponse() *APIError
	Error() string
}
