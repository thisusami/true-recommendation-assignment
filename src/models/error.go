package models

import "github.com/gofiber/fiber/v2"

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code"`
}

func NewError(code int, message string, err string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Error:   err,
	}
}
func (e *Error) Set(message string) *Error {
	e.Message = message
	return e
}
func (e *Error) Clone() *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Error:   e.Error,
	}
}
func (e *Error) ToFiberMap() (int, *fiber.Map) {
	return e.Code, &fiber.Map{
		"message": e.Message,
		"error":   e.Error,
	}
}

var (
	InternalServerError     = NewError(500, "An unexpected error occurred", "internal_server_error")
	BadRequestError         = NewError(400, "Invalid limit parameter", "bad_request")
	NotFoundError           = NewError(404, "Not Found", "not_found")
	ServiceUnavailableError = NewError(503, "Recommendation model is temporarily unavailable", "service_unavailable")
)
