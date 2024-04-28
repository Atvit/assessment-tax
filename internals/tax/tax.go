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
	defaultPersonalAllowance := 60000.00
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

	if taxableIncome <= 150000 {
		return 0, taxLevels
	}

	if taxableIncome > 150000 {
		tax := calculateTaxBracket(taxableIncome, 150000, 500000, 0.10)

		taxAmount += tax
		updateTaxLevel(taxLevelsMap, level2, utils.Round(tax, precision))
	}

	if taxableIncome > 500000 {
		tax := calculateTaxBracket(taxableIncome, 500000, 1000000, 0.15)

		taxAmount += tax
		updateTaxLevel(taxLevelsMap, level3, utils.Round(tax, precision))
	}
	if taxableIncome > 1000000 {
		tax := calculateTaxBracket(taxableIncome, 1000000, 2000000, 0.20)

		taxAmount += tax
		updateTaxLevel(taxLevelsMap, level4, utils.Round(tax, precision))
	}
	if taxableIncome > 2000000 {
		tax := (taxableIncome - 2000000) * 0.35

		taxAmount += tax
		updateTaxLevel(taxLevelsMap, level5, utils.Round(tax, precision))
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
			if allowance.Amount > 100000 {
				allowance.Amount = 100000
			}
		}

		if allowance.AllowanceType == kReceipt {
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
	if taxLevel, exists := taxLevels[level]; exists {
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
