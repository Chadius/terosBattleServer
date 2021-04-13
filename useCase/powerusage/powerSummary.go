package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetPowerSummary returns a summary.
func GetPowerSummary(power *power.Power, user *squaddie.Squaddie, targetSquaddies []*squaddie.Squaddie) *powerusagecontext.PowerForecast {
	summary := &powerusagecontext.PowerForecast{
		UserSquaddieID: user.ID,
		PowerID: power.ID,
		AttackEffectSummary: []*powerusagecontext.AttackingPowerForecast{},
	}

	for _, target := range targetSquaddies {
		summary.AttackEffectSummary = append(summary.AttackEffectSummary, GetExpectedDamage(nil, &powerusagecontext.AttackContext{
			Power:           power,
			Attacker:        user,
			Target:          target,
			IsCounterAttack: false,
		}))
	}
	return summary
}
