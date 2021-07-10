package squaddiestats_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
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

	repos *repositories.RepositoryCollection
}

var _ = Suite(&squaddieOffense{})

func (suite *squaddieOffense) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.ID = "powerSpear"

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell
	suite.blot.ID = "powerBlot"

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo: suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)
}

func (suite *squaddieOffense) TestSquaddieMeasuresAim(checker *C) {
	suite.teros.Offense.Aim = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                    1,
		DamageBonus:                   1,
		ExtraBarrierBurn:              0,
		CanBeEquipped:                 true,
		CanCounterAttack:              true,
		CounterAttackPenaltyReduction: -2,
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

	spearAim, spearErr := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)

	blotAim, blotErr := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotAim, Equals, 3)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfSquaddieDoesNotExist(checker *C) {
	_, err := squaddiestats.GetSquaddieAimWithPower("does not exist", suite.spear.ID, suite.repos)
	checker.Assert(err, ErrorMatches, "squaddie could not be found, ID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotExist(checker *C) {
	_, err := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, "does not exist", suite.repos)
	checker.Assert(err, ErrorMatches, "power could not be found, ID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerHasNoAttackEffect(checker *C) {
	wait := power.NewPower("wait")
	wait.PowerType = power.Physical
	wait.ID = "powerWait"

	suite.powerRepo.AddSlicePowerSource([]*power.Power{wait})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
			wait.GetReference(),
		},
		suite.repos,
	)

	_, err := squaddiestats.GetSquaddieAimWithPower(suite.teros.Identification.ID, wait.ID, suite.repos)
	checker.Assert(err, ErrorMatches, "cannot attack with power, ID: powerWait")
}

func (suite *squaddieOffense) TestGetRawDamageOfPhysicalPower(checker *C) {
	suite.teros.Offense.Strength = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                    1,
		DamageBonus:                   1,
		ExtraBarrierBurn:              0,
		CanBeEquipped:                 true,
		CanCounterAttack:              true,
		CounterAttackPenaltyReduction: -2,
		CriticalEffect:            &power.CriticalEffect{
			CriticalHitThresholdBonus: 3,
			Damage:                    5,
		},
	}

	spearDamage, spearErr := squaddiestats.GetSquaddieRawDamageWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
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

	blotDamage, blotErr := squaddiestats.GetSquaddieRawDamageWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 3)
}

func (suite *squaddieOffense) TestGetCriticalThresholdOfPower(checker *C) {
	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                    1,
		DamageBonus:                   1,
		ExtraBarrierBurn:              0,
		CanBeEquipped:                 true,
		CanCounterAttack:              true,
		CounterAttackPenaltyReduction: -2,
		CriticalEffect:            &power.CriticalEffect{
			CriticalHitThresholdBonus: 2,
			Damage:                    5,
		},
	}

	spearCritThreat, critErr := squaddiestats.GetSquaddieCriticalThresholdWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(critErr, IsNil)
	checker.Assert(spearCritThreat, Equals, 4)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotCrit(checker *C) {
	suite.blot.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                    1,
		DamageBonus:                   1,
		ExtraBarrierBurn:              0,
		CanBeEquipped:                 true,
		CanCounterAttack:              true,
		CounterAttackPenaltyReduction: -2,
	}

	_, critErr := squaddiestats.GetSquaddieCriticalThresholdWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(critErr, ErrorMatches, "cannot critical hit with power, ID: powerBlot")
}

func (suite *squaddieOffense) TestGetCriticalDamageOfPower(checker *C) {
	suite.teros.Offense.Strength = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                    1,
		DamageBonus:                   1,
		ExtraBarrierBurn:              0,
		CanBeEquipped:                 true,
		CanCounterAttack:              true,
		CounterAttackPenaltyReduction: -2,
		CriticalEffect:            &power.CriticalEffect{
			CriticalHitThresholdBonus: 2,
			Damage:                    5,
		},
	}

	spearDamage, damageErr := squaddiestats.GetSquaddieCriticalRawDamageWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(damageErr, IsNil)
	checker.Assert(spearDamage, Equals, 7)
}

func (suite *squaddieOffense) TestSquaddieCanCounterAttackWithPower(checker *C) {
	suite.spear.AttackEffect = &power.AttackingEffect{
		CanCounterAttack:          true,
	}

	suite.blot.AttackEffect = &power.AttackingEffect{
		CanCounterAttack:          false,
	}

	spearCanCounter, spearErr := squaddiestats.GetSquaddieCanCounterAttackWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearCanCounter, Equals, true)

	blotCanCounter, blotErr := squaddiestats.GetSquaddieCanCounterAttackWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotCanCounter, Equals, false)
}

func (suite *squaddieOffense) TestSquaddieShowsCounterAttackToHit(checker *C) {
	suite.teros.Offense.Aim = 2

	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                    1,
		CanCounterAttack:              true,
		CounterAttackPenaltyReduction: 1,
	}

	spearAim, spearErr := squaddiestats.GetSquaddieCounterAttackAimWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)
}

func (suite *squaddieOffense) TestGetTotalBarrierBurnOfAttacks(checker *C) {
	suite.teros.Offense.Mind = 3

	suite.blot.AttackEffect = &power.AttackingEffect{
		ToHitBonus:                2,
		DamageBonus:               0,
		ExtraBarrierBurn:          2,
		CanBeEquipped:             true,
		CanCounterAttack:          false,
	}

	blotDamage, blotErr := squaddiestats.GetSquaddieExtraBarrierBurnWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 2)
}

func (suite *squaddieOffense) TestCanCriticallyHitWithPower(checker *C) {
	suite.spear.AttackEffect = &power.AttackingEffect{
		CriticalEffect:            &power.CriticalEffect{
			Damage:                    1,
		},
	}

	spearCanCrit, spearCanCritErr := squaddiestats.GetSquaddieCanCriticallyHitWithPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearCanCritErr, IsNil)
	checker.Assert(spearCanCrit, Equals, true)

	suite.blot.AttackEffect = &power.AttackingEffect{}

	blotCanCrit, blotCanCritErr := squaddiestats.GetSquaddieCanCriticallyHitWithPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotCanCritErr, IsNil)
	checker.Assert(blotCanCrit, Equals, false)
}

type healingPower struct {
	lini			*squaddie.Squaddie
	healingStaff *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	repos *repositories.RepositoryCollection
}

var _ = Suite(&healingPower{})

func (suite *healingPower) SetUpTest(checker *C) {
	suite.lini = squaddie.NewSquaddie("lini")
	suite.lini.Identification.Name = "lini"

	suite.healingStaff = power.NewPower("healing_staff")
	suite.healingStaff.PowerType = power.Spell

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.lini})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo: suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.lini,
		[]*power.Reference{
			suite.healingStaff.GetReference(),
		},
		suite.repos,
	)
}

func (suite *healingPower) TestSquaddieKnowsHealingPotential(checker *C) {
	suite.lini.Offense.Mind = 1
	suite.healingStaff.HealingEffect = &power.HealingEffect{
		HitPointsHealed: 3,
		HealingAdjustmentBasedOnUserMind: power.Full,
	}

	staffHeal, staffErr := squaddiestats.GetHitPointsHealedWithPower(suite.lini.Identification.ID, suite.healingStaff.ID, suite.repos)
	checker.Assert(staffErr, IsNil)
	checker.Assert(staffHeal, Equals, 4)
}