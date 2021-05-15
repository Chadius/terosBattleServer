package powercommit

import (
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
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
	Damage *powerattackforecast.DamageDistribution
}

// Commit tries to use the power and records the effects.
func (result *Result) Commit() {
	for _, forecast := range result.Forecast.ForecastedResultPerTarget {
		resultPerTarget := &ResultPerTarget{
			UserID: forecast.Setup.UserID,
			TargetID: forecast.Setup.Targets[0],
			PowerID: forecast.Setup.PowerID,
			Attack: &AttackResult{},
		}

		toHitChance := forecast.Attack.VersusContext.ToHitBonus
		attackRoll, defendRoll := result.DieRoller.RollTwoDice()
		resultPerTarget.Attack.HitTarget = attackRoll + toHitChance >= defendRoll

		if !resultPerTarget.Attack.HitTarget {
			resultPerTarget.Attack.Damage = &powerattackforecast.DamageDistribution{
				DamageAbsorbedByArmor:   0,
				DamageAbsorbedByBarrier: 0,
				DamageDealt:             0,
				ExtraBarrierBurnt:       0,
				TotalBarrierBurnt:       0,
			}
		}

		result.ResultPerTarget = append(result.ResultPerTarget, resultPerTarget)
	}
}