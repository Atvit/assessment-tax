package tax

import (
	"github.com/Atvit/assessment-tax/errs"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name    string
		income  float64
		wantTax float64
		wantErr error
	}{
		{
			name:    "Negative Income",
			income:  -1000,
			wantTax: 0,
			wantErr: errs.ErrValueMustBePositive,
		},
		{
			name:    "Zero Income",
			income:  0,
			wantTax: 0,
			wantErr: nil,
		},
		{
			name:    "Income Below Personal Allowance",
			income:  50000,
			wantTax: 0,
			wantErr: nil,
		},
		{
			name:    "Income Equal to Personal Allowance",
			income:  60000,
			wantTax: 0,
			wantErr: nil,
		},
		{
			name:    "Income Above Personal Allowance",
			income:  150001,
			wantTax: 0,
			wantErr: nil,
		},
		{
			name:    "Income at Upper Bracket Limit",
			income:  500000,
			wantTax: 29000,
			wantErr: nil,
		},
		{
			name:    "Income Above Upper Bracket Limit",
			income:  500001,
			wantTax: 29000.100000000002,
			wantErr: nil,
		},
		{
			name:    "Income at Second Upper Bracket Limit",
			income:  1000000,
			wantTax: 101000,
			wantErr: nil,
		},
		{
			name:    "Income Above Second Upper Bracket Limit",
			income:  1000001,
			wantTax: 101000.15,
			wantErr: nil,
		},
		{
			name:    "Income at Third Upper Bracket Limit",
			income:  2000000,
			wantTax: 298000,
			wantErr: nil,
		},
		{
			name:    "Income Above Third Upper Bracket Limit",
			income:  2000001,
			wantTax: 298000.2,
			wantErr: nil,
		},
		{
			name:    "High Income",
			income:  3000000,
			wantTax: 639000,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax, gotErr := Calculate(&Tax{
				Income: tt.income,
			})
			if gotErr != tt.wantErr {
				t.Errorf("calculateTax error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if gotTax != tt.wantTax {
				t.Errorf("calculateTax = %v, want %v", gotTax, tt.wantTax)
			}
		})
	}
}
