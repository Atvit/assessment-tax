package tax

func getMockTaxLevels(levelsToUpdate ...TaxLevel) []TaxLevel {
	levelElementPositionMap := map[string]int{
		level1: 0,
		level2: 1,
		level3: 2,
		level4: 3,
		level5: 4,
	}

	mockTaxLevels := []TaxLevel{
		{getLevelDescription(level1), new(float64)},
		{getLevelDescription(level2), new(float64)},
		{getLevelDescription(level3), new(float64)},
		{getLevelDescription(level4), new(float64)},
		{getLevelDescription(level5), new(float64)},
	}

	for _, v := range levelsToUpdate {
		mockTaxLevels[levelElementPositionMap[v.Level]].Tax = v.Tax
	}

	return mockTaxLevels
}
