package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	. "gopkg.in/check.v1"
)

type LevelUpBenefitSuite struct{}

var _ = Suite(&LevelUpBenefitSuite{})

func (s *LevelUpBenefitSuite) TestRaisesAnErrorIfNoBenefitType(checker *C) {
	badLevel := levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{ClassID: "class0"},
	}
	err := badLevel.CheckForErrors()
	checker.Assert(err, ErrorMatches, `unknown level up benefit type: ""`)
}

func (s *LevelUpBenefitSuite) TestRaisesAnErrorIfNoClassID(checker *C) {
	badLevel := levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            "",
		},
	}
	err := badLevel.CheckForErrors()
	checker.Assert(err, ErrorMatches, `no classID found for LevelUpBenefit`)
}

func (s *LevelUpBenefitSuite) TestFiltersAList(checker *C) {
	listToTest := []*levelupbenefit.LevelUpBenefit{
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            "class0",
				ID:                 "level0",
			},
			Offense: &levelupbenefit.Offense{
				Aim: 1,
			},
		},
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            "class0",
				ID:                 "level1",
			},
			Defense: &levelupbenefit.Defense{
				MaxHitPoints: 1,
			},
		},
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            "class0",
				ID:                 "level2",
			},
			Offense: &levelupbenefit.Offense{
				Aim: 1,
			},
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
		return benefit.Identification.LevelUpBenefitType == levelupbenefit.Small
	})
	checker.Assert(onlySmallLevels, HasLen, 2)

	onlyBigLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.Identification.LevelUpBenefitType == levelupbenefit.Big
	})
	checker.Assert(onlyBigLevels, HasLen, 1)

	increasesAimLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.Offense != nil && benefit.Offense.Aim > 0
	})
	checker.Assert(increasesAimLevels, HasLen, 2)
}