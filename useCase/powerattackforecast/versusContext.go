package powerattackforecast

import "github.com/cserrant/terosBattleServer/entity/power"

// VersusContext stores the results of an AttackerContext and DefenderContext.
type VersusContext struct {
	ToHitBonus	int
	DamageDealt int
}

func (context *VersusContext) getPower(setup ForecastSetup) *power.Power {
	return setup.PowerRepo.GetPowerByID(setup.PowerID)
}

func (context *VersusContext)calculate(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.ToHitBonus = context.calculateToHitBonus(attackerContext, defenderContext)
	context.DamageDealt = context.calculateDamageBreakdown(attackerContext, defenderContext)
}

func (context *VersusContext) calculateToHitBonus(attackerContext AttackerContext, defenderContext DefenderContext) int {
	return attackerContext.TotalToHitBonus - defenderContext.TotalToHitPenalty
}

func (context *VersusContext) calculateDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) int {

	damageDealt := attackerContext.RawDamage - defenderContext.BarrierResistance
	if attackerContext.DamageType == power.Physical {
		damageDealt -= defenderContext.ArmorResistance
	}
	if damageDealt < 0 {
		damageDealt = 0
	}

	return damageDealt
}