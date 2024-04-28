package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGte(t *testing.T) {
	type testcase struct {
		Name     string
		Value1   float64
		Value2   float64
		Expected bool
	}

	tcs := []testcase{
		{"value 1 greater than value 2", 5, 4, true},
		{"value 1 equal value 2", 5, 5, true},
		{"value 1 less than value 2", 4, 5, false},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			result := Gte(tc.Value1, tc.Value2)

			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestLte(t *testing.T) {
	type testcase struct {
		Name     string
		Value1   float64
		Value2   float64
		Expected bool
	}

	tcs := []testcase{
		{"value 1 greater than value 2", 5, 4, false},
		{"value 1 equal value 2", 5, 5, true},
		{"value 1 less than value 2", 4, 5, true},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			result := Lte(tc.Value1, tc.Value2)

			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestGt(t *testing.T) {
	type testcase struct {
		Name     string
		Value1   float64
		Value2   float64
		Expected bool
	}

	tcs := []testcase{
		{"value 1 greater than value 2", 5, 4, true},
		{"value 1 equal value 2", 5, 5, false},
		{"value 1 less than value 2", 4, 5, false},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			result := Gt(tc.Value1, tc.Value2)

			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestOneof(t *testing.T) {
	type testcase struct {
		Name     string
		Value    string
		Items    []string
		Expected bool
	}

	tcs := []testcase{
		{"contains 1 value", "donation", []string{"personal", "donation", "k-receipt"}, true},
		{"does not contains", "shop", []string{"personal", "donation", "k-receipt"}, false},
		{"no values in slice", "personal", []string{}, false},
		{"empty value", "", []string{"personal", "donation", "k-receipt"}, false},
		{"empty both", "", []string{}, false},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			result := Oneof(tc.Value, tc.Items...)

			assert.Equal(t, tc.Expected, result)
		})
	}
}
