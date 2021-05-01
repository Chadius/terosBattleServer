package powerattackforecast

import "github.com/cserrant/terosBattleServer/entity/power"

// VersusContext stores the results of an AttackerContext and DefenderContext.
type VersusContext struct {
	ToHitBonus int

	DamageAbsorbedByArmor   int
	DamageAbsorbedByBarrier int
	DamageDealt             int

	ExtraBarrierBurnt int
	TotalBarrierBurnt int
}

func (context *VersusContext) getPower(setup ForecastSetup) *power.Power {
	return setup.PowerRepo.GetPowerByID(setup.PowerID)
}

func (context *VersusContext)calculate(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.ToHitBonus = context.calculateToHitBonus(attackerContext, defenderContext)
	context.setDamageBreakdown(attackerContext, defenderContext)
}

func (context *VersusContext) calculateToHitBonus(attackerContext AttackerContext, defenderContext DefenderContext) int {
	return attackerContext.TotalToHitBonus - defenderContext.TotalToHitPenalty
}

func (context *VersusContext) setDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	damageDealtToTarget := attackerContext.RawDamage

	context.setBarrierBurntAndDamageAbsorbed(attackerContext, defenderContext, damageDealtToTarget)
	damageDealtToTarget -= context.DamageAbsorbedByBarrier
	context.TotalBarrierBurnt = context.DamageAbsorbedByBarrier + context.ExtraBarrierBurnt

	context.DamageAbsorbedByArmor = context.calculateDamageAbsorbedByArmor(attackerContext, defenderContext, damageDealtToTarget)
	damageDealtToTarget -= context.DamageAbsorbedByArmor

	context.DamageDealt = damageDealtToTarget
}

func (context *VersusContext) calculateDamageAbsorbedByArmor(attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) int {
	if attackerContext.DamageType != power.Physical {
		return 0
	}

	armorAbsorbsAllDamage := damageDealtToTarget <= defenderContext.ArmorResistance
	if armorAbsorbsAllDamage {
		return damageDealtToTarget
	} else {
		return defenderContext.ArmorResistance
	}
}

func (context *VersusContext) setBarrierBurntAndDamageAbsorbed(attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) {
	barrierAbsorbsAllDamageAndExtraBurn := damageDealtToTarget + attackerContext.ExtraBarrierBurn <= defenderContext.BarrierResistance
	if barrierAbsorbsAllDamageAndExtraBurn {
		context.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
		context.DamageAbsorbedByBarrier = damageDealtToTarget
		context.TotalBarrierBurnt = context.DamageAbsorbedByBarrier + context.ExtraBarrierBurnt
		return
	}

	barrierAbsorbsExtraBarrierBurn := attackerContext.ExtraBarrierBurn <= defenderContext.BarrierResistance
	if !barrierAbsorbsExtraBarrierBurn {
		context.ExtraBarrierBurnt = defenderContext.BarrierResistance
		context.DamageAbsorbedByBarrier = 0
		context.TotalBarrierBurnt = context.ExtraBarrierBurnt
		return
	}

	barrierRemainingAfterExtraBarrierBurn := defenderContext.BarrierResistance - attackerContext.ExtraBarrierBurn

	remainingBarrierAbsorbsDamage := damageDealtToTarget <= barrierRemainingAfterExtraBarrierBurn
	if remainingBarrierAbsorbsDamage {
		context.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
		context.DamageAbsorbedByBarrier = damageDealtToTarget
		context.TotalBarrierBurnt = context.DamageAbsorbedByBarrier + context.ExtraBarrierBurnt
		return
	}

	context.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
	context.DamageAbsorbedByBarrier = barrierRemainingAfterExtraBarrierBurn
	context.TotalBarrierBurnt = defenderContext.BarrierResistance
	return
}