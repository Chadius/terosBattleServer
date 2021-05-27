package powerattackforecast_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CounterAttackCalculate struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	mysticMage		*squaddie.Squaddie

	spear    *power.Power
	fireball *power.Power
	axe      *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastSpearOnMysticMage *powerattackforecast.Forecast
}

var _ = Suite(&CounterAttackCalculate{})

func (suite *CounterAttackCalculate) SetUpTest(checker *C) {
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

	suite.axe = power.NewPower("axe")
	suite.axe.PowerType = power.Physical
	suite.axe.AttackEffect.ToHitBonus = 1
	suite.axe.AttackEffect.DamageBonus = 1
	suite.axe.AttackEffect.CanBeEquipped = true

	suite.fireball = power.NewPower("fireball")
	suite.fireball.PowerType = power.Spell
	suite.fireball.AttackEffect.DamageBonus = 3
	suite.fireball.AttackEffect.CanBeEquipped = true

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.axe, suite.fireball})

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

	suite.forecastSpearOnMysticMage = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.mysticMage.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &powerusagescenario.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}
}

func (suite *CounterAttackCalculate) TestNoCounterAttackHappensIfNoEquippedPower(checker *C) {
	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget[0].CounterAttack, IsNil)
}

func (suite *CounterAttackCalculate) TestNoCounterAttackHappensIfEquippedPowerCannotCounter(checker *C) {
	powerAddedErrors := suite.mysticMage.PowerCollection.AddInnatePower(suite.fireball)
	checker.Assert(powerAddedErrors, IsNil)

	mysticMageEquipsFireball := powerequip.SquaddieEquipPower(suite.mysticMage, suite.fireball.ID, suite.powerRepo)
	checker.Assert(mysticMageEquipsFireball, Equals, true)

	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget[0].CounterAttack, IsNil)
}

func (suite *CounterAttackCalculate) TestCounterAttackHappensIfPossible(checker *C) {
	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.CounterAttackToHitPenalty = 2
	powerAddedErrors := suite.bandit.PowerCollection.AddInnatePower(suite.axe)
	checker.Assert(powerAddedErrors, IsNil)

	banditEquipsAxe := powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.powerRepo)
	checker.Assert(banditEquipsAxe, Equals, true)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].CounterAttack.VersusContext.ToHitBonus, Equals, -1)
}
