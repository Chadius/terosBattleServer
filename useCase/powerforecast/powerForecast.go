package powerforecast

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powercounter"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
)

// CalculatePowerForecast applies logic to the context and shows the expected forecast of its effect.
func CalculatePowerForecast(context *powerusagecontext.PowerUsageContext) *powerusagecontext.PowerForecast {
	summary := &powerusagecontext.PowerForecast{
		UserSquaddieID:      context.ActingSquaddieID,
		PowerID:             context.PowerID,
		AttackPowerForecast: []*powerusagecontext.AttackingPowerForecast{},
	}

	for _, targetID := range context.TargetSquaddieIDs {
		summary.AttackPowerForecast = append(summary.AttackPowerForecast, GetExpectedDamage(context, &powerusagecontext.AttackContext{
			PowerID:         context.PowerID,
			AttackerID:      context.ActingSquaddieID,
			TargetID:        targetID,
			IsCounterAttack: false,
		}))
	}
	return summary
}

// GetExpectedDamage provides a forecast of what the attacker's attackingPower will do against the given target.
func GetExpectedDamage(
	context *powerusagecontext.PowerUsageContext,
	attackContext *powerusagecontext.AttackContext) (battleSummary *powerusagecontext.AttackingPowerForecast) {

	attackingPower := context.PowerRepo.GetPowerByID(attackContext.PowerID)
	attacker := context.SquaddieRepo.GetOriginalSquaddieByID(attackContext.AttackerID)
	target := context.SquaddieRepo.GetOriginalSquaddieByID(attackContext.TargetID)
	isCounterAttack := attackContext.IsCounterAttack

	toHitBonus := GetPowerToHitBonusWhenUsedBySquaddie(attackingPower, attacker, isCounterAttack)
	toHitPenalty := GetPowerToHitPenaltyAgainstSquaddie(attackingPower, target)
	totalChanceToHit := power.GetChanceToHitBasedOnHitRate(toHitBonus - toHitPenalty)

	healthDamage, barrierDamage, extraBarrierDamage := GetHowTargetDistributesDamage(attackingPower, attacker, target)

	chanceToCritical := power.GetChanceToCriticalBasedOnThreshold(attackingPower.AttackEffect.CriticalHitThreshold)
	var criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage int
	if chanceToCritical > 0 {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = GetHowTargetDistributesCriticalDamage(attackingPower, attacker, target)
	} else {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = 0, 0, 0
	}

	var counterAttackSummary *powerusagecontext.AttackingPowerForecast = nil
	if isCounterAttack == false && powercounter.CanTargetSquaddieCounterAttack(context, attackContext) {
		counterAttackContext := attackContext.Clone()
		counterAttackContext.AttackerID = attackContext.TargetID
		counterAttackContext.TargetID = attackContext.AttackerID
		counterAttackContext.IsCounterAttack = true
		counterAttackContext.PowerID = powerequip.GetEquippedPower(target, context.PowerRepo).ID
		counterAttackSummary = GetExpectedDamage(context, counterAttackContext)
	}

	return &powerusagecontext.AttackingPowerForecast{
		AttackingSquaddieID:			attacker.Identification.ID,
		PowerID:						attackingPower.ID,
		TargetSquaddieID: 				target.Identification.ID,
		CriticalHitThreshold:			attackingPower.AttackEffect.CriticalHitThreshold,
		HitRate:						toHitBonus - toHitPenalty,
		ChanceToHit:					totalChanceToHit,
		DamageTaken:					healthDamage,
		ExpectedDamage:					totalChanceToHit * healthDamage,
		BarrierDamageTaken:				barrierDamage + extraBarrierDamage,
		ExpectedBarrierDamage:			totalChanceToHit * (barrierDamage + extraBarrierDamage),
		ChanceToCritical:				chanceToCritical,
		CriticalDamageTaken:			criticalHealthDamage,
		CriticalBarrierDamageTaken:		criticalBarrierDamage + criticalExtraBarrierDamage,
		CriticalExpectedDamage:			totalChanceToHit * criticalHealthDamage,
		CriticalExpectedBarrierDamage:	totalChanceToHit * (criticalBarrierDamage + criticalExtraBarrierDamage),
		CounterAttack:					counterAttackSummary,
		IsACounterAttack:				isCounterAttack,
	}
}

// GetPowerToHitBonusWhenUsedBySquaddie calculates the total to hit bonus for the attacking squaddie and attacking power
func GetPowerToHitBonusWhenUsedBySquaddie(attackingPower *power.Power, squaddie *squaddie.Squaddie, isCounterAttack bool) (toHit int) {
	counterAttackPenalty := 0
	if isCounterAttack {
		counterAttackPenalty = attackingPower.AttackEffect.CounterAttackToHitPenalty
	}
	return attackingPower.AttackEffect.ToHitBonus + squaddie.Offense.Aim + counterAttackPenalty
}

