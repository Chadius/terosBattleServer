package levelup_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/usecase/levelup"
	. "gopkg.in/check.v1"
)

type SquaddieUsesLevelUpBenefitSuite struct {
	mageClass *squaddieclass.Class
	statBooster levelupbenefit.LevelUpBenefit
	teros *squaddie.Squaddie
	improveAllMovement *levelupbenefit.LevelUpBenefit
	upgradeToLightMovement *levelupbenefit.LevelUpBenefit
}

var _ = Suite(&SquaddieUsesLevelUpBenefitSuite{})

func (suite *SquaddieUsesLevelUpBenefitSuite) SetUpTest(checker *C) {
	suite.mageClass = &squaddieclass.Class{
		ID:   "ffffffff",
		Name: "Mage",
	}
	suite.teros = squaddie.NewSquaddie("suite.teros")
	suite.teros.MaxHitPoints = 5
	suite.teros.Aim = 0
	suite.teros.Strength = 1
	suite.teros.Mind = 2
	suite.teros.Dodge = 3
	suite.teros.Deflect = 4
	suite.teros.MaxBarrier = 6
	suite.teros.Armor = 7
	suite.teros.AddClass(suite.mageClass)
	suite.teros.SetHPToMax()
	suite.teros.SetBarrierToMax()

	suite.statBooster = levelupbenefit.LevelUpBenefit{
		ID:           "deadbeef",
		ClassID:      suite.mageClass.ID,
		MaxHitPoints: 0,
		Aim :         7,
		Strength :    6,
		Mind :        5,
		Dodge :       4,
		Deflect :     3,
		MaxBarrier :  2,
		Armor :       1,
	}

	suite.improveAllMovement = &levelupbenefit.LevelUpBenefit{
		ID:      "aaaaaaa0",
		ClassID: suite.mageClass.ID,
		Movement: &squaddie.Movement{
			Distance: 1,
			Type: "fly",
			HitAndRun: true,
		},
	}

	suite.upgradeToLightMovement = &levelupbenefit.LevelUpBenefit{
		ID:      "aaaaaaa1",
		ClassID: suite.mageClass.ID,
		Movement: &squaddie.Movement{
			Type: "light",
		},
	}
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestIncreaseStats(checker *C) {
	err := levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.MaxHitPoints, Equals, 5)
	checker.Assert(suite.teros.Aim, Equals, 7)
	checker.Assert(suite.teros.Strength, Equals, 7)
	checker.Assert(suite.teros.Mind, Equals, 7)
	checker.Assert(suite.teros.Dodge, Equals, 7)
	checker.Assert(suite.teros.Deflect, Equals, 7)
	checker.Assert(suite.teros.MaxBarrier, Equals, 8)
	checker.Assert(suite.teros.Armor, Equals, 8)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieRecordsLevel(checker *C) {
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID), Equals, false)
	err := levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID: 1})
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorForNonexistentClass(checker *C) {
	mushroomClassLevel := levelupbenefit.LevelUpBenefit{
		ID:           "deedbeeg",
		ClassID:      "bad ID",
		MaxHitPoints: 0,
		Aim :         7,
		Strength :    6,
		Mind :        5,
		Dodge :       4,
		Deflect :     3,
		MaxBarrier :  2,
		Armor :       1,
	}
	err := levelup.ImproveSquaddie(&mushroomClassLevel, suite.teros, nil)
	checker.Assert(err.Error(), Equals, `squaddie "suite.teros" cannot add levels to unknown class "bad ID"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorIfReusingLevel(checker *C) {
	err := levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"ffffffff": 1})
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID), Equals, true)

	err = levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err.Error(), Equals, `suite.teros already consumed LevelUpBenefit - class:"ffffffff" id:"deadbeef"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestUsingLevelSetsBaseClassIfBaseClassIsUnset(checker *C) {
	checker.Assert(suite.teros.BaseClassID, Equals, "")
	levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(suite.teros.BaseClassID, Equals, suite.mageClass.ID)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieChangeMovement(checker *C) {
	startingMovement := suite.teros.GetMovementDistancePerRound()

	err := levelup.ImproveSquaddie(suite.improveAllMovement, suite.teros, nil)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.GetMovementDistancePerRound(), Equals, startingMovement + 1)
	checker.Assert(suite.teros.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.CanHitAndRun(), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieCannotDowngradeMovement(checker *C) {
	startingMovement := suite.teros.GetMovementDistancePerRound()
	levelup.ImproveSquaddie(suite.improveAllMovement, suite.teros, nil)

	err := levelup.ImproveSquaddie(suite.upgradeToLightMovement, suite.teros, nil)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.GetMovementDistancePerRound(), Equals, startingMovement + 1)
	checker.Assert(suite.teros.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.CanHitAndRun(), Equals, true)
}

type SquaddieChangePowersWithLevelUpBenefitsSuite struct {
	mageClass *squaddieclass.Class
	teros        *squaddie.Squaddie
	powerRepo    *power.Repository
	gainPower    levelupbenefit.LevelUpBenefit
	upgradePower levelupbenefit.LevelUpBenefit
	spear        *power.Power
	spearLevel2  *power.Power
}

var _ = Suite(&SquaddieChangePowersWithLevelUpBenefitsSuite{})

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) SetUpTest(checker *C) {
	suite.mageClass = &squaddieclass.Class{
		ID:   "ffffffff",
		Name: "Mage",
	}
	suite.teros = squaddie.NewSquaddie("suite.teros")
	suite.teros.MaxHitPoints = 5
	suite.teros.Aim = 0
	suite.teros.Strength = 1
	suite.teros.Mind = 2
	suite.teros.Dodge = 3
	suite.teros.Deflect = 4
	suite.teros.MaxBarrier = 6
	suite.teros.Armor = 7
	suite.teros.AddClass(&squaddieclass.Class{
		ID:   suite.mageClass.ID,
		Name: "Mage",
	})
	suite.teros.SetHPToMax()
	suite.teros.SetBarrierToMax()

	suite.powerRepo = power.NewPowerRepository()

	suite.spear = power.NewPower("suite.spear")
	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.ToHitBonus = 1
	suite.spear.ID = "suite.spearlvl1"
	suite.teros.PowerReferences = []*power.Reference{{Name: "suite.spear", ID: "suite.spearlvl1"}}

	suite.spearLevel2 = power.NewPower("suite.spear")
	suite.spearLevel2.PowerType = power.Physical
	suite.spearLevel2.AttackEffect.ToHitBonus = 1
	suite.spearLevel2.ID = "suite.spearlvl2"
	newPowers := []*power.Power{suite.spear, suite.spearLevel2}
	suite.powerRepo.AddSlicePowerSource(newPowers)

	suite.gainPower = levelupbenefit.LevelUpBenefit{
		ID:                 "aaab1234",
		LevelUpBenefitType: levelupbenefit.Big,
		ClassID:            suite.mageClass.ID,
		PowerIDGained:      []*power.Reference{{Name: "suite.spear", ID: suite.spear.ID}},
	}

	suite.upgradePower = levelupbenefit.LevelUpBenefit{
		ID:                 "aaaa1235",
		LevelUpBenefitType: levelupbenefit.Big,
		ClassID:            suite.mageClass.ID,
		PowerIDLost:        []*power.Reference{{Name: "suite.spear", ID: suite.spear.ID}},
		PowerIDGained:      []*power.Reference{{Name: "suite.spear", ID: suite.spearLevel2.ID}},
	}
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieGainPowers(checker *C) {
	err := levelup.ImproveSquaddie(&suite.gainPower, suite.teros, suite.powerRepo)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "suite.spear")
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.spear.ID)
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieLosePowers(checker *C) {
	levelup.ImproveSquaddie(&suite.gainPower, suite.teros, suite.powerRepo)
	suite.teros.GetInnatePowerIDNames()

	err := levelup.ImproveSquaddie(&suite.upgradePower, suite.teros, suite.powerRepo)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "suite.spear")
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.spearLevel2.ID)
}
