package squaddie

// ClassLevelsConsumed tracks information about how a Squaddie uses a class as well
//    as the progress made in a given class.
type ClassLevelsConsumed struct {
	ClassID string `json:"id" yaml:"id"`
	ClassName string `json:"name" yaml:"name"`
	LevelsConsumed []string `json:"levels_gained" yaml:"levels_gained"`
}

// IsLevelAlreadyConsumed returns true if the level ID has already been used.
func (progress *ClassLevelsConsumed) IsLevelAlreadyConsumed(levelUpBenefitLevelID string) bool {
	return progress.AnyLevelsConsumed(func(consumedLevelID string) bool {
		return consumedLevelID == levelUpBenefitLevelID
	})
}

// AnyLevelsConsumed returns true if at least 1 levelID satisfies the condition.
func (progress *ClassLevelsConsumed) AnyLevelsConsumed(condition func(consumedLevelID string) bool) bool {
	for _, levelID := range progress.LevelsConsumed {
		if condition(levelID) {
			return true
		}
	}
	return false
}

// AccumulateLevelsConsumed calls the calculate function on each Level consumed and adds it to a sum.
//   The sum is returned after processing all of the levels.
func (progress *ClassLevelsConsumed) AccumulateLevelsConsumed(calculate func(consumedLevelID string) int) int {
	count := 0
	for _, levelID := range progress.LevelsConsumed {
		count = count + calculate(levelID)
	}
	return count
}