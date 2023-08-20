package textbot

import "errors"

var InvalidUsageError = errors.New("invalid command usage")
var NotImplementedError = errors.New("not implemented")
var AuthenticationRequired = errors.New("authentication required")
