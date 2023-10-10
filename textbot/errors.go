package textbot

import "errors"

var (
	InvalidUsageError      = errors.New("invalid command usage")
	NotImplementedError    = errors.New("not implemented")
	AuthenticationRequired = errors.New("authentication required")
)
