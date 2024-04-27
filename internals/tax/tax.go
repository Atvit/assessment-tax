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

var Calculate = func(t *Tax) (float64, float64, error) {
	err := validate(t)
	if err != nil {
		return 0, 0, err
	}

	addPersonalAllowance(t)
	deductAmount := getDeductAmount(t.Allowances)
	taxableIncome := t.Income - deductAmount

	taxAmount, refundAmount := calculateTax(taxableIncome, t.Wht)

	return round(taxAmount, precision), round(refundAmount, precision), nil
}

func addPersonalAllowance(t *Tax) {
	personalAllowance := 60000.00
	t.Allowances = append(t.Allowances, TaxAllowance{
		AllowanceType: "personal",
		Amount:        personalAllowance,
	})
}

func calculateTax(taxableIncome, wht float64) (float64, float64) {
	taxAmount := 0.0
	refundAmount := 0.0

	if taxableIncome <= 150000 {
		return 0, wht
	}

	taxAmount = calculateProgressiveTax(taxableIncome)
	taxAmount -= wht

	if taxAmount < 0 {
		refundAmount = math.Abs(taxAmount)
		taxAmount = 0
	}

	return taxAmount, refundAmount
}

func calculateProgressiveTax(taxableIncome float64) float64 {
	taxAmount := 0.0

	if taxableIncome > 150000 {
		taxAmount += calculateTaxAmount(taxableIncome, 150000, 500000, 0.10)
	}
	if taxableIncome > 500000 {
		taxAmount += calculateTaxAmount(taxableIncome, 500000, 1000000, 0.15)
	}
	if taxableIncome > 1000000 {
		taxAmount += calculateTaxAmount(taxableIncome, 1000000, 2000000, 0.20)
	}
	if taxableIncome > 2000000 {
		taxAmount += (taxableIncome - 2000000) * 0.35
	}

	return taxAmount
}

func calculateTaxAmount(income, lower, upper float64, rate float64) float64 {
	if income <= upper {
		return (income - lower) * rate
	}
	return (upper - lower) * rate
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

func round(num float64, precision int) float64 {
	power := math.Pow10(precision)
	rounded := math.Round(num*power) / power
	return rounded
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
