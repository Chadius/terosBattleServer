package squaddiestats_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/squaddiestats"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type squaddieOffense struct {
	teros			*squaddie.Squaddie

	spear    *power.Power
	blot    *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository
}

var _ = Suite(&squaddieOffense{})

func (suite *squaddieOffense) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.powerRepo,
	)
}

func (suite *squaddieOffense) TestSquaddieMeasuresAim(checker *C) {
	suite.teros.Offense.Aim = 1
	//suite.teros.Offense.Strength = 2
	//suite.teros.Offense.Mind = 3

	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                1,
		DamageBonus:               1,
		ExtraBarrierBurn:          0,
		CanBeEquipped:             true,
		CanCounterAttack:          true,
		CounterAttackToHitPenalty: -2,
		CriticalEffect:            &power.CriticalEffect{
			CriticalHitThresholdBonus: 3,
			Damage:                    5,
		},
	}

	suite.blot.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                2,
		DamageBonus:               0,
		ExtraBarrierBurn:          2,
		CanBeEquipped:             true,
		CanCounterAttack:          false,
	}

	spearAim, spearErr := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)

	blotAim, blotErr := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotAim, Equals, 3)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfSquaddieDoesNotExist(checker *C) {
	_, err := squaddiestats.GetSquaddieAimWithPower("does not exist", suite.spear.ID, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(err, ErrorMatches, "squaddie could not be found, ID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotExist(checker *C) {
	_, err := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, "does not exist", suite.squaddieRepo, suite.powerRepo)
	checker.Assert(err, ErrorMatches, "power could not be found, ID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerHasNoOffense(checker *C) {
	wait := power.NewPower("wait")
	wait.PowerType = power.Physical
	wait.ID = "powerWait"

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{wait})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
			wait.GetReference(),
		},
		suite.powerRepo,
	)

	_, err := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, wait.ID, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(err, ErrorMatches, "cannot attack with power, ID: powerWait")
}

func (suite *squaddieOffense) TestGetRawDamageOfPhysicalPower(checker *C) {
	suite.teros.Offense.Strength = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                1,
		DamageBonus:               1,
		ExtraBarrierBurn:          0,
		CanBeEquipped:             true,
		CanCounterAttack:          true,
		CounterAttackToHitPenalty: -2,
		CriticalEffect:            &power.CriticalEffect{
			CriticalHitThresholdBonus: 3,
			Damage:                    5,
		},
	}

	spearDamage, spearErr := squaddiestats.GetSquaddieRawDamageWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDamage, Equals, 2)
}

func (suite *squaddieOffense) TestGetRawDamageOfSpellPower(checker *C) {
	suite.teros.Offense.Mind = 3

	suite.blot.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                2,
		DamageBonus:               0,
		ExtraBarrierBurn:          2,
		CanBeEquipped:             true,
		CanCounterAttack:          false,
	}

	blotDamage, blotErr := squaddiestats.GetSquaddieRawDamageWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 3)
}

// func (suite *squaddieOffense) TestGetCriticalThresholdOfPower(checker *C) {
// func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotCrit(checker *C) {
// func (suite *squaddieOffense) TestGetCriticalDamageOfPower(checker *C) {
// func (suite *squaddieOffense) TestSquaddieCanCounterAttackWithPower(checker *C) {
// func (suite *squaddieOffense) TestSquaddieShowsCounterAttackToHit(checker *C) {
// func (suite *squaddieOffense) TestGetTotalBarrierBurnOfAttacks(checker *C) {