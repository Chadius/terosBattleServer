package powerattackforecast_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	. "gopkg.in/check.v1"
)

type AttackContext struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	spear			*power.Power
	blot			*power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastBlotOnBandit *powerattackforecast.Forecast
}

var _ = Suite(&AttackContext{})


func (suite *AttackContext) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
	suite.teros.Offense.Aim = 2
	suite.teros.Offense.Strength = 2
	suite.teros.Offense.Mind = 2

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.ToHitBonus = 1
	suite.spear.AttackEffect.DamageBonus = 1

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 3

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

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
}

func (suite *AttackContext) TestGetAttackerHitBonus(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.TotalToHitBonus, Equals, 3)
}

func (suite *AttackContext) TestGetAttackerPhysicalRawDamage(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.DamageType, Equals, power.Type(power.Physical))
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.RawDamage, Equals, 3)
}

func (suite *AttackContext) TestGetAttackerSpellDamage(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.DamageType, Equals, power.Type(power.Spell))
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.RawDamage, Equals, 5)
}

func (suite *AttackContext) TestCriticalHits(checker *C) {
	suite.spear.AttackEffect.CriticalHitThreshold = 8
	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CanCritical, Equals, true)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CriticalHitThreshold, Equals, 8)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CriticalHitDamage, Equals, 6)
}