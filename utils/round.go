package utils

import "math"

func Round(num float64, precision int) float64 {
	power := math.Pow10(precision)
	rounded := math.Round(num*power) / power
	return rounded
}
