package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRound(t *testing.T) {
	type testcase struct {
		Name      string
		Value     float64
		Precision int
		Expected  float64
	}

	tcs := []testcase{
		{"round positive float", 54.23456, 2, 54.23},
		{"round negative float", -54.761, 2, -54.76},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			result := Round(tc.Value, tc.Precision)

			assert.Equal(t, tc.Expected, result)
		})
	}
}
