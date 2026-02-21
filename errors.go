package hyperacc

import (
	"errors"
	"fmt"
)

type AccessError struct {
	Reason string
	Cause  error
}

func (e *AccessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("access error: %s: %v", e.Reason, e.Cause)
	}
	return "access error: " + e.Reason
}

func (e *AccessError) Unwrap() error { return e.Cause }

func NewAccessError(reason string) *AccessError {
	return &AccessError{Reason: reason}
}

func WrapAccessError(reason string, cause error) *AccessError {
	return &AccessError{
		Reason: reason,
		Cause:  cause,
	}
}

func AsAccessError(err error) (*AccessError, bool) {
	var e *AccessError
	return e, errors.As(err, &e)
}
