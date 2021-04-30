package powerattackforecast_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	. "gopkg.in/check.v1"
)

type VersusContext struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	spear			*power.Power
	blot			*power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastBlotOnBandit *powerattackforecast.Forecast
}

var _ = Suite(&VersusContext{})


func (suite *VersusContext) SetUpTest(checker *C) {
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
	suite.bandit.Defense.Dodge = 1
	suite.bandit.Defense.Deflect = 2
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.MaxBarrier = 3
	suite.bandit.Defense.SetBarrierToMax()

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

func (suite *VersusContext) TestNetToHitReliesOnToHitMinusDodgeOrDeflect(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.ToHitBonus, Equals, 2)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.ToHitBonus, Equals, 0)
}

func (suite *VersusContext) TestTargetTakesFullDamageAgainstPhysicalWhenNoArmor(checker *C) {
	suite.bandit.Defense.Armor = 0
	suite.bandit.Defense.CurrentBarrier = 0

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.DamageDealt, Equals, 3)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.DamageDealt, Equals, 5)
}

func (suite *VersusContext) TestTargetUsesArmorResistAgainstPhysicalOnly(checker *C) {
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.CurrentBarrier = 0

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.DamageDealt, Equals, 2)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.DamageDealt, Equals, 5)
}

func (suite *VersusContext) TestTargetUsesBarrierToResistDamageFromAllAttacks(checker *C) {
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.CurrentBarrier = 3

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.DamageDealt, Equals, 0)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.DamageDealt, Equals, 2)
}
