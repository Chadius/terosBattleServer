package powerattackforecast

import "github.com/cserrant/terosBattleServer/entity/power"

// VersusContext stores the results of an AttackerContext and DefenderContext.
type VersusContext struct {
	ToHitBonus int

	NormalDamage *DamageDistribution
	CriticalHitDamage *DamageDistribution

	CanCritical bool
	CriticalHitThreshold int
}

// DamageDistribution tracks how damage is distributed.
type DamageDistribution struct {
	DamageAbsorbedByArmor   int
	DamageAbsorbedByBarrier int
	DamageDealt             int
	ExtraBarrierBurnt int
	TotalBarrierBurnt int
}

func (context *VersusContext) calculate(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.ToHitBonus = context.calculateToHitBonus(attackerContext, defenderContext)
	context.setNormalDamageBreakdown(attackerContext, defenderContext)

	context.setCriticalHitChance(attackerContext)
	context.setCriticalDamageBreakdown(attackerContext, defenderContext)
}

func (context *VersusContext) calculateToHitBonus(attackerContext AttackerContext, defenderContext DefenderContext) int {
	return attackerContext.TotalToHitBonus - defenderContext.TotalToHitPenalty
}

func (context *VersusContext) setNormalDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.NormalDamage = context.setDamageBreakdown(attackerContext.RawDamage, attackerContext, defenderContext)
}

func (context *VersusContext) setCriticalDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	if context.CanCritical {
		context.CriticalHitDamage = context.setDamageBreakdown(attackerContext.CriticalHitDamage, attackerContext, defenderContext)
	}
}

func (context *VersusContext) setDamageBreakdown(damageDealtToTarget int, attackerContext AttackerContext, defenderContext DefenderContext) *DamageDistribution {
	distribution := &DamageDistribution{}

	context.setBarrierBurntAndDamageAbsorbed(distribution, attackerContext, defenderContext, damageDealtToTarget)

	damageDealtToTarget -= distribution.DamageAbsorbedByBarrier
	distribution.TotalBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt

	distribution.DamageAbsorbedByArmor = context.calculateDamageAbsorbedByArmor(attackerContext, defenderContext, damageDealtToTarget)
	damageDealtToTarget -= distribution.DamageAbsorbedByArmor

	distribution.DamageDealt = damageDealtToTarget

	return distribution
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

func (context *VersusContext) setBarrierBurntAndDamageAbsorbed(distribution *DamageDistribution, attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) {
	barrierAbsorbsAllDamageAndExtraBurn := damageDealtToTarget + attackerContext.ExtraBarrierBurn <= defenderContext.BarrierResistance
	if barrierAbsorbsAllDamageAndExtraBurn {
		distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
		distribution.DamageAbsorbedByBarrier = damageDealtToTarget
		distribution.TotalBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt
		return
	}

	barrierAbsorbsExtraBarrierBurn := attackerContext.ExtraBarrierBurn <= defenderContext.BarrierResistance
	if !barrierAbsorbsExtraBarrierBurn {
		distribution.ExtraBarrierBurnt = defenderContext.BarrierResistance
		distribution.DamageAbsorbedByBarrier = 0
		distribution.TotalBarrierBurnt = distribution.ExtraBarrierBurnt
		return
	}

	barrierRemainingAfterExtraBarrierBurn := defenderContext.BarrierResistance - attackerContext.ExtraBarrierBurn

	remainingBarrierAbsorbsDamage := damageDealtToTarget <= barrierRemainingAfterExtraBarrierBurn
	if remainingBarrierAbsorbsDamage {
		distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
		distribution.DamageAbsorbedByBarrier = damageDealtToTarget
		distribution.TotalBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt
		return
	}

	distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
	distribution.DamageAbsorbedByBarrier = barrierRemainingAfterExtraBarrierBurn
	distribution.TotalBarrierBurnt = defenderContext.BarrierResistance
}

func (context *VersusContext) setCriticalHitChance(attackerContext AttackerContext) {
	context.CanCritical = attackerContext.CanCritical
	if context.CanCritical {
		context.CriticalHitThreshold = attackerContext.CriticalHitThreshold
	} else {
		context.CriticalHitThreshold = 0
	}
}