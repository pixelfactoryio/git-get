package internal

import (
	"fmt"
)

// Error represents an error that could be wrapping another error, it includes a code for determining what
// triggered the error.
type Error struct {
	orig error
	msg  string
	code ErrorCode
}

// ErrorCode defines supported error codes.
type ErrorCode uint

// Error codes is also used as exit code
// Must always start at 1
const (
	// ErrorCodeUnknown represents an unknown error
	ErrorCodeUnknown ErrorCode = iota + 1
	// ErrorCodeInvalidArgument represents an invalid argument error
	ErrorCodeInvalidArgument
	// ErrorMissingArgument represents an missing argument error
	ErrorMissingArgument
)

// WrapErrorf returns a wrapped error.
func WrapErrorf(orig error, code ErrorCode, format string, a ...interface{}) error {
	return &Error{
		code: code,
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

// NewErrorf instantiates a new error.
func NewErrorf(code ErrorCode, format string, a ...interface{}) error {
	return WrapErrorf(nil, code, format, a...)
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.orig
}

// Code returns the code representing this error.
func (e *Error) Code() ErrorCode {
	return e.code
}
