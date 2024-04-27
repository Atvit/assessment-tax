package errs

import "errors"

var (
	ErrValueMustBePositive           = errors.New("value must be positive")
	ErrWhtMustLowerThanOrEqualIncome = errors.New("with holding tax must be lower than or equal to income")
)
