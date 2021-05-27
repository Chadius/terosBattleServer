package powercommit

import (
	"github.com/cserrant/terosBattleServer/entity/damagedistribution"
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
		result.calculateResultForThisTarget(calculation.Setup, calculation.Attack)

		if calculation.CounterAttack != nil {
			result.calculateResultForThisTarget(calculation.CounterAttackSetup, calculation.CounterAttack)
		}
	}
}

func (result *Result) calculateResultForThisTarget(setup *powerattackforecast.ForecastSetup, attack	*powerattackforecast.AttackForecast) *ResultPerTarget {
	results := &ResultPerTarget{
		UserID:   setup.UserID,
		TargetID: setup.Targets[0],
		PowerID:  setup.PowerID,
		Attack:   &AttackResult{},
	}

	attackingSquaddie := setup.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	powerequip.SquaddieEquipPower(attackingSquaddie, setup.PowerID, setup.PowerRepo)

	toHitChance := attack.VersusContext.ToHitBonus
	attackRoll, defendRoll := result.DieRoller.RollTwoDice()
	results.Attack.HitTarget = attackRoll + toHitChance >= defendRoll

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

	result.ResultPerTarget = append(result.ResultPerTarget, results)

	targetSquaddie := setup.SquaddieRepo.GetOriginalSquaddieByID(results.TargetID)
	targetSquaddie.Defense.TakeDamageDistribution(results.Attack.Damage)

	return results
}