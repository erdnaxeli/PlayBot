package ldjson

import "errors"

// Errors returned by the Extract method
var (
	ErrHTTPNotOk = errors.New("received an HTTP Error")
	ErrNoLDJSON  = errors.New("no LD+JSON data found")
)
