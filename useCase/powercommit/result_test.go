package powercommit_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
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
	fireball *power.Power
	axe      *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastSpearOnMysticMage *powerattackforecast.Forecast

	resultSpearOnBandit *powercommit.Result
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
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.axe, suite.fireball})

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

	suite.resultSpearOnBandit = &powercommit.Result{
		Forecast: suite.forecastSpearOnBandit,
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