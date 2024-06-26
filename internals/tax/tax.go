package tax

import (
	"github.com/Atvit/assessment-tax/errs"
	"github.com/Atvit/assessment-tax/utils"
	"github.com/shopspring/decimal"
	"math"
)

const precision = 1

const (
	personal = "personal"
	donation = "donation"
	kReceipt = "k-receipt"
)

const (
	level1 = "level1"
	level2 = "level2"
	level3 = "level3"
	level4 = "level4"
	level5 = "level5"
)

const (
	defaultPersonalAllowance = 60000.00
	defaultKReceiptAllowance = 50000.00
	maxDonationAllowance     = 100000.00
)

type TaxLevelMap map[string]TaxLevel

type TaxLevel struct {
	Level string   `json:"level"`
	Tax   *float64 `json:"tax"`
}

type AllowanceSetting struct {
	Personal float64
	KReceipt float64
}

type Allowance struct {
	AllowanceType string
	Amount        float64
}

type Tax struct {
	Income           float64
	Wht              float64
	Allowances       []Allowance
	AllowanceSetting AllowanceSetting
}

var Calculate = func(t *Tax) (float64, float64, []TaxLevel, error) {
	err := validate(t)
	if err != nil {
		return 0, 0, nil, err
	}

	addPersonalAllowance(t)
	deductAmount := getDeductAmount(t.Allowances, t.AllowanceSetting)
	taxableIncome := t.Income - deductAmount

	taxAmount, refundAmount, taxLevels := calculateTax(taxableIncome, t.Wht)

	return taxAmount, refundAmount, taxLevels, nil
}

func addPersonalAllowance(t *Tax) {
	personalAllowance := t.AllowanceSetting.Personal
	if decimal.NewFromFloat(personalAllowance).IsZero() {
		personalAllowance = defaultPersonalAllowance
	}

	t.Allowances = append(t.Allowances, Allowance{
		AllowanceType: "personal",
		Amount:        personalAllowance,
	})
}

func calculateTax(taxableIncome, wht float64) (float64, float64, []TaxLevel) {
	taxAmount := 0.0
	refundAmount := 0.0

	taxAmount, taxLevels := calculateProgressiveTax(taxableIncome)
	taxAmount -= wht

	if taxAmount < 0 {
		refundAmount = math.Abs(taxAmount)
		taxAmount = 0
	}

	return taxAmount, utils.Round(refundAmount, precision), taxLevels
}

func calculateProgressiveTax(taxableIncome float64) (float64, []TaxLevel) {
	taxAmount := 0.0
	taxLevelsMap := initializeTaxLevelsMap()
	taxLevels := getTaxLevels(taxLevelsMap)

	taxBrackets := []struct {
		lowerBound   float64
		upperBound   float64
		rate         float64
		level        string
		noUpperLimit bool
	}{
		{150000, 500000, 0.10, level2, false},
		{500000, 1000000, 0.15, level3, false},
		{1000000, 2000000, 0.20, level4, false},
		{2000000, 0, 0.35, level5, true},
	}

	for _, bracket := range taxBrackets {
		if taxableIncome > bracket.lowerBound {
			tax := 0.0
			if bracket.noUpperLimit {
				tax = (taxableIncome - bracket.lowerBound) * bracket.rate
			} else {
				tax = calculateTaxBracket(taxableIncome, bracket.lowerBound, bracket.upperBound, bracket.rate)
			}

			taxAmount += tax
			updateTaxLevel(taxLevelsMap, bracket.level, utils.Round(tax, precision))
		}
	}

	return utils.Round(taxAmount, precision), taxLevels
}

func calculateTaxBracket(income, lower, upper float64, rate float64) float64 {
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

func getDeductAmount(allowances []Allowance, setting AllowanceSetting) float64 {
	amount := 0.00

	for _, allowance := range allowances {
		if allowance.AllowanceType == donation {
			if allowance.Amount > maxDonationAllowance {
				allowance.Amount = maxDonationAllowance
			}
		}

		if allowance.AllowanceType == kReceipt {
			if decimal.NewFromFloat(setting.KReceipt).IsZero() {
				setting.KReceipt = defaultKReceiptAllowance
			}

			if allowance.Amount > setting.KReceipt {
				allowance.Amount = setting.KReceipt
			}
		}

		amount += allowance.Amount
	}

	return amount
}

func initializeTaxLevelsMap() TaxLevelMap {
	return map[string]TaxLevel{
		level1: {Level: getLevelDescription(level1), Tax: new(float64)},
		level2: {Level: getLevelDescription(level2), Tax: new(float64)},
		level3: {Level: getLevelDescription(level3), Tax: new(float64)},
		level4: {Level: getLevelDescription(level4), Tax: new(float64)},
		level5: {Level: getLevelDescription(level5), Tax: new(float64)},
	}
}

func updateTaxLevel(taxLevels TaxLevelMap, level string, amount float64) {
	if taxLevel, ok := taxLevels[level]; ok {
		*taxLevel.Tax += amount
	}
}

func getLevelDescription(level string) string {
	levelDescriptionMap := map[string]string{
		level1: "0-150,000",
		level2: "150,001-500,000",
		level3: "500,001-1,000,000",
		level4: "1,000,001-2,000,000",
		level5: "2,000,001 ขึ้นไป",
	}

	if description, ok := levelDescriptionMap[level]; ok {
		return description
	}

	return ""
}

func getTaxLevels(taxLevels TaxLevelMap) []TaxLevel {
	var levels []TaxLevel
	keys := []string{level1, level2, level3, level4, level5}
	for _, k := range keys {
		levels = append(levels, taxLevels[k])
	}

	return levels
}
