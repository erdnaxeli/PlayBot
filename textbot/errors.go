package textbot

import "errors"

var (
	// ErrAuthenticationRequired is the error when the user is not authenticated and must be authenticated.
	ErrAuthenticationRequired = errors.New("authentication required")
	// ErrInvalidUsage is the error when a command is called with incorrect parameters.
	ErrInvalidUsage = errors.New("invalid command usage")
	// ErrNotImplemented is the error when a command is not implemented.
	ErrNotImplemented = errors.New("not implemented")
	// ErrOffsetToBig is the error when the offset given to `!get` or other commands is too big.
	// The offset cannot exceed -1*.
	ErrOffsetToBig = errors.New("offset too big")
)
