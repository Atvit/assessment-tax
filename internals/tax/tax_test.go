package tax

import (
	"github.com/Atvit/assessment-tax/errs"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name       string
		income     float64
		wht        float64
		wantTax    float64
		wantErr    error
		allowances []TaxAllowance
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
			wantTax: 29000.1,
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
			wantTax: 101000.2,
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
		{
			name:    "With Holding Tax",
			income:  500000,
			wht:     25000,
			wantTax: 4000,
			wantErr: nil,
		},
		{
			name:    "Negative With Holding Tax",
			income:  500000,
			wht:     -25000,
			wantTax: 0,
			wantErr: errs.ErrValueMustBePositive,
		},
		{
			name:    "With Holding Tax Greater Than Income",
			income:  500000,
			wht:     500001,
			wantTax: 0,
			wantErr: errs.ErrWhtMustLowerThanOrEqualIncome,
		},
		{
			name:    "Invalid Allowance Type",
			income:  500000,
			wht:     0,
			wantTax: 0,
			allowances: []TaxAllowance{
				{AllowanceType: personal, Amount: 40000},
				{AllowanceType: "invalid", Amount: 50000},
			},
			wantErr: errs.ErrIncorrectAllowanceType,
		},
		{
			name:       "Allowance Amount Equal Zero",
			income:     500000,
			wht:        0,
			wantTax:    0,
			allowances: []TaxAllowance{{AllowanceType: personal, Amount: 0}},
			wantErr:    errs.ErrValueMustBePositive,
		},
		{
			name:       "Allowance Amount Less Than Zero",
			income:     500000,
			wht:        0,
			allowances: []TaxAllowance{{AllowanceType: personal, Amount: -1000}},
			wantTax:    0,
			wantErr:    errs.ErrValueMustBePositive,
		},
		{
			name:       "Empty Allowance Slice",
			income:     500000,
			wht:        0,
			allowances: []TaxAllowance{},
			wantTax:    29000,
			wantErr:    nil,
		},
		{
			name:       "With Donation Allowance",
			income:     500000,
			wht:        0,
			allowances: []TaxAllowance{{AllowanceType: donation, Amount: 200000}},
			wantTax:    19000,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax, gotErr := Calculate(&Tax{
				Income:     tt.income,
				Wht:        tt.wht,
				Allowances: tt.allowances,
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
