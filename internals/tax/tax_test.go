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
		wantRefund float64
		wantErr    error
		allowances []TaxAllowance
	}{
		{
			name:       "negative income",
			income:     -1000,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    errs.ErrValueMustBePositive,
		},
		{
			name:       "zero income",
			income:     0,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income below personal allowance",
			income:     50000,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income equal to personal allowance",
			income:     60000,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income above personal allowance",
			income:     150001,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income at upper bracket limit",
			income:     500000,
			wantTax:    29000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income above upper bracket limit",
			income:     500001,
			wantTax:    29000.1,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income at second upper bracket limit",
			income:     1000000,
			wantTax:    101000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income above second upper bracket limit",
			income:     1000001,
			wantTax:    101000.2,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income at third upper bracket limit",
			income:     2000000,
			wantTax:    298000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "income above third upper bracket limit",
			income:     2000001,
			wantTax:    298000.2,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "high income",
			income:     3000000,
			wantTax:    639000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "with holding tax",
			income:     500000,
			wht:        25000,
			wantTax:    4000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "negative with holding tax",
			income:     500000,
			wht:        -25000,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    errs.ErrValueMustBePositive,
		},
		{
			name:       "with holding tax greater than income",
			income:     500000,
			wht:        500001,
			wantTax:    0,
			wantRefund: 0,
			wantErr:    errs.ErrWhtMustLowerThanOrEqualIncome,
		},
		{
			name:       "invalid allowance type",
			income:     500000,
			wht:        0,
			wantTax:    0,
			wantRefund: 0,
			allowances: []TaxAllowance{
				{AllowanceType: personal, Amount: 40000},
				{AllowanceType: "invalid", Amount: 50000},
			},
			wantErr: errs.ErrIncorrectAllowanceType,
		},
		{
			name:       "allowance amount less than zero",
			income:     500000,
			wht:        0,
			allowances: []TaxAllowance{{AllowanceType: personal, Amount: -1000}},
			wantTax:    0,
			wantRefund: 0,
			wantErr:    errs.ErrValueMustBePositive,
		},
		{
			name:       "empty allowance slice",
			income:     500000,
			wht:        0,
			allowances: []TaxAllowance{},
			wantTax:    29000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "donation allowance",
			income:     500000,
			wht:        0,
			allowances: []TaxAllowance{{AllowanceType: donation, Amount: 200000}},
			wantTax:    19000,
			wantRefund: 0,
			wantErr:    nil,
		},
		{
			name:       "get a refund if tax-exempt and have withholding tax",
			income:     150000,
			wht:        1000,
			wantTax:    0,
			wantRefund: 1000,
		},
		{
			name:       "get a refund if withholding tax is greater than tax to pay",
			income:     500000,
			wht:        30000,
			allowances: []TaxAllowance{{AllowanceType: donation, Amount: 200000}},
			wantTax:    0,
			wantRefund: 11000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax, gotRefund, gotErr := Calculate(&Tax{
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
			if gotRefund != tt.wantRefund {
				t.Errorf("calculateTax = %v, want %v", gotRefund, tt.wantRefund)
			}
		})
	}
}
