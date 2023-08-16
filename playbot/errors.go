package playbot

import (
	"errors"
	"fmt"
)

var NoRecordFoundError = errors.New("no matching record found")

type SearchCanceledError struct {
	err error
}

func (e SearchCanceledError) Error() string {
	return fmt.Sprintf("search canceled: %s", e.err)
}

func (e SearchCanceledError) Unwrap() error {
	return e.err
}

func (e SearchCanceledError) Is(target error) bool {
	_, ok := target.(SearchCanceledError)
	return ok
}

var InvalidOffsetError = errors.New("invalid offset")
