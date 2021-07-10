package powercommit

import (
	"github.com/cserrant/terosBattleServer/entity/damagedistribution"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/usecase/squaddiestats"
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
	Healing *HealResult
}

// AttackResult shows what happens when the power was an attack.
type AttackResult struct {
	AttackRoll, DefendRoll int
	AttackerToHitBonus, DefenderToHitPenalty int
	AttackerTotal, DefenderTotal int
	HitTarget           bool
	CriticallyHitTarget bool
	Damage *damagedistribution.DamageDistribution
}

// HealResult shows the effects of recovery abilities.
type HealResult struct {
	HitPointsRestored int
}

// Commit tries to use the power and records the effects.
func (result *Result) Commit() {
	for _, calculation := range result.Forecast.ForecastedResultPerTarget {
		attackResultForTarget := result.getAttackResult(&calculation)
		if attackResultForTarget != nil {
			result.ResultPerTarget = append(result.ResultPerTarget, attackResultForTarget)
		}

		healResultForTarget := result.getHealingResult(&calculation)
		if healResultForTarget != nil {
			result.ResultPerTarget = append(result.ResultPerTarget, healResultForTarget)
		}
	}
	for _, calculation := range result.Forecast.ForecastedResultPerTarget {
		if result.isCounterAttackPossible(calculation) {
			counterAttackResultForTarget := result.calculateAttackResultForThisTarget(calculation.CounterAttackSetup, calculation.CounterAttack, result.Forecast.Repositories)
			result.ResultPerTarget = append(result.ResultPerTarget, counterAttackResultForTarget)
		}
	}
}

func (result *Result) getAttackResult(calculation *powerattackforecast.Calculation) *ResultPerTarget {
	if calculation.Attack == nil {
		return nil
	}
	return result.calculateAttackResultForThisTarget(calculation.Setup, calculation.Attack, result.Forecast.Repositories)
}

func (result *Result) getHealingResult(calculation *powerattackforecast.Calculation) *ResultPerTarget {
	if calculation.HealingForecast == nil {
		return nil
	}
	return result.calculateHealingResultForThisTarget(calculation.Setup, calculation.HealingForecast, result.Forecast.Repositories)
}

func (result *Result) isCounterAttackPossible(calculation powerattackforecast.Calculation) bool {
	if calculation.CounterAttack == nil {
		return false
	}

	counterAttacker := calculation.Repositories.SquaddieRepo.GetOriginalSquaddieByID(calculation.CounterAttackSetup.UserID)
	if counterAttacker.Defense.IsDead() {
		return false
	}

	return true
}

func (result *Result) calculateAttackResultForThisTarget(setup *powerusagescenario.Setup, attack *powerattackforecast.AttackForecast, repositories *repositories.RepositoryCollection) *ResultPerTarget {
	results := &ResultPerTarget{
		UserID:   setup.UserID,
		TargetID: setup.Targets[0],
		PowerID:  setup.PowerID,
		Attack:   &AttackResult{},
	}

	attackingSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	powerequip.SquaddieEquipPower(attackingSquaddie, setup.PowerID, repositories)

	attackRoll, defendRoll := result.DieRoller.RollTwoDice()
	results.Attack.AttackerToHitBonus = attack.VersusContext.ToHit.AttackerToHitBonus
	results.Attack.DefenderToHitPenalty = attack.VersusContext.ToHit.DefenderToHitPenalty

	results.Attack.AttackRoll = attackRoll
	results.Attack.DefendRoll = defendRoll

	results.Attack.AttackerTotal = results.Attack.AttackRoll + results.Attack.AttackerToHitBonus
	results.Attack.DefenderTotal = results.Attack.DefendRoll + results.Attack.DefenderToHitPenalty

	results.Attack.HitTarget = results.Attack.AttackerTotal >= results.Attack.DefenderTotal
	results.Attack.CriticallyHitTarget = attack.AttackerContext.CanCritical && results.Attack.AttackerTotal >= results.Attack.DefenderTotal + attack.AttackerContext.CriticalHitThreshold

	if !results.Attack.HitTarget {
		results.Attack.Damage = &damagedistribution.DamageDistribution{
			DamageAbsorbedByArmor:   0,
			DamageAbsorbedByBarrier: 0,
			RawDamageDealt:          0,
			ExtraBarrierBurnt:       0,
			TotalRawBarrierBurnt:    0,
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

func (result *Result) calculateHealingResultForThisTarget(setup *powerusagescenario.Setup, forecast *powerattackforecast.HealingForecast, repositories *repositories.RepositoryCollection) *ResultPerTarget {
	resultForThisTarget := &ResultPerTarget{
		UserID:   setup.UserID,
		TargetID: setup.Targets[0],
		PowerID:  setup.PowerID,
		Healing: &HealResult{
			HitPointsRestored: 0,
		},
	}

	healingSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	powerequip.SquaddieEquipPower(healingSquaddie, setup.PowerID, repositories)

	targetSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(resultForThisTarget.TargetID)
	maximumHealing, err := squaddiestats.GetHitPointsHealedWithPower(setup.UserID, setup.PowerID, repositories)
	if err != nil {
		return resultForThisTarget
	}
	hitPointsRestored := targetSquaddie.Defense.GainHitPoints(maximumHealing)
	resultForThisTarget.Healing.HitPointsRestored = hitPointsRestored
	return resultForThisTarget
}