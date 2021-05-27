package powerattackforecast_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	. "gopkg.in/check.v1"
)

type DefenderContext struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	spear			*power.Power
	blot			*power.Power
	axe			*power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastBlotOnBandit *powerattackforecast.Forecast
}

var _ = Suite(&DefenderContext{})


func (suite *DefenderContext) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"

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

	suite.axe = power.NewPower("axe")
	suite.axe.PowerType = power.Physical
	suite.axe.AttackEffect.ToHitBonus = 1
	suite.axe.AttackEffect.DamageBonus = 4
	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.CanBeEquipped = true

	suite.bandit.PowerCollection.AddInnatePower(suite.axe)

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.axe})

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &powerusagescenario.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &powerusagescenario.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}
}

func (suite *DefenderContext) TestGetDefenderDodge(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.TotalToHitPenalty, Equals, 1)
}

func (suite *DefenderContext) TestGetDefenderArmorResistance(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.ArmorResistance, Equals, 1)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.ArmorResistance, Equals, 0)
}

func (suite *DefenderContext) TestGetDefenderDeflect(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.TotalToHitPenalty, Equals, 2)
}

func (suite *DefenderContext) TestGetDefenderBarrierAbsorb(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.BarrierResistance, Equals, 3)

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.BarrierResistance, Equals, 3)
}
