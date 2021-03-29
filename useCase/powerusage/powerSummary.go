package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// PowerSummary showcases the expected results of using a given power.
type PowerSummary struct {
	UserSquaddieID 		string
	PowerID        		string
	AttackEffectSummary []*AttackingPowerSummary
}

// GetPowerSummary returns a summary.
func GetPowerSummary(power *power.Power, user *squaddie.Squaddie, targetSquaddies []*squaddie.Squaddie) *PowerSummary {
	summary := &PowerSummary{
		UserSquaddieID: user.ID,
		PowerID: power.ID,
		AttackEffectSummary: []*AttackingPowerSummary{},
	}

	for _, target := range targetSquaddies {
		summary.AttackEffectSummary = append(summary.AttackEffectSummary, GetExpectedDamage(&AttackContext{
			Power:           power,
			Attacker:        user,
			Target:          target,
			IsCounterattack: false,
		}))
	}
	return summary
}
