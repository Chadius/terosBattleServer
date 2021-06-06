package powercommit

import (
	"github.com/cserrant/terosBattleServer/entity/damagedistribution"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/utility"
)

// Result applies the Forecast given to determine what actually happened.
//  changes are committed.
type Result struct {
	Forecast *powerattackforecast.Forecast
	DieRoller utility.SixSideGenerator
	ResultPerTarget []*ResultPerTarget
}

// ResultPerTarget shows what happened to each target.
type ResultPerTarget struct {
	UserID string
	PowerID string
	TargetID string
	Attack *AttackResult
	AttackRoll, DefendRoll int
}

// AttackResult shows what happens when the power was an attack.
type AttackResult struct {
	HitTarget           bool
	CriticallyHitTarget bool
	Damage *damagedistribution.DamageDistribution
}

// Commit tries to use the power and records the effects.
func (result *Result) Commit() {
	for _, calculation := range result.Forecast.ForecastedResultPerTarget {
		resultForTarget := result.calculateResultForThisTarget(calculation.Setup, calculation.Attack, result.Forecast.Repositories)
		result.ResultPerTarget = append(result.ResultPerTarget, resultForTarget)
	}
	for _, calculation := range result.Forecast.ForecastedResultPerTarget {
		if result.isCounterAttackPossible(calculation) {
			counterAttackResultForTarget := result.calculateResultForThisTarget(calculation.CounterAttackSetup, calculation.CounterAttack, result.Forecast.Repositories)
			result.ResultPerTarget = append(result.ResultPerTarget, counterAttackResultForTarget)
		}
	}
}

func (result *Result) isCounterAttackPossible(calculation powerattackforecast.Calculation) bool {
	if calculation.CounterAttack == nil {
		return false
	}

	counterattacker := calculation.Repositories.SquaddieRepo.GetOriginalSquaddieByID(calculation.CounterAttackSetup.UserID)
	if counterattacker.Defense.IsDead() {
		return false
	}

	return true
}

func (result *Result) calculateResultForThisTarget(setup *powerusagescenario.Setup, attack *powerattackforecast.AttackForecast, repositories *powerusagescenario.RepositoryCollection) *ResultPerTarget {
	results := &ResultPerTarget{
		UserID:   setup.UserID,
		TargetID: setup.Targets[0],
		PowerID:  setup.PowerID,
		Attack:   &AttackResult{},
	}

	attackingSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	powerequip.SquaddieEquipPower(attackingSquaddie, setup.PowerID, repositories.PowerRepo)

	toHitChance := attack.VersusContext.ToHitBonus
	attackRoll, defendRoll := result.DieRoller.RollTwoDice()
	results.Attack.HitTarget = attackRoll + toHitChance >= defendRoll
	results.AttackRoll = attackRoll
	results.DefendRoll = defendRoll

	if results.Attack.HitTarget {
		roll1, roll2 := result.DieRoller.RollTwoDice()
		results.Attack.CriticallyHitTarget = roll1 + roll2 < attack.VersusContext.CriticalHitThreshold
	}

	if !results.Attack.HitTarget {
		results.Attack.Damage = &damagedistribution.DamageDistribution{
			DamageAbsorbedByArmor:   0,
			DamageAbsorbedByBarrier: 0,
			DamageDealt:             0,
			ExtraBarrierBurnt:       0,
			TotalBarrierBurnt:       0,
		}
	} else if results.Attack.CriticallyHitTarget {
		results.Attack.Damage = attack.VersusContext.CriticalHitDamage
	} else {
		results.Attack.Damage = attack.VersusContext.NormalDamage
	}

	targetSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(results.TargetID)
	targetSquaddie.Defense.TakeDamageDistribution(results.Attack.Damage)

	return results
}