package errs

import "errors"

var (
	ErrValueMustBePositive = errors.New("value must be positive")
)
