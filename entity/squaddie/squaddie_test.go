package squaddie_test

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieStatsSuite struct{
	teros *squaddie.Squaddie
	mageClass *squaddieclass.Class
	mushroomClass *squaddieclass.Class
}

var _ = Suite(&SquaddieStatsSuite{})

func (suite *SquaddieStatsSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.mageClass = &squaddieclass.Class{ID: "1", Name: "Mage"}
	suite.mushroomClass = &squaddieclass.Class{ID: "2", Name: "Mushroom"}
}

func (suite *SquaddieStatsSuite) TestNameIsSet(checker *C) {
	checker.Assert(suite.teros.Name, Equals, "teros")
}

func (suite *SquaddieStatsSuite) TestGetARandomIDUponCreation(checker *C) {
	checker.Assert(suite.teros.ID, NotNil)
	checker.Assert(suite.teros.ID, Not(Equals), "")
}

func (suite *SquaddieStatsSuite) TestGetANewID(checker *C) {
	initialID := suite.teros.ID
	suite.teros.SetNewIDToRandom()
	checker.Assert(suite.teros.ID, Not(Equals), initialID)
}

func (suite *SquaddieStatsSuite) TestSetMaxHPAndMatchToCurrentHP(checker *C) {
	maxHP := suite.teros.MaxHitPoints
	suite.teros.SetHPToMax()
	checker.Assert(suite.teros.CurrentHitPoints, Equals, maxHP)
}

func (suite *SquaddieStatsSuite) TestCanSetCurrentBarrierToMax(checker *C) {
	suite.teros.MaxBarrier = 2
	suite.teros.SetBarrierToMax()
	checker.Assert(suite.teros.CurrentBarrier, Equals, 2)
}

func (suite *SquaddieStatsSuite) TestDefaultHitPoints(checker *C) {
	checker.Assert(suite.teros.MaxHitPoints, Equals, 5)
	checker.Assert(suite.teros.CurrentHitPoints, Equals, 5)
}

func (suite *SquaddieStatsSuite) TestDefaultMovement(checker *C) {
	checker.Assert(suite.teros.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(suite.teros.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
}

func (suite *SquaddieStatsSuite) TestRaisesErrorIfSquaddieHasUnknownAffiliation(checker *C) {
	newSquaddie := squaddie.NewSquaddie("teros")
	newSquaddie.Affiliation = "Unknown Affiliation"
	err := squaddie.CheckSquaddieForErrors(newSquaddie)
	checker.Assert(err, NotNil)
	checker.Assert(err, ErrorMatches,"Squaddie has unknown affiliation: 'Unknown Affiliation'")
}

func (suite *SquaddieStatsSuite) TestDefenseAgainstPhysicalAttacks(checker *C) {
	suite.teros.Armor = 2
	suite.teros.Dodge = 3
	suite.teros.Deflect = 4
	suite.teros.MaxBarrier = 1
	suite.teros.SetBarrierToMax()
	evasion, barrierDamageReduction, armorDamageReduction := suite.teros.GetDefensiveStatsAgainstPhysical()
	checker.Assert(evasion, Equals, 3)
	checker.Assert(barrierDamageReduction, Equals, 1)
	checker.Assert(armorDamageReduction, Equals, 2)
}

func (suite *SquaddieStatsSuite) TestDefenseAgainstSpellAttacks(checker *C) {
	suite.teros.Armor = 2
	suite.teros.Dodge = 3
	suite.teros.Deflect = 4
	suite.teros.MaxBarrier = 1
	suite.teros.SetBarrierToMax()
	evasion, barrierDamageReduction, armorDamageReduction := suite.teros.GetDefensiveStatsAgainstSpell()
	checker.Assert(evasion, Equals, 4)
	checker.Assert(barrierDamageReduction, Equals, 1)
	checker.Assert(armorDamageReduction, Equals, 0)
}

func (suite *SquaddieStatsSuite) TestOffenseWithPhysicalAttacks(checker *C) {
	suite.teros.Aim = 2
	suite.teros.Strength = 3
	suite.teros.Mind = 4
	toHitBonus, damageBonus := suite.teros.GetOffensiveStatsWithPhysical()
	checker.Assert(toHitBonus, Equals, 2)
	checker.Assert(damageBonus, Equals, 3)
}

func (suite *SquaddieStatsSuite) TestOffenseWithSpellAttacks(checker *C) {
	suite.teros.Aim = 2
	suite.teros.Strength = 3
	suite.teros.Mind = 4
	toHitBonus, damageBonus := suite.teros.GetOffensiveStatsWithSpell()
	checker.Assert(toHitBonus, Equals, 2)
	checker.Assert(damageBonus, Equals, 4)
}

func (suite *SquaddieStatsSuite) TestGainInnatePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.AddInnatePower(attackA)

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].ID, Equals, attackA.ID)
}

