package errz

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/withstack"
)

type PrettyError struct {
	HttpStatusCode int
	Code           string
	Message        string
	Details        map[string]string // Optional additional details
	cause          error

	// TODO: add error details
}

func (e PrettyError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}

	return e.Message
}

func (e PrettyError) Unwrap() error {
	return e.cause
}

func (e PrettyError) Cause() error {
	return e.cause
}

func (e PrettyError) Format(s fmt.State, verb rune) {
	errors.FormatError(e, s, verb)
}

func NewPrettyError(httpStatusCode int, code, message string, cause error) PrettyError {

	return PrettyError{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
		cause:          withstack.WithStackDepth(cause, 1),
	}
}

func NewPrettyErrorDetail(httpStatusCode int, code, message string, cause error, details map[string]string) PrettyError {

	return PrettyError{
		HttpStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
		Details:        details,
		cause:          withstack.WithStackDepth(cause, 1),
	}
}