// GetPowerDamageBonusWhenUsedBySquaddie calculates the total Damage bonus for the attacking squaddie and attacking power
func GetPowerDamageBonusWhenUsedBySquaddie(attackingPower *power.Power, squaddie *squaddie.Squaddie) (damageBonus int) {
	if attackingPower.PowerType == power.Physical {
		return attackingPower.AttackEffect.DamageBonus + squaddie.Offense.Strength
	}
	return attackingPower.AttackEffect.DamageBonus + squaddie.Offense.Mind
}

// GetPowerCriticalDamageBonusWhenUsedBySquaddie calculates the total Critical Hit Damage bonus for the attacking squaddie and attacking power
func GetPowerCriticalDamageBonusWhenUsedBySquaddie(attackingPower *power.Power, squaddie *squaddie.Squaddie) (damageBonus int) {
	return 2 * GetPowerDamageBonusWhenUsedBySquaddie(attackingPower, squaddie)
}

// GetHowTargetDistributesDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesDamage(attackingPower *power.Power, attacker *squaddie.Squaddie, target *squaddie.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerDamageBonusWhenUsedBySquaddie(attackingPower, attacker)
	return calculateHowTargetTakesDamage(attackingPower, target, damageToAbsorb)
}

// GetPowerToHitPenaltyAgainstSquaddie calculates how much the target can reduce the chance of getting hit by the attacking power.
func GetPowerToHitPenaltyAgainstSquaddie(attackingPower *power.Power, target *squaddie.Squaddie) (toHitPenalty int) {
	if attackingPower.PowerType == power.Physical {
		return target.Defense.Dodge
	}
	return target.Defense.Deflect
}

// GetHowTargetDistributesCriticalDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesCriticalDamage(attackingPower *power.Power, attacker *squaddie.Squaddie, target *squaddie.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerCriticalDamageBonusWhenUsedBySquaddie(attackingPower, attacker)
	return calculateHowTargetTakesDamage(attackingPower, target, damageToAbsorb)
}

// calculateHowTargetTakesDamage factors the target's damage reduction to figure out how the damage is split between barrier, armor and health.
func calculateHowTargetTakesDamage(attackingPower *power.Power, target *squaddie.Squaddie, damageToAbsorb int) (healthDamage, barrierDamage, extraBarrierDamage int) {
	remainingBarrier := target.Defense.CurrentBarrier

	damageToAbsorb, barrierDamage, remainingBarrier = calculateDamageAfterInitialBarrierAbsorption(target, damageToAbsorb, barrierDamage, remainingBarrier)

	extraBarrierDamage = calculateDamageAfterExtraBarrierDamage(attackingPower, remainingBarrier, extraBarrierDamage)

	healthDamage = calculateDamageAfterArmorAbsorption(attackingPower, target, damageToAbsorb, healthDamage)

	return healthDamage, barrierDamage, extraBarrierDamage
}

func calculateDamageAfterArmorAbsorption(attackingPower *power.Power, target *squaddie.Squaddie, damageToAbsorb int, healthDamage int) int {
	var armorCanAbsorbDamage bool = attackingPower.PowerType == power.Physical
	if armorCanAbsorbDamage {

		var armorFullyAbsorbsDamage bool = target.Defense.Armor > damageToAbsorb
		if armorFullyAbsorbsDamage {
			healthDamage = 0
		} else {
			healthDamage = damageToAbsorb - target.Defense.Armor
		}
	} else {
		healthDamage = damageToAbsorb
	}
	return healthDamage
}

func calculateDamageAfterExtraBarrierDamage(attackingPower *power.Power, remainingBarrier int, extraBarrierDamage int) int {
	if attackingPower.AttackEffect.ExtraBarrierBurn > 0 {
		var barrierFullyAbsorbsExtraBarrierDamage bool = remainingBarrier > attackingPower.AttackEffect.ExtraBarrierBurn
		if barrierFullyAbsorbsExtraBarrierDamage {
			extraBarrierDamage = attackingPower.AttackEffect.ExtraBarrierBurn
			remainingBarrier = remainingBarrier - attackingPower.AttackEffect.ExtraBarrierBurn
		} else {
			extraBarrierDamage = remainingBarrier
			remainingBarrier = 0
		}
	}
	return extraBarrierDamage
}

func calculateDamageAfterInitialBarrierAbsorption(target *squaddie.Squaddie, damageToAbsorb int, barrierDamage int, remainingBarrier int) (int, int, int) {
	var barrierFullyAbsorbsDamage bool = target.Defense.CurrentBarrier > damageToAbsorb
	if barrierFullyAbsorbsDamage {
		barrierDamage = damageToAbsorb
		remainingBarrier = remainingBarrier - barrierDamage
		damageToAbsorb = 0
	} else {
		barrierDamage = target.Defense.CurrentBarrier
		remainingBarrier = 0
		damageToAbsorb = damageToAbsorb - target.Defense.CurrentBarrier
	}
	return damageToAbsorb, barrierDamage, remainingBarrier
}
