package tax

import (
	"github.com/Atvit/assessment-tax/errs"
)

type Tax struct {
	Income float64
	Wht    float64
}

func validateInput(t *Tax) error {
	if t.Income < 0 {
		return errs.ErrValueMustBePositive
	}

	if t.Wht < 0 {
		return errs.ErrValueMustBePositive
	}

	if t.Wht > t.Income {
		return errs.ErrWhtMustLowerThanOrEqualIncome
	}

	return nil
}

var Calculate = func(t *Tax) (float64, error) {
	err := validateInput(t)
	if err != nil {
		return 0, err
	}

	personalAllowance := 60000.00

	if t.Income <= 150000 {
		return 0, nil
	}

	taxableIncome := t.Income - personalAllowance
	taxAmount := 0.0

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

	return taxAmount, nil
}
