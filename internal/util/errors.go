package util

import "fmt"

const (
	ErrCodeUnknown = iota
	ErrCodeNotFound
	ErrCodeInvalidArgument
	ErrCodeUnauthorized
)

type Error struct {
	orig    error
	message string
	code    int
}

func WrapErrorf(orig error, code int, format string, a ...interface{}) error {
	return &Error{
		orig:    orig,
		code:    code,
		message: fmt.Sprintf(format, a...),
	}
}

func NewErrorf(code int, format string, a ...interface{}) error {
	return WrapErrorf(nil, code, format, a...)
}

func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.message, e.orig)
	}

	return e.message
}

func (e *Error) Code() int {
	return e.code
}
