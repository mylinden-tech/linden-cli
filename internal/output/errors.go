package output

import (
	"errors"
	"fmt"
)

// Error is a structured CLI error with code, message, and optional hint.
type Error struct {
	Code       string
	Message    string
	Hint       string
	HTTPStatus int
	Retryable  bool
	Cause      error
}

func (e *Error) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s\nHint: %s", e.Message, e.Hint)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Cause }

// ExitCode returns the OS exit code for this error.
func (e *Error) ExitCode() int {
	switch e.Code {
	case CodeUsage:
		return 1
	case CodeNotFound:
		return 2
	case CodeAuth:
		return 3
	case CodeForbidden:
		return 4
	case CodeRateLimit:
		return 5
	case CodeNetwork:
		return 6
	case CodeAPI:
		return 7
	default:
		return 1
	}
}

// Error code constants.
const (
	CodeUsage     = "usage"
	CodeNotFound  = "not_found"
	CodeAuth      = "auth"
	CodeForbidden = "forbidden"
	CodeRateLimit = "rate_limit"
	CodeNetwork   = "network"
	CodeAPI       = "api"
)

func ErrUsage(msg string) *Error {
	return &Error{Code: CodeUsage, Message: msg}
}

func ErrUsageHint(msg, hint string) *Error {
	return &Error{Code: CodeUsage, Message: msg, Hint: hint}
}

func ErrNotFound(resource, id string) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s %q not found", resource, id),
	}
}

func ErrAuth(msg string) *Error {
	return &Error{
		Code:    CodeAuth,
		Message: msg,
		Hint:    "Run: linden auth login",
	}
}

func ErrForbidden(msg string) *Error {
	return &Error{Code: CodeForbidden, Message: msg, HTTPStatus: 403}
}

func ErrAPI(status int, msg string) *Error {
	return &Error{Code: CodeAPI, Message: msg, HTTPStatus: status}
}

func ErrNetwork(cause error) *Error {
	return &Error{Code: CodeNetwork, Message: "network error: " + cause.Error(), Cause: cause}
}

// AsError wraps any error into *Error. If err is already *Error, returns it as-is.
func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{Code: CodeAPI, Message: err.Error(), Cause: err}
}

// errJQSentinel is the sentinel cause for all jq-related errors.
var errJQSentinel = errors.New("jq unsupported")

func ErrJQValidation(cause error) *Error {
	return &Error{
		Code:    CodeUsage,
		Message: fmt.Sprintf("invalid --jq expression: %s", cause),
		Cause:   errJQSentinel,
	}
}

func IsJQError(err error) bool {
	return errors.Is(err, errJQSentinel)
}
