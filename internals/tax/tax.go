package tax

import (
	"github.com/Atvit/assessment-tax/errs"
)

var calculate = func(income float64) (float64, error) {
	if income < 0 {
		return 0, errs.ErrValueMustBePositive
	}

	personalAllowance := 60000.00

	if income <= 150000 {
		return 0, nil
	}

	taxableIncome := income - personalAllowance
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

	return taxAmount, nil
}
