package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetPowerSummary returns a summary.
func GetPowerSummary(context *powerusagecontext.PowerUsageContext, power *power.Power, user *squaddie.Squaddie, targetSquaddies []*squaddie.Squaddie) *powerusagecontext.PowerForecast {
	summary := &powerusagecontext.PowerForecast{
		UserSquaddieID: user.ID,
		PowerID: power.ID,
		AttackEffectSummary: []*powerusagecontext.AttackingPowerForecast{},
	}

	for _, target := range targetSquaddies {
		summary.AttackEffectSummary = append(summary.AttackEffectSummary, GetExpectedDamage(context, &powerusagecontext.AttackContext{
			PowerID:           power.ID,
			AttackerID:        user.ID,
			TargetID:          target.ID,
			IsCounterAttack: false,
		}))
	}
	return summary
}
