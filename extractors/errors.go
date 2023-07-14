package extractors

import "fmt"

type UnknownRecordError struct {
	id string
}

func (err *UnknownRecordError) Error() string {
	return fmt.Sprintf("unknown record with id '%s'", err.id)
}

type UnknownRecordSourceError struct{}

func (err *UnknownRecordSourceError) Error() string {
	return "unknown record source error"
}
