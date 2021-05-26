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
		resultPerTarget := &ResultPerTarget{
			UserID:   calculation.Setup.UserID,
			TargetID: calculation.Setup.Targets[0],
			PowerID:  calculation.Setup.PowerID,
			Attack:   &AttackResult{},
		}

		attackingSquaddie := calculation.Setup.SquaddieRepo.GetOriginalSquaddieByID(calculation.Setup.UserID)
		powerequip.SquaddieEquipPower(attackingSquaddie, calculation.Setup.PowerID, calculation.Setup.PowerRepo)

		toHitChance := calculation.Attack.VersusContext.ToHitBonus
		attackRoll, defendRoll := result.DieRoller.RollTwoDice()
		resultPerTarget.Attack.HitTarget = attackRoll + toHitChance >= defendRoll

		if !resultPerTarget.Attack.HitTarget {
			resultPerTarget.Attack.Damage = &damagedistribution.DamageDistribution{
				DamageAbsorbedByArmor:   0,
				DamageAbsorbedByBarrier: 0,
				DamageDealt:             0,
				ExtraBarrierBurnt:       0,
				TotalBarrierBurnt:       0,
			}
		} else {
			resultPerTarget.Attack.Damage = calculation.Attack.VersusContext.NormalDamage
		}

		result.ResultPerTarget = append(result.ResultPerTarget, resultPerTarget)

		targetSquaddie := calculation.Setup.SquaddieRepo.GetOriginalSquaddieByID(resultPerTarget.TargetID)
		targetSquaddie.Defense.TakeDamageDistribution(resultPerTarget.Attack.Damage)
	}
}
