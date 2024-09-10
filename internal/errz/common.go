package errz

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/withstack"
)



type PrettyError struct {
	HttpStatusCode int
	Code string
	Message string
	cause error
}

func (e PrettyError) Error() string {
	return e.cause.Error()
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

	errors.New("aa")

	return PrettyError{
		HttpStatusCode: httpStatusCode,
		Code: code,
		Message: message,
		cause: withstack.WithStackDepth(cause, 1),
	}
}
