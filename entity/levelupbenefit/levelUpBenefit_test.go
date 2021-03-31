package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/kindrid/gotest"
	"github.com/kindrid/gotest/should"
	"testing"
)

func TestRaisesAnErrorIfNoBenefitType(t *testing.T) {
	badLevel := levelupbenefit.LevelUpBenefit{
		ClassID:            "class0",
	}
	err := badLevel.CheckForErrors()
	gotest.Assert(t, err.Error(), should.Equal, `unknown level up benefit type: ""`)
}

func TestRaisesAnErrorIfNoClassID(t *testing.T) {
	badLevel := levelupbenefit.LevelUpBenefit{
		LevelUpBenefitType: levelupbenefit.Small,
		ClassID:            "",
	}
	err := badLevel.CheckForErrors()
	gotest.Assert(t, err.Error(), should.Equal, `no classID found for LevelUpBenefit`)
}

func TestFiltersAList(t *testing.T) {
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
	gotest.Assert(t, noLevelsFound, should.HaveLength, 0)

	allLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return true
	})
	gotest.Assert(t, allLevelsFound, should.HaveLength, 3)


	onlySmallLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.LevelUpBenefitType == levelupbenefit.Small
	})
	gotest.Assert(t, onlySmallLevels, should.HaveLength, 2)

	onlyBigLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.LevelUpBenefitType == levelupbenefit.Big
	})
	gotest.Assert(t, onlyBigLevels, should.HaveLength, 1)

	increasesAimLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.Aim > 0
	})
	gotest.Assert(t, increasesAimLevels, should.HaveLength, 2)
}