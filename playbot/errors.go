package playbot

import (
	"errors"
	"fmt"
)

var (
	// ErrNoRecordFound is the error when no record can be found.
	ErrNoRecordFound = errors.New("no matching record found")
	// ErrInvalidOffset is the error when a offset is invalid.
	ErrInvalidOffset = errors.New("invalid offset")
)

// SearchCancelError is the error when a search has been canceled and cannot be iterated anymore.
type SearchCanceledError struct {
	err error
}

func (e SearchCanceledError) Error() string {
	return fmt.Sprintf("search canceled: %s", e.err)
}

// Unwrap returns the wrapped error.
func (e SearchCanceledError) Unwrap() error {
	return e.err
}

// Is returns true if the error is of the same type of the target.
func (e SearchCanceledError) Is(target error) bool {
	_, ok := target.(SearchCanceledError)
	return ok
}
