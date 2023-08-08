package playbot

import "fmt"

type NoRecordFoundError struct{}

func (NoRecordFoundError) Error() string {
	return "no matching record found"
}

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
