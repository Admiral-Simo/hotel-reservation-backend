package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *Error) Error() string {
	return e.Err
}

func NewError(code int, err string) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var (
		apiError *Error
		ok       bool
	)
	if apiError, ok = err.(*Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError = NewError(http.StatusInternalServerError, "internal server error")
	return c.Status(apiError.Code).JSON(apiError)
}

func ErrUnAuthorized() *Error {
	return NewError(http.StatusUnauthorized, "unauthorized request")
}

func ErrInvalidId() *Error {
	return NewError(http.StatusBadRequest, "invalid id given")
}

func ErrBadRequest() *Error {
	return NewError(http.StatusBadRequest, "invalid JSON request")
}

func ErrNotFound(resource string) *Error {
	return NewError(http.StatusNotFound, resource+" not found")
}

func ErrUnavailable(resource string) *Error {
	return NewError(http.StatusNotFound, resource+" unavailable at the moment")
}
