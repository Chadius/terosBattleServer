package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/damagedistribution"
	"github.com/cserrant/terosBattleServer/entity/power"
)

// VersusContext compares an AttackerContext and DefenderContext to determine the possible results.
type VersusContext struct {
	ToHit *damagedistribution.ToHitComparison

	NormalDamage *damagedistribution.DamageDistribution
	CriticalHitDamage *damagedistribution.DamageDistribution

	CanCritical bool
	CriticalHitThreshold int
}

func (context *VersusContext) calculate(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.calculateToHitBonus(attackerContext, defenderContext)
	context.setNormalDamageBreakdown(attackerContext, defenderContext)

	context.setCriticalHitChance(attackerContext)
	context.setCriticalDamageBreakdown(attackerContext, defenderContext)
}

func (context *VersusContext) calculateToHitBonus(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.ToHit = &damagedistribution.ToHitComparison{}
	context.ToHit.AttackerToHitBonus = attackerContext.TotalToHitBonus
	context.ToHit.DefenderToHitPenalty = defenderContext.TotalToHitPenalty
	context.ToHit.ToHitBonus = context.ToHit.AttackerToHitBonus - context.ToHit.DefenderToHitPenalty
}

func (context *VersusContext) setNormalDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.NormalDamage = context.setDamageBreakdown(attackerContext.RawDamage, attackerContext, defenderContext)
}

func (context *VersusContext) setCriticalDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	if context.CanCritical {
		context.CriticalHitDamage = context.setDamageBreakdown(attackerContext.CriticalHitDamage, attackerContext, defenderContext)
	}
}

func (context *VersusContext) setDamageBreakdown(damageDealtToTarget int, attackerContext AttackerContext, defenderContext DefenderContext) *damagedistribution.DamageDistribution {
	distribution := &damagedistribution.DamageDistribution{}

	context.setBarrierBurntAndDamageAbsorbed(distribution, attackerContext, defenderContext, damageDealtToTarget)

	damageDealtToTarget -= distribution.DamageAbsorbedByBarrier
	distribution.TotalRawBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt

	distribution.DamageAbsorbedByArmor = context.calculateDamageAbsorbedByArmor(attackerContext, defenderContext, damageDealtToTarget)
	damageDealtToTarget -= distribution.DamageAbsorbedByArmor

	distribution.RawDamageDealt = damageDealtToTarget
	distribution.IsFatalToTarget = distribution.RawDamageDealt >= defenderContext.HitPoints

	return distribution
}

func (context *VersusContext) calculateDamageAbsorbedByArmor(attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) int {
	if attackerContext.DamageType != power.Physical {
		return 0
	}

	armorAbsorbsAllDamage := damageDealtToTarget <= defenderContext.ArmorResistance
	if armorAbsorbsAllDamage {
		return damageDealtToTarget
	}
	return defenderContext.ArmorResistance
}

func (context *VersusContext) setBarrierBurntAndDamageAbsorbed(distribution *damagedistribution.DamageDistribution, attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) {
	barrierAbsorbsAllDamageAndExtraBurn := damageDealtToTarget + attackerContext.ExtraBarrierBurn <= defenderContext.BarrierResistance
	if barrierAbsorbsAllDamageAndExtraBurn {
		distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
		distribution.DamageAbsorbedByBarrier = damageDealtToTarget
		distribution.TotalRawBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt
		return
	}

	barrierAbsorbsExtraBarrierBurn := attackerContext.ExtraBarrierBurn <= defenderContext.BarrierResistance
	if !barrierAbsorbsExtraBarrierBurn {
		distribution.ExtraBarrierBurnt = defenderContext.BarrierResistance
		distribution.DamageAbsorbedByBarrier = 0
		distribution.TotalRawBarrierBurnt = distribution.ExtraBarrierBurnt
		return
	}

	barrierRemainingAfterExtraBarrierBurn := defenderContext.BarrierResistance - attackerContext.ExtraBarrierBurn

	remainingBarrierAbsorbsDamage := damageDealtToTarget <= barrierRemainingAfterExtraBarrierBurn
	if remainingBarrierAbsorbsDamage {
		distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
		distribution.DamageAbsorbedByBarrier = damageDealtToTarget
		distribution.TotalRawBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt
		return
	}

	distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn
	distribution.DamageAbsorbedByBarrier = barrierRemainingAfterExtraBarrierBurn
	distribution.TotalRawBarrierBurnt = defenderContext.BarrierResistance
}

func (context *VersusContext) setCriticalHitChance(attackerContext AttackerContext) {
	context.CanCritical = attackerContext.CanCritical
	if context.CanCritical {
		context.CriticalHitThreshold = attackerContext.CriticalHitThreshold
	} else {
		context.CriticalHitThreshold = 0
	}
}