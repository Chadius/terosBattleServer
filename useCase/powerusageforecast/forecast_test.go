package powerattackusage_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	powerattackusage "github.com/cserrant/terosBattleServer/usecase/powerusageforecast"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ForecastDamage struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	spear			*power.Power
	blot			*power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackusage.Forecast
	//forecastBlotOnBandit *powerattackusage.Forecast
}

var _ = Suite(&ForecastDamage{})


func (suite *ForecastDamage) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
	suite.teros.Offense.Aim = 2

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.ToHitBonus = 1

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

	suite.forecastSpearOnBandit = &powerattackusage.Forecast{
		Setup: powerattackusage.ForecastSetup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
			IsCounterAttack: false,
		},
	}
	//suite.forecastBlotOnBandit = *powerattackusage.Forecast
}

func (suite *ForecastDamage) TestUserHitBonus(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.TotalToHitBonus, Equals, 3)
}