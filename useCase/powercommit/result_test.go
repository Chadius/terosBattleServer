package powercommit_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type resultOnAttack struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	mysticMage		*squaddie.Squaddie

	spear    *power.Power
	blot    *power.Power
	fireball *power.Power
	axe      *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	resultSpearOnBandit *powercommit.Result

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit *powercommit.Result

	forecastSpearOnMysticMage *powerattackforecast.Forecast

	forecastFireballOnBandit *powerattackforecast.Forecast
	resultFireballOnBandit *powercommit.Result
}

var _ = Suite(&resultOnAttack{})

func (suite *resultOnAttack) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
	suite.teros.Offense.Aim = 2
	suite.teros.Offense.Strength = 2
	suite.teros.Offense.Mind = 2

	suite.mysticMage = squaddie.NewSquaddie("mysticMage")
	suite.mysticMage.Identification.Name = "mysticMage"
	suite.mysticMage.Offense.Mind = 2

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.ToHitBonus = 1
	suite.spear.AttackEffect.DamageBonus = 1
	suite.spear.AttackEffect.CanBeEquipped = true
	suite.spear.AttackEffect.CanCounterAttack = true

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.axe = power.NewPower("axe")
	suite.axe.PowerType = power.Physical
	suite.axe.AttackEffect.ToHitBonus = 1
	suite.axe.AttackEffect.DamageBonus = 1
	suite.axe.AttackEffect.CanBeEquipped = true

	suite.fireball = power.NewPower("fireball")
	suite.fireball.PowerType = power.Spell
	suite.fireball.AttackEffect.DamageBonus = 3

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.axe, suite.fireball})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.powerRepo,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.powerRepo,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.mysticMage,
		[]*power.Reference{
			suite.fireball.GetReference(),
		},
		suite.powerRepo,
	)

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	suite.resultSpearOnBandit = &powercommit.Result{
		Forecast: suite.forecastSpearOnBandit,
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	suite.resultBlotOnBandit = &powercommit.Result{
		Forecast: suite.forecastBlotOnBandit,
	}

	suite.forecastSpearOnMysticMage = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.mysticMage.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}

	suite.forecastFireballOnBandit = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.mysticMage.Identification.ID,
			PowerID:         suite.fireball.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	suite.resultFireballOnBandit = &powercommit.Result{
		Forecast: suite.forecastFireballOnBandit,
	}
}

func (suite *resultOnAttack) TestAttackCanMiss(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget, HasLen, 1)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageDealt, Equals, 0)
}

func (suite *resultOnAttack) TestAttackCanHitButNotCritically(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Mind = 2

	suite.blot.AttackEffect.DamageBonus = 3

	suite.bandit.Defense.CurrentBarrier = 3
	suite.bandit.Defense.Armor = 1

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].PowerID, Equals, suite.blot.ID)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageDealt, Equals, 2)

	checker.Assert(
		suite.bandit.Defense.CurrentHitPoints,
		Equals,
		suite.bandit.Defense.MaxHitPoints - suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageDealt,
	)
}

func (suite *resultOnAttack) TestAttackCanHitCritically(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}
	suite.blot.AttackEffect.CriticalHitThreshold = 900

	suite.teros.Offense.Mind = 2

	suite.blot.AttackEffect.DamageBonus = 3

	suite.bandit.Defense.CurrentBarrier = 3
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.MaxHitPoints = 1
	suite.bandit.Defense.SetHPToMax()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].PowerID, Equals, suite.blot.ID)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageDealt, Equals, 7)

	checker.Assert(
		suite.bandit.Defense.CurrentHitPoints,
		Equals,
		0,
	)
}

func (suite *resultOnAttack) TestCounterAttacks(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Strength = 2
	suite.teros.Defense.Armor = 0
	suite.teros.Defense.CurrentBarrier = 0

	suite.spear.AttackEffect.DamageBonus = 3

	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.DamageBonus = 3
	suite.bandit.Offense.Strength = 0
	suite.bandit.Defense.Armor = 1
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.powerRepo)

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.Identification.ID)

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].PowerID, Equals, suite.axe.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].UserID, Equals, suite.bandit.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].TargetID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageDealt, Equals, 3)

	checker.Assert(
		suite.teros.Defense.CurrentHitPoints,
		Equals,
		suite.teros.Defense.MaxHitPoints - suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageDealt,
	)
}

type EquipPowerWhenCommitting struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	mysticMage		*squaddie.Squaddie

	spear    *power.Power
	blot    *power.Power
	fireball *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	resultSpearOnBandit *powercommit.Result

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit *powercommit.Result

	forecastFireballOnBandit *powerattackforecast.Forecast
	resultFireballOnBandit *powercommit.Result
}

var _ = Suite(&EquipPowerWhenCommitting{})

func (suite *EquipPowerWhenCommitting) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
	suite.teros.Offense.Aim = 2
	suite.teros.Offense.Strength = 2
	suite.teros.Offense.Mind = 2

	suite.mysticMage = squaddie.NewSquaddie("mysticMage")
	suite.mysticMage.Identification.Name = "mysticMage"
	suite.mysticMage.Offense.Mind = 2

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.ToHitBonus = 1
	suite.spear.AttackEffect.DamageBonus = 1
	suite.spear.AttackEffect.CanBeEquipped = true
	suite.spear.AttackEffect.CanCounterAttack = true

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.fireball = power.NewPower("fireball")
	suite.fireball.PowerType = power.Spell
	suite.fireball.AttackEffect.DamageBonus = 3

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.fireball})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.powerRepo,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.mysticMage,
		[]*power.Reference{
			suite.fireball.GetReference(),
		},
		suite.powerRepo,
	)

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	suite.resultSpearOnBandit = &powercommit.Result{
		Forecast: suite.forecastSpearOnBandit,
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	suite.resultBlotOnBandit = &powercommit.Result{
		Forecast: suite.forecastBlotOnBandit,
	}

	suite.forecastFireballOnBandit = &powerattackforecast.Forecast{
		Setup: powerattackforecast.ForecastSetup{
			UserID:          suite.mysticMage.Identification.ID,
			PowerID:         suite.fireball.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	suite.resultFireballOnBandit = &powercommit.Result{
		Forecast: suite.forecastFireballOnBandit,
	}
}

func (suite *EquipPowerWhenCommitting) TestCommitWillTryToEquipPower(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.teros.PowerCollection.CurrentlyEquippedPowerID, Equals, suite.spear.ID)
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillKeepPreviousPowerIfCommitPowerCannotBeEquipped(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.spear.ID, suite.powerRepo)
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(powerequip.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.spear.ID)
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillNotEquipPowerIfNoneExistAfterCommitting(checker *C) {
	suite.resultFireballOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastFireballOnBandit.CalculateForecast()
	suite.resultFireballOnBandit.Commit()
	checker.Assert(powerequip.GetEquippedPower(suite.mysticMage, suite.powerRepo), IsNil)
}