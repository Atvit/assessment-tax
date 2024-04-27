package utils

import (
	"golang.org/x/exp/slices"
)

func Gte(value1, value2 float64) bool {
	return value1 >= value2
}

func Lte(value1, value2 float64) bool {
	return value1 <= value2
}

func Gt(value1, value2 float64) bool {
	return value1 > value2
}

func Oneof(value string, items ...string) bool {
	return slices.Contains(items, value)
}