func (suite *SquaddieStatsSuite) TestClearInnatePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.AddInnatePower(attackA)
	suite.teros.ClearInnatePowers()

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, DeepEquals, []*power.Reference{})
}

func (suite *SquaddieStatsSuite) TestClearPowerReferences(checker *C) {
	suite.teros.PowerReferences = []*power.Reference{{Name: "Pow pow", ID: "Power Wheels"}}
	suite.teros.ClearTemporaryPowerReferences()
	checker.Assert(suite.teros.PowerReferences, DeepEquals, []*power.Reference{})
}

func (suite *SquaddieStatsSuite) TestRemoveInnatePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.AddInnatePower(attackA)
	suite.teros.RemovePowerByID(attackA.ID)

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, DeepEquals, []*power.Reference{})
}

func (suite *SquaddieStatsSuite) TestRaiseErrorIfTryToRegainSamePower(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	err := suite.teros.AddInnatePower(attackA)
	checker.Assert(err, IsNil)
	err = suite.teros.AddInnatePower(attackA)
	expectedErrorMessage := fmt.Sprintf(`squaddie "teros" already has innate power with ID "%s"`, attackA.ID)
	checker.Assert(err, ErrorMatches, expectedErrorMessage)

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].ID, Equals, attackA.ID)
}

func (suite *SquaddieStatsSuite) TestNewSquaddieHasNoClassesOrLevels(checker *C) {
	checker.Assert(suite.teros.CurrentClass, Equals, "")
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{})
}

func (suite *SquaddieStatsSuite) TestAddClassToSquaddie(checker *C) {
	suite.teros.AddClass(suite.mageClass)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID: 0})
}

func (suite *SquaddieStatsSuite) TestCanTellIfSquaddieAddedClass(checker *C) {
	suite.teros.AddClass(suite.mageClass)
	checker.Assert(suite.teros.HasAddedClass(suite.mageClass.ID), Equals, true)
	checker.Assert(suite.teros.HasAddedClass(suite.mushroomClass.ID), Equals, false)
}

func (suite *SquaddieStatsSuite) TestChangeCurrentClass(checker *C) {
	suite.teros.AddClass(suite.mageClass)
	checker.Assert(suite.teros.CurrentClass, Equals, "")
	err := suite.teros.SetClass(suite.mageClass.ID)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.CurrentClass, Equals, suite.mageClass.ID)
}

func (suite *SquaddieStatsSuite) TestCanSetBaseClass(checker *C) {
	suite.teros.AddClass(suite.mageClass)
	checker.Assert(suite.teros.BaseClassID, Equals, "")
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	checker.Assert(suite.teros.BaseClassID, Equals, suite.mageClass.ID)
}

func (suite *SquaddieStatsSuite) TestRaiseErrorIfClassDoesNotExist(checker *C) {
	suite.teros.AddClass(suite.mageClass)
	checker.Assert(suite.teros.CurrentClass, Equals, "")
	err := suite.teros.SetClass(suite.mushroomClass.ID)
	checker.Assert(err.Error(), Equals, `cannot switch "teros" to unknown class "2"`)
}
