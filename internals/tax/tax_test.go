package tax

import (
	"github.com/Atvit/assessment-tax/errs"
	"github.com/Atvit/assessment-tax/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name             string
		income           float64
		wht              float64
		expectedTax      float64
		expectedRefund   float64
		expectedErr      error
		expectedLevels   []TaxLevel
		allowances       []Allowance
		allowanceSetting AllowanceSetting
	}{
		{
			name:           "negative income",
			income:         -1000,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    errs.ErrValueMustBePositive,
		},
		{
			name:           "zero income",
			income:         0,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(),
		},
		{
			name:           "income below personal allowance",
			income:         50000,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(),
		},
		{
			name:           "income equal to personal allowance",
			income:         60000,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(),
		},
		{
			name:           "income above personal allowance",
			income:         150001,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(),
		},
		{
			name:           "income at upper bracket limit",
			income:         500000,
			expectedTax:    29000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(TaxLevel{Level: level2, Tax: utils.ToPointer(29000.0)}),
		},
		{
			name:           "income above upper bracket limit",
			income:         500001,
			expectedTax:    29000.1,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(TaxLevel{Level: level2, Tax: utils.ToPointer(29000.1)}),
		},
		{
			name:           "income at second upper bracket limit",
			income:         1000000,
			expectedTax:    101000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(
				TaxLevel{level2, utils.ToPointer(35000.0)},
				TaxLevel{level3, utils.ToPointer(66000.0)},
			),
		},
		{
			name:           "income above second upper bracket limit",
			income:         1000001,
			expectedTax:    101000.2,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(
				TaxLevel{level2, utils.ToPointer(35000.0)},
				TaxLevel{level3, utils.ToPointer(66000.2)},
			),
		},
		{
			name:           "income at third upper bracket limit",
			income:         2000000,
			expectedTax:    298000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(
				TaxLevel{level2, utils.ToPointer(35000.0)},
				TaxLevel{level3, utils.ToPointer(75000.0)},
				TaxLevel{level4, utils.ToPointer(188000.0)},
			),
		},
		{
			name:           "income above third upper bracket limit",
			income:         2000001,
			expectedTax:    298000.2,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(
				TaxLevel{level2, utils.ToPointer(35000.0)},
				TaxLevel{level3, utils.ToPointer(75000.0)},
				TaxLevel{level4, utils.ToPointer(188000.2)},
			),
		},
		{
			name:           "high income",
			income:         3000000,
			expectedTax:    639000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(
				TaxLevel{level2, utils.ToPointer(35000.0)},
				TaxLevel{level3, utils.ToPointer(75000.0)},
				TaxLevel{level4, utils.ToPointer(200000.0)},
				TaxLevel{level5, utils.ToPointer(329000.0)},
			),
		},
		{
			name:           "with holding tax",
			income:         500000,
			wht:            25000,
			expectedTax:    4000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(TaxLevel{level2, utils.ToPointer(29000.0)}),
		},
		{
			name:           "negative with holding tax",
			income:         500000,
			wht:            -25000,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    errs.ErrValueMustBePositive,
		},
		{
			name:           "with holding tax greater than income",
			income:         500000,
			wht:            500001,
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    errs.ErrWhtMustLowerThanOrEqualIncome,
		},
		{
			name:           "invalid allowance type",
			income:         500000,
			wht:            0,
			expectedTax:    0,
			expectedRefund: 0,
			allowances:     []Allowance{{personal, 40000}, {"invalid", 50000}},
			expectedErr:    errs.ErrIncorrectAllowanceType,
		},
		{
			name:           "allowance amount less than zero",
			income:         500000,
			wht:            0,
			allowances:     []Allowance{{personal, -1000}},
			expectedTax:    0,
			expectedRefund: 0,
			expectedErr:    errs.ErrValueMustBePositive,
		},
		{
			name:           "empty allowance slice",
			income:         500000,
			wht:            0,
			allowances:     []Allowance{},
			expectedTax:    29000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(TaxLevel{level2, utils.ToPointer(29000.0)}),
		},
		{
			name:           "donation allowance",
			income:         500000,
			wht:            0,
			allowances:     []Allowance{{donation, 200000}},
			expectedTax:    19000,
			expectedRefund: 0,
			expectedErr:    nil,
			expectedLevels: getMockTaxLevels(TaxLevel{level2, utils.ToPointer(19000.0)}),
		},
		{
			name:           "get a refund if tax-exempt and have withholding tax",
			income:         150000,
			wht:            1000,
			expectedTax:    0,
			expectedRefund: 1000,
			expectedLevels: getMockTaxLevels(),
		},
		{
			name:           "get a refund if withholding tax is greater than tax to pay",
			income:         500000,
			wht:            30000,
			allowances:     []Allowance{{donation, 200000}},
			expectedTax:    0,
			expectedRefund: 11000,
			expectedLevels: getMockTaxLevels(TaxLevel{level2, utils.ToPointer(19000.0)}),
		},
		{
			name:       "k-receipt allowance",
			income:     500000,
			wht:        0,
			allowances: []Allowance{{kReceipt, 200000}, {donation, 100000}},
			allowanceSetting: AllowanceSetting{
				Personal: 60000.0,
				KReceipt: 50000.0,
			},
			expectedTax:    14000,
			expectedRefund: 0,
			expectedLevels: getMockTaxLevels(TaxLevel{level2, utils.ToPointer(14000.0)}),
		},
		{
			name:           "default k-receipt allowance",
			income:         500000,
			wht:            0,
			allowances:     []Allowance{{kReceipt, 200000}, {donation, 100000}},
			expectedTax:    14000,
			expectedRefund: 0,
			expectedLevels: getMockTaxLevels(TaxLevel{level2, utils.ToPointer(14000.0)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taxAmount, refundAmount, taxLevels, err := Calculate(&Tax{
				Income:           tt.income,
				Wht:              tt.wht,
				Allowances:       tt.allowances,
				AllowanceSetting: tt.allowanceSetting,
			})

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTax, taxAmount)
			assert.Equal(t, tt.expectedRefund, refundAmount)

			for i, level := range taxLevels {
				assert.Equal(t, *tt.expectedLevels[i].Tax, *level.Tax)
			}
		})
	}
}
