package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	. "gopkg.in/check.v1"
)

type LevelUpBenefitSuite struct{}

var _ = Suite(&LevelUpBenefitSuite{})

func (s *LevelUpBenefitSuite) TestRaisesAnErrorIfNoBenefitType(checker *C) {
	badLevel := levelupbenefit.LevelUpBenefit{
		ClassID:            "class0",
	}
	err := badLevel.CheckForErrors()
	checker.Assert(err, ErrorMatches, `unknown level up benefit type: ""`)
}

func (s *LevelUpBenefitSuite) TestRaisesAnErrorIfNoClassID(checker *C) {
	badLevel := levelupbenefit.LevelUpBenefit{
		LevelUpBenefitType: levelupbenefit.Small,
		ClassID:            "",
	}
	err := badLevel.CheckForErrors()
	checker.Assert(err, ErrorMatches, `no classID found for LevelUpBenefit`)
}

func (s *LevelUpBenefitSuite) TestFiltersAList(checker *C) {
	listToTest := []*levelupbenefit.LevelUpBenefit{
		{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            "class0",
			ID:                 "level0",
			Aim:                1,
		},
		{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            "class0",
			ID:                 "level1",
			MaxHitPoints:       1,
		},
		{
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            "class0",
			ID:                 "level2",
			Aim:                1,
		},
	}

	noLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return false
	})
	checker.Assert(noLevelsFound, HasLen, 0)

	allLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return true
	})
	checker.Assert(allLevelsFound, HasLen, 3)


	onlySmallLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.LevelUpBenefitType == levelupbenefit.Small
	})
	checker.Assert(onlySmallLevels, HasLen, 2)

	onlyBigLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.LevelUpBenefitType == levelupbenefit.Big
	})
	checker.Assert(onlyBigLevels, HasLen, 1)

	increasesAimLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.Aim > 0
	})
	checker.Assert(increasesAimLevels, HasLen, 2)
}