package powerusage_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
	. "gopkg.in/check.v1"
)

type SquaddieEquipPowersFromRepo struct {
	teros *squaddie.Squaddie
	spear *power.Power
	scimitar *power.Power
	blot *power.Power
	powerRepo *power.Repository
}

var _ = Suite(&SquaddieEquipPowersFromRepo{})

func (suite *SquaddieEquipPowersFromRepo) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("Teros")
	suite.spear = power.NewPower("Spear")
	suite.spear.AttackEffect.CanBeEquipped = true
	suite.spear.AttackEffect.CanCounterAttack = true

	suite.scimitar = power.NewPower("suite.scimitar the second")
	suite.scimitar.AttackEffect.CanBeEquipped = true

	suite.blot = power.NewPower("Magic spell Blot")

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.spear,
		suite.scimitar,
		suite.blot,
	})
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieEquipsFirstPowerByDefault (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieSkipsUnequippablePowersWhenDefaultEquipping (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.blot.GetReference(),
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieWillNotEquipIfNoEquippablePowers (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo), IsNil)
}

func (suite *SquaddieEquipPowersFromRepo) TestCanChangeEquippedPower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	success := powerusage.SquaddieEquipPower(suite.teros, suite.scimitar.ID, suite.powerRepo)
	checker.Assert(success, Equals, true)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.scimitar.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipUnequibbablePower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	success := powerusage.SquaddieEquipPower(suite.teros, suite.blot.ID, suite.powerRepo)
	checker.Assert(success, Equals, false)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipNonexistentPowers (checker *C) {
	success := powerusage.SquaddieEquipPower(suite.teros, "kwyjibo", suite.powerRepo)
	checker.Assert(success, Equals, false)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo), IsNil)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipUnownedPower (checker *C) {
	notTerosPower := power.NewPower("Does not belong to Teros")
	notTerosPower.AttackEffect.CanBeEquipped = true
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		notTerosPower,
	})

	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	success := powerusage.SquaddieEquipPower(suite.teros, notTerosPower.ID, suite.powerRepo)
	checker.Assert(success, Equals, false)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCanCounter (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	checker.Assert(powerusage.CanSquaddieCounterWithEquippedWeapon(suite.teros, suite.powerRepo), Equals, true)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCannotCounterWithUncounterablePower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	powerusage.SquaddieEquipPower(suite.teros, suite.scimitar.ID, suite.powerRepo)
	checker.Assert(powerusage.CanSquaddieCounterWithEquippedWeapon(suite.teros, suite.powerRepo), Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCannotCounterWithUnequippablePower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)
	checker.Assert(powerusage.CanSquaddieCounterWithEquippedWeapon(suite.teros, suite.powerRepo), Equals, false)
}
