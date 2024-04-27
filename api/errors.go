package api

import "net/http"

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

func ErrUnAuthorized() *Error {
	return NewError(http.StatusUnauthorized, "unauthorized")
}

func ErrInvalidId() *Error {
	return NewError(http.StatusBadRequest, "invalid id given")
}
