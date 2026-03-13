package errors

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

var (
	ErrUnauthorized = New("UNAUTHORIZED", "Unauthorized")
	ErrBadRequest   = New("BAD_REQUEST", "Bad request")
	ErrNotFound     = New("NOT_FOUND", "Not found")
)
