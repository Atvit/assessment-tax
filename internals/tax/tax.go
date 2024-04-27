package tax

import (
	"github.com/Atvit/assessment-tax/errs"
	"github.com/Atvit/assessment-tax/utils"
	"math"
)

const (
	personal = "personal"
	donation = "donation"
	kReceipt = "k-receipt"
)

type TaxAllowance struct {
	AllowanceType string
	Amount        float64
}

type Tax struct {
	Income     float64
	Wht        float64
	Allowances []TaxAllowance
}

const precision = 1

func round(num float64, precision int) float64 {
	power := math.Pow10(precision)
	rounded := math.Round(num*power) / power
	return rounded
}

var Calculate = func(t *Tax) (float64, float64, error) {
	err := validate(t)
	if err != nil {
		return 0, 0, err
	}

	personalAllowance := 60000.00
	t.Allowances = append(t.Allowances, TaxAllowance{
		AllowanceType: "personal",
		Amount:        personalAllowance,
	})
	deductAmount := getDeductAmount(t.Allowances)
	taxableIncome := t.Income - deductAmount
	taxAmount := 0.0
	refundAmount := 0.0

	if t.Income <= 150000 {
		return 0, t.Wht, nil
	}

	if taxableIncome > 150000 {
		if taxableIncome <= 500000 {
			taxAmount += (taxableIncome - 150000) * 0.10
		} else {
			taxAmount += (500000 - 150000) * 0.10
		}
	}

	if taxableIncome > 500000 {
		if taxableIncome <= 1000000 {
			taxAmount += (taxableIncome - 500000) * 0.15
		} else {
			taxAmount += (1000000 - 500000) * 0.15
		}
	}

	if taxableIncome > 1000000 {
		if taxableIncome <= 2000000 {
			taxAmount += (taxableIncome - 1000000) * 0.20
		} else {
			taxAmount += (2000000 - 1000000) * 0.20
		}
	}

	if taxableIncome > 2000000 {
		taxAmount += (taxableIncome - 2000000) * 0.35
	}

	taxAmount -= t.Wht
	if taxAmount < 0 {
		refundAmount = math.Abs(taxAmount)
		taxAmount = 0
	}

	return round(taxAmount, precision), round(refundAmount, precision), nil
}

func validate(t *Tax) error {
	if ok := utils.Gte(t.Income, 0); !ok {
		return errs.ErrValueMustBePositive
	}

	if ok := utils.Gte(t.Wht, 0); !ok {
		return errs.ErrValueMustBePositive
	}

	if ok := utils.Lte(t.Wht, t.Income); !ok {
		return errs.ErrWhtMustLowerThanOrEqualIncome
	}

	for _, allowance := range t.Allowances {
		if ok := utils.Gte(allowance.Amount, 0); !ok {
			return errs.ErrValueMustBePositive
		}

		if ok := utils.Oneof(allowance.AllowanceType, personal, donation, kReceipt); !ok {
			return errs.ErrIncorrectAllowanceType
		}
	}

	return nil
}

func getDeductAmount(allowances []TaxAllowance) float64 {
	amount := 0.00

	for _, allowance := range allowances {
		if allowance.AllowanceType == donation {
			if allowance.Amount > 100000 {
				allowance.Amount = 100000
			}
		}

		amount += allowance.Amount
	}

	return amount
}
