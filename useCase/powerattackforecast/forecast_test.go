package powerattackforecast_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
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
	repos *repositories.RepositoryCollection

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
	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus: 1,
		DamageBonus: 1,
		CanBeEquipped: true,
		CanCounterAttack: true,
	}

	suite.axe = power.NewPower("axe")
	suite.axe.PowerType = power.Physical
	suite.axe.AttackEffect = &power.AttackingEffect{
		ToHitBonus: 1,
		DamageBonus: 1,
		CanBeEquipped: true,
		CanCounterAttack: true,
		CounterAttackPenaltyReduction: 0,
	}

	suite.fireball = power.NewPower("fireball")
	suite.fireball.PowerType = power.Spell
	suite.fireball.AttackEffect = &power.AttackingEffect{
		DamageBonus: 3,
		CanBeEquipped: true,
	}

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.axe, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
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
		Repositories: &repositories.RepositoryCollection{
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

	mysticMageEquipsFireball := powerequip.SquaddieEquipPower(suite.mysticMage, suite.fireball.ID, suite.repos)
	checker.Assert(mysticMageEquipsFireball, Equals, true)

	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget[0].CounterAttack, IsNil)
}

func (suite *CounterAttackCalculate) TestCounterAttackHappensIfPossible(checker *C) {
	suite.axe.AttackEffect.CanCounterAttack = true
	powerAddedErrors := suite.bandit.PowerCollection.AddInnatePower(suite.axe)
	checker.Assert(powerAddedErrors, IsNil)

	banditEquipsAxe := powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.repos)
	checker.Assert(banditEquipsAxe, Equals, true)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].CounterAttack.VersusContext.ToHit.ToHitBonus, Equals, -1)
}

type HealingEffectForecast struct {
	lini *squaddie.Squaddie
	teros *squaddie.Squaddie

	healingStaff *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository
	repos *repositories.RepositoryCollection

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
}

var _ = Suite(&HealingEffectForecast{})

func (suite *HealingEffectForecast) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("Teros")
	suite.teros.Identification.ID = "squaddie_teros"
	suite.teros.Identification.Name = "Teros"

	suite.lini = squaddie.NewSquaddie("Lini")
	suite.lini.Identification.ID = "squaddie_lini"
	suite.lini.Identification.Name = "Lini"

	suite.healingStaff = power.NewPower("healing_staff")
	suite.healingStaff.PowerType = power.Spell
	suite.healingStaff.HealingEffect = &power.HealingEffect{
		HitPointsHealed: 3,
	}

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.lini})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.Identification.ID,
			PowerID:         suite.healingStaff.ID,
			Targets:         []string{suite.teros.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}
}

func (suite *HealingEffectForecast) TestForecastedHealingUsesHealingEffect(checker *C) {
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast, NotNil)
	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed)
}

func (suite *HealingEffectForecast) TestForecastedHealingAppliesMindStat(checker *C) {
	suite.lini.Offense.Mind = 3
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed + suite.lini.Offense.Mind)
}

func (suite *HealingEffectForecast) TestForecastedHealingCanBeHalved(checker *C) {
	suite.lini.Offense.Mind = 3
	suite.healingStaff.HealingEffect.HealingAdjustmentBasedOnUserMind = power.Half
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed + (suite.lini.Offense.Mind) / 2)
}

func (suite *HealingEffectForecast) TestForecastedHealingCanBeZeroed(checker *C) {
	suite.lini.Offense.Mind = 3
	suite.healingStaff.HealingEffect.HealingAdjustmentBasedOnUserMind = power.Zero
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed)
}
