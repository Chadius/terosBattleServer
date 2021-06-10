package levelup_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/usecase/levelup"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SquaddieQualifiesForClassSuite struct {
	classRepo *squaddieclass.Repository
	levelRepo *levelupbenefit.Repository

	mageClass *squaddieclass.Class
	mageLevel0 *levelupbenefit.LevelUpBenefit
	mageLevel1 *levelupbenefit.LevelUpBenefit

	dimensionWalkerClass *squaddieclass.Class
	dimensionWalkerLevel0 *levelupbenefit.LevelUpBenefit
	dimensionWalkerLevel1 *levelupbenefit.LevelUpBenefit

	ancientTomeClass *squaddieclass.Class
	ancientTomeClassLevel0 *levelupbenefit.LevelUpBenefit

	atLeastTenLevelsBaseClass *squaddieclass.Class
	lotsOfLevels []*levelupbenefit.LevelUpBenefit

	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieQualifiesForClassSuite{})

func (suite *SquaddieQualifiesForClassSuite) SetUpTest(checker *C) {
	suite.mageClass = &squaddieclass.Class{
		ID:                "class1",
		Name:              "Mage",
		BaseClassRequired: false,
	}

	suite.dimensionWalkerClass = &squaddieclass.Class{
		ID:                "class2",
		Name:              "Dimension Walker",
		BaseClassRequired: true,
	}

	suite.ancientTomeClass = &squaddieclass.Class{
		ID:                "class3",
		Name:              "Ancient Tome",
		BaseClassRequired: true,
	}

	suite.atLeastTenLevelsBaseClass = &squaddieclass.Class{
		ID:                "class4",
		Name:              "Base with many levels",
		BaseClassRequired: false,
	}

	suite.classRepo = squaddieclass.NewRepository()
	suite.classRepo.AddListOfClasses([]*squaddieclass.Class{suite.mageClass, suite.dimensionWalkerClass, suite.ancientTomeClass})

	suite.mageLevel0 = &levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            suite.mageClass.ID,
			ID:                 "mageLevel0",
		},
		Defense: &levelupbenefit.Defense{
			MaxHitPoints: 1,
			Dodge:        1,
			Deflect:      1,
			MaxBarrier:   1,
			Armor:        1,
		},
		Offense: &levelupbenefit.Offense{
			Aim:      1,
			Strength: 1,
			Mind:     1,
		},
		PowerIDGained:      nil,
		PowerIDLost:        nil,
		Movement:           nil,
	}

	suite.mageLevel1 = &levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            suite.mageClass.ID,
			ID:                 "mageLevel1",
		},
		Defense: &levelupbenefit.Defense{
			MaxHitPoints: 0,
			Dodge:        0,
			Deflect:      0,
			MaxBarrier:   0,
			Armor:        0,
		},
		Offense: &levelupbenefit.Offense{
			Aim:      0,
			Strength: 0,
			Mind:     0,
		},
		PowerIDGained:      nil,
		PowerIDLost:        nil,
		Movement:           &squaddie.Movement{
			Distance:  1,
			Type:      "",
			HitAndRun: false,
		},
	}

	suite.dimensionWalkerLevel0 = &levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            suite.dimensionWalkerClass.ID,
			ID:                 "dwLevel0",
		},
		Movement:           &squaddie.Movement{
			Distance:  1,
			Type:      "light",
			HitAndRun: false,
		},
	}

	suite.dimensionWalkerLevel1 = &levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            suite.dimensionWalkerClass.ID,
			ID:                 "dwLevel1",
		},
		Defense: &levelupbenefit.Defense{
			MaxHitPoints: 1,
		},
		Offense: &levelupbenefit.Offense{
			Aim:      0,
			Strength: 0,
			Mind:     0,
		},
	}

	suite.ancientTomeClassLevel0 = &levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            suite.ancientTomeClass.ID,
			ID:                 "ancientTomeLevel0",
		},
		Defense: &levelupbenefit.Defense{
			Dodge:   1,
			Deflect: 1,
		},
		Offense: &levelupbenefit.Offense{
			Mind: 1,
		},
	}

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		suite.mageLevel0, suite.mageLevel1, suite.dimensionWalkerLevel0, suite.dimensionWalkerLevel1, suite.ancientTomeClassLevel0,
	})

	suite.lotsOfLevels = append(
		(&testutility.LevelGenerator{
			Instructions: &testutility.LevelGeneratorInstruction{
				NumberOfLevels: 11,
				ClassID:        suite.atLeastTenLevelsBaseClass.ID,
				PrefixLevelID:  "lotsLevelsSmall",
				Type:           levelupbenefit.Small,
			},
		}).Build(),
		(&testutility.LevelGenerator{
			Instructions: &testutility.LevelGeneratorInstruction{
				NumberOfLevels: 11,
				ClassID:        suite.atLeastTenLevelsBaseClass.ID,
				PrefixLevelID:  "lotsLevelsBig",
				Type:           levelupbenefit.Big,
			},
		}).Build()...
	)

	suite.levelRepo.AddLevels(suite.lotsOfLevels)

	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.ClassProgress.AddClass(suite.mageClass)
	suite.teros.ClassProgress.AddClass(suite.dimensionWalkerClass)
	suite.teros.ClassProgress.AddClass(suite.ancientTomeClass)
	suite.teros.ClassProgress.AddClass(suite.atLeastTenLevelsBaseClass)
}

func (suite *SquaddieQualifiesForClassSuite) TestNewSquaddieCanSwitchToBaseClass(checker *C) {
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieCannotSwitchToCurrentClass(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieMustFinishBaseClassBeforeSwitching(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	levelup.ImproveSquaddie(suite.mageLevel0, suite.teros, nil)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieCanSwitchFromFinishedBaseToAdvancedClass(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	levelup.ImproveSquaddie(suite.mageLevel0, suite.teros, nil)
	levelup.ImproveSquaddie(suite.mageLevel1, suite.teros, nil)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieCanSwitchToNewAdvancedClassAtLevelTen(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	levelup.ImproveSquaddie(suite.mageLevel0, suite.teros, nil)
	levelup.ImproveSquaddie(suite.mageLevel1, suite.teros, nil)
	suite.teros.ClassProgress.SetClass(suite.dimensionWalkerClass.ID)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestCannotSwitchToCompletedClass(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	levelup.ImproveSquaddie(suite.mageLevel0, suite.teros, nil)
	levelup.ImproveSquaddie(suite.mageLevel1, suite.teros, nil)
	suite.teros.ClassProgress.SetClass(suite.dimensionWalkerClass.ID)
	levelup.ImproveSquaddie(suite.dimensionWalkerLevel0, suite.teros, nil)
	levelup.ImproveSquaddie(suite.dimensionWalkerLevel1, suite.teros, nil)
	suite.teros.ClassProgress.SetClass(suite.ancientTomeClass.ID)
	checker.Assert(suite.teros.ClassProgress.CurrentClass, Equals, suite.ancientTomeClass.ID)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestQualifyForAdvancedClassesWhenBaseClassCompletes(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	levelup.ImproveSquaddie(suite.mageLevel0, suite.teros, nil)
	levelup.ImproveSquaddie(suite.mageLevel1, suite.teros, nil)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieStaysInAdvancedClassUntilCompletion(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	levelup.ImproveSquaddie(suite.mageLevel0, suite.teros, nil)
	levelup.ImproveSquaddie(suite.mageLevel1, suite.teros, nil)
	suite.teros.ClassProgress.SetClass(suite.dimensionWalkerClass.ID)
	levelup.ImproveSquaddie(suite.dimensionWalkerLevel0, suite.teros, nil)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID, suite.classRepo, suite.levelRepo), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestCanSwitchClassAfterTenLevels(checker *C) {
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.atLeastTenLevelsBaseClass.ID)
	suite.teros.ClassProgress.SetClass(suite.atLeastTenLevelsBaseClass.ID)
	for index, _ := range [10]int{} {
		levelup.ImproveSquaddie(suite.lotsOfLevels[index], suite.teros, nil)
	}
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass()[suite.atLeastTenLevelsBaseClass.ID], Equals, 10)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
	checker.Assert(levelup.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID, suite.classRepo, suite.levelRepo), Equals, true)
}
