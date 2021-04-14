package powerusage

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetPowerToHitBonusWhenUsedBySquaddie calculates the total to hit bonus for the attacking squaddie and attacking power
func GetPowerToHitBonusWhenUsedBySquaddie(attackingPower *power.Power, squaddie *squaddie.Squaddie, isCounterAttack bool) (toHit int) {
	counterAttackPenalty := 0
	if isCounterAttack {
		counterAttackPenalty = attackingPower.AttackEffect.CounterAttackToHitPenalty
	}
	return attackingPower.AttackEffect.ToHitBonus + squaddie.Aim + counterAttackPenalty
}

// GetPowerDamageBonusWhenUsedBySquaddie calculates the total Damage bonus for the attacking squaddie and attacking power
func GetPowerDamageBonusWhenUsedBySquaddie(attackingPower *power.Power, squaddie *squaddie.Squaddie) (damageBonus int) {
	if attackingPower.PowerType == power.Physical {
		return attackingPower.AttackEffect.DamageBonus + squaddie.Strength
	}
	return attackingPower.AttackEffect.DamageBonus + squaddie.Mind
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

// GetHowTargetDistributesCriticalDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesCriticalDamage(attackingPower *power.Power, attacker *squaddie.Squaddie, target *squaddie.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerCriticalDamageBonusWhenUsedBySquaddie(attackingPower, attacker)
	return calculateHowTargetTakesDamage(attackingPower, target, damageToAbsorb)
}

// calculateHowTargetTakesDamage factors the target's damage reduction to figure out how the damage is split between barrier, armor and health.
func calculateHowTargetTakesDamage(attackingPower *power.Power, target *squaddie.Squaddie, damageToAbsorb int) (healthDamage, barrierDamage, extraBarrierDamage int) {
	remainingBarrier := target.CurrentBarrier

	damageToAbsorb, barrierDamage, remainingBarrier = calculateDamageAfterInitialBarrierAbsorption(target, damageToAbsorb, barrierDamage, remainingBarrier)

	extraBarrierDamage = calculateDamageAfterExtraBarrierDamage(attackingPower, remainingBarrier, extraBarrierDamage)

	healthDamage = calculateDamageAfterArmorAbsorption(attackingPower, target, damageToAbsorb, healthDamage)

	return healthDamage, barrierDamage, extraBarrierDamage
}

func calculateDamageAfterArmorAbsorption(attackingPower *power.Power, target *squaddie.Squaddie, damageToAbsorb int, healthDamage int) int {
	var armorCanAbsorbDamage bool = attackingPower.PowerType == power.Physical
	if armorCanAbsorbDamage {

		var armorFullyAbsorbsDamage bool = target.Armor > damageToAbsorb
		if armorFullyAbsorbsDamage {
			healthDamage = 0
		} else {
			healthDamage = damageToAbsorb - target.Armor
		}
	} else {
		healthDamage = damageToAbsorb
	}
	return healthDamage
}

func calculateDamageAfterExtraBarrierDamage(attackingPower *power.Power, remainingBarrier int, extraBarrierDamage int) int {
	if attackingPower.AttackEffect.ExtraBarrierDamage > 0 {
		var barrierFullyAbsorbsExtraBarrierDamage bool = remainingBarrier > attackingPower.AttackEffect.ExtraBarrierDamage
		if barrierFullyAbsorbsExtraBarrierDamage {
			extraBarrierDamage = attackingPower.AttackEffect.ExtraBarrierDamage
			remainingBarrier = remainingBarrier - attackingPower.AttackEffect.ExtraBarrierDamage
		} else {
			extraBarrierDamage = remainingBarrier
			remainingBarrier = 0
		}
	}
	return extraBarrierDamage
}

func calculateDamageAfterInitialBarrierAbsorption(target *squaddie.Squaddie, damageToAbsorb int, barrierDamage int, remainingBarrier int) (int, int, int) {
	var barrierFullyAbsorbsDamage bool = target.CurrentBarrier > damageToAbsorb
	if barrierFullyAbsorbsDamage {
		barrierDamage = damageToAbsorb
		remainingBarrier = remainingBarrier - barrierDamage
		damageToAbsorb = 0
	} else {
		barrierDamage = target.CurrentBarrier
		remainingBarrier = 0
		damageToAbsorb = damageToAbsorb - target.CurrentBarrier
	}
	return damageToAbsorb, barrierDamage, remainingBarrier
}

// GetExpectedDamage provides a summary of what the attacker's attackingPower will do against the given target.
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
	if isCounterAttack == false && CanTargetSquaddieCounterAttack(context, attackContext) {
		counterAttackContext := attackContext.Clone()
		counterAttackContext.AttackerID = attackContext.TargetID
		counterAttackContext.TargetID = attackContext.AttackerID
		counterAttackContext.IsCounterAttack = true
		counterAttackContext.PowerID = GetEquippedPower(target, context.PowerRepo).ID
		counterAttackSummary = GetExpectedDamage(context, counterAttackContext)
	}

	return &powerusagecontext.AttackingPowerForecast{
		AttackingSquaddieID:			attacker.ID,
		PowerID:						attackingPower.ID,
		TargetSquaddieID: 				target.ID,
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

// CanTargetSquaddieCounterAttack returns true if the target can counterAttack the attacker.
func CanTargetSquaddieCounterAttack(context *powerusagecontext.PowerUsageContext, attackContext *powerusagecontext.AttackContext) bool {
	target := context.SquaddieRepo.GetOriginalSquaddieByID(attackContext.TargetID)
	return CanSquaddieCounterWithEquippedWeapon(target, context.PowerRepo)
}

// GetPowerToHitPenaltyAgainstSquaddie calculates how much the target can reduce the chance of getting hit by the attacking power.
func GetPowerToHitPenaltyAgainstSquaddie(attackingPower *power.Power, target *squaddie.Squaddie) (toHitPenalty int) {
	if attackingPower.PowerType == power.Physical {
		return target.Dodge
	}
	return target.Deflect
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func LoadAllOfSquaddieInnatePowers(squaddie *squaddie.Squaddie, powerReferencesToLoad []*power.Reference, repo *power.Repository) (int, error) {
	numberOfPowersAdded := 0

	squaddie.ClearInnatePowers()
	squaddie.ClearTemporaryPowerReferences()

	for _, powerIDName := range powerReferencesToLoad {
		powerToAdd := repo.GetPowerByID(powerIDName.ID)
		if powerToAdd == nil {
			return numberOfPowersAdded, fmt.Errorf("squaddie '%s' tried to add Power '%s' but it does not exist", squaddie.Name, powerIDName.Name)
		}

		err := squaddie.AddInnatePower(powerToAdd)
		if err == nil {
			numberOfPowersAdded = numberOfPowersAdded + 1
		}
	}

	return numberOfPowersAdded, nil
}