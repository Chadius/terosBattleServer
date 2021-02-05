package powerAttack

import (
	"fmt"
	powerPackage "github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetPowerToHitBonusWhenUsedBySquaddie calculates the total to hit bonus for the attacking squaddie and attacking power
func GetPowerToHitBonusWhenUsedBySquaddie(power *powerPackage.Power, squaddie *squaddie.Squaddie) (toHit int) {
	return power.AttackingEffect.ToHitBonus + squaddie.Aim
}

// GetPowerDamageBonusWhenUsedBySquaddie calculates the total Damage bonus for the attacking squaddie and attacking power
func GetPowerDamageBonusWhenUsedBySquaddie(power *powerPackage.Power, squaddie *squaddie.Squaddie) (damageBonus int) {
	if power.PowerType == powerPackage.Physical {
		return power.AttackingEffect.DamageBonus + squaddie.Strength
	}
	return power.AttackingEffect.DamageBonus + squaddie.Mind
}

// GetPowerCriticalDamageBonusWhenUsedBySquaddie calculates the total Critical Hit Damage bonus for the attacking squaddie and attacking power
func GetPowerCriticalDamageBonusWhenUsedBySquaddie(power *powerPackage.Power, squaddie *squaddie.Squaddie) (damageBonus int) {
	return 2 * GetPowerDamageBonusWhenUsedBySquaddie(power, squaddie)
}

// GetHowTargetDistributesDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesDamage(power *powerPackage.Power, attacker *squaddie.Squaddie, target *squaddie.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerDamageBonusWhenUsedBySquaddie(power, attacker)
	return calculateHowTargetTakesDamage(power, target, damageToAbsorb)
}

// GetHowTargetDistributesCriticalDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesCriticalDamage(power *powerPackage.Power, attacker *squaddie.Squaddie, target *squaddie.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerCriticalDamageBonusWhenUsedBySquaddie(power, attacker)
	return calculateHowTargetTakesDamage(power, target, damageToAbsorb)
}

// calculateHowTargetTakesDamage factors the target's damage reduction to figure out how the damage is split between barrier, armor and health.
func calculateHowTargetTakesDamage(power *powerPackage.Power, target *squaddie.Squaddie, damageToAbsorb int) (healthDamage, barrierDamage, extraBarrierDamage int) {
	remainingBarrier := target.CurrentBarrier

	damageToAbsorb, barrierDamage, remainingBarrier = calculateDamageAfterInitialBarrierAbsorption(target, damageToAbsorb, barrierDamage, remainingBarrier)

	extraBarrierDamage = calculateDamageAfterExtraBarrierDamage(power, remainingBarrier, extraBarrierDamage)

	healthDamage = calculateDamageAfterArmorAbsorption(power, target, damageToAbsorb, healthDamage)

	return healthDamage, barrierDamage, extraBarrierDamage
}

func calculateDamageAfterArmorAbsorption(power *powerPackage.Power, target *squaddie.Squaddie, damageToAbsorb int, healthDamage int) int {
	var armorCanAbsorbDamage bool = power.PowerType == powerPackage.Physical
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

func calculateDamageAfterExtraBarrierDamage(power *powerPackage.Power, remainingBarrier int, extraBarrierDamage int) int {
	if power.ExtraBarrierDamage > 0 {
		var barrierFullyAbsorbsExtraBarrierDamage bool = remainingBarrier > power.ExtraBarrierDamage
		if barrierFullyAbsorbsExtraBarrierDamage {
			extraBarrierDamage = power.ExtraBarrierDamage
			remainingBarrier = remainingBarrier - power.ExtraBarrierDamage
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

// AttackingPowerSummary gives a summary of the chance to hit and damage dealt by attacks. Expected damage counts the number of 36ths so we can use ints for fractional math.
type AttackingPowerSummary struct {
	ChanceToHit                   int
	DamageTaken                   int
	ExpectedDamage                int
	BarrierDamageTaken            int
	ExpectedBarrierDamage         int
	ChanceToCrit                  int
	CriticalDamageTaken           int
	CriticalBarrierDamageTaken    int
	CriticalExpectedDamage        int
	CriticalExpectedBarrierDamage int
}

// GetExpectedDamage provides a quick summary of an attack as well as the multiplied estimate
func GetExpectedDamage(power *powerPackage.Power, attacker *squaddie.Squaddie, target *squaddie.Squaddie) (battleSummary *AttackingPowerSummary) {
	toHitBonus := GetPowerToHitBonusWhenUsedBySquaddie(power, attacker)
	toHitPenalty := GetPowerToHitPenaltyAgainstSquaddie(power, target)
	totalChanceToHit := powerPackage.GetChanceToHitBasedOnHitRate(toHitBonus - toHitPenalty)

	healthDamage, barrierDamage, extraBarrierDamage := GetHowTargetDistributesDamage(power, attacker, target)

	chanceToCritical := powerPackage.GetChanceToCriticalBasedOnThreshold(power.CriticalHitThreshold)
	var criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage int
	if chanceToCritical > 0 {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = GetHowTargetDistributesCriticalDamage(power, attacker, target)
	} else {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = 0, 0, 0
	}

	return &AttackingPowerSummary{
		ChanceToHit:                   totalChanceToHit,
		DamageTaken:                   healthDamage,
		ExpectedDamage:                totalChanceToHit * healthDamage,
		BarrierDamageTaken:            barrierDamage + extraBarrierDamage,
		ExpectedBarrierDamage:         totalChanceToHit * (barrierDamage + extraBarrierDamage),
		ChanceToCrit:                  chanceToCritical,
		CriticalDamageTaken:           criticalHealthDamage,
		CriticalBarrierDamageTaken:    criticalBarrierDamage + criticalExtraBarrierDamage,
		CriticalExpectedDamage:        totalChanceToHit * criticalHealthDamage,
		CriticalExpectedBarrierDamage: totalChanceToHit * (criticalBarrierDamage + criticalExtraBarrierDamage),
	}
}

// GetPowerToHitPenaltyAgainstSquaddie calculates how much the target can reduce the chance of getting hit by the attacking power.
func GetPowerToHitPenaltyAgainstSquaddie(power *powerPackage.Power, target *squaddie.Squaddie) (toHitPenalty int) {
	if power.PowerType == powerPackage.Physical {
		return target.Dodge
	}
	return target.Deflect
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func LoadAllOfSquaddieInnatePowers(squaddie *squaddie.Squaddie, powerReferencesToLoad []*powerPackage.Reference, repo *powerPackage.Repository) (int, error) {
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