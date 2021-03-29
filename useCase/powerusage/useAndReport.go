package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/report"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/utility"
)

// UsePowerAgainstSquaddiesAndGetResults will make the actingSquaddie use the powerUsed against all targetSquaddies.
//   Returns a report indicating what happened to each target.
func UsePowerAgainstSquaddiesAndGetResults(powerUsed *power.Power, actingSquaddie *squaddie.Squaddie, targetSquaddies []*squaddie.Squaddie, d6generator utility.SixSideGenerator) *report.PowerReport {
	powerResults := &report.PowerReport{
		AttackerID:            actingSquaddie.ID,
		PowerID:               powerUsed.ID,
		AttackingPowerResults: []*report.AttackingPowerReport{},
	}

	for _, targetSquaddie := range targetSquaddies {
		attackingResult := GetAttackEffectResults(powerUsed, actingSquaddie, targetSquaddie, d6generator)
		powerResults.AttackingPowerResults = append(powerResults.AttackingPowerResults, attackingResult)
	}
	return powerResults
}

// GetAttackEffectResults looks at the actingSquaddie's powerUsed's AttackingEffect to figure out what happened to the targetSquaddie.
func GetAttackEffectResults(powerUsed *power.Power, actingSquaddie *squaddie.Squaddie, targetSquaddie *squaddie.Squaddie, d6generator utility.SixSideGenerator) *report.AttackingPowerReport {
	attackSummary := GetExpectedDamage(&AttackContext{
		Power:           powerUsed,
		Attacker:        actingSquaddie,
		Target:          targetSquaddie,
		IsCounterAttack: false,
	})

	didItHit := DetermineIfItHit(attackSummary, d6generator)
	if !didItHit {
		return &report.AttackingPowerReport{
			TargetID:        targetSquaddie.ID,
			DamageTaken:     0,
			BarrierDamage:   0,
			WasAHit:         false,
			WasACriticalHit: false,
		}
	}

	didItCrit := DetermineIfItWasACriticalHit(attackSummary, d6generator)
	if !didItCrit {
		return &report.AttackingPowerReport{
			TargetID:        targetSquaddie.ID,
			DamageTaken:     attackSummary.DamageTaken,
			BarrierDamage:   attackSummary.BarrierDamageTaken,
			WasAHit:         true,
			WasACriticalHit: false,
		}
	}

	return &report.AttackingPowerReport{
		TargetID:        targetSquaddie.ID,
		DamageTaken:     attackSummary.CriticalDamageTaken,
		BarrierDamage:   attackSummary.CriticalBarrierDamageTaken,
		WasAHit:         true,
		WasACriticalHit: true,
	}
}

// DetermineIfItHit rolls attacks and determines if the attack hit.
func DetermineIfItHit(summary *AttackingPowerSummary, d6generator utility.SixSideGenerator) bool {
	hitRate := summary.HitRate
	attackRoll, defendRoll := d6generator.RollTwoDice()
	return attackRoll + hitRate >= defendRoll
}

// DetermineIfItWasACriticalHit rolls and determines if the attack was a crit.
func DetermineIfItWasACriticalHit(summary *AttackingPowerSummary, d6generator utility.SixSideGenerator) bool {
	criticalHitThreshold := summary.CriticalHitThreshold
	roll1, roll2 := d6generator.RollTwoDice()
	return roll1 + roll2 < criticalHitThreshold
}