package powerattackforecast

import "github.com/cserrant/terosBattleServer/entity/power"

// AttackerContext lists the attacker's relevant information when attacking
type AttackerContext struct {
	AttackerID		string

	TotalToHitBonus int

	RawDamage       int
	DamageType      power.Type

	ExtraBarrierBurn int

	CanCritical bool
	CriticalHitThreshold int
	CriticalHitDamage int
}

func (context *AttackerContext)calculate(setup ForecastSetup) {
	power := setup.PowerRepo.GetPowerByID(setup.PowerID)

	context.DamageType = power.PowerType
	context.ExtraBarrierBurn = power.AttackEffect.ExtraBarrierBurn

	context.RawDamage = context.calculateRawDamage(setup)
	context.TotalToHitBonus = context.calculateToHitBonus(setup)

	context.calculateCriticalHit(setup)
}

func (context *AttackerContext) calculateToHitBonus(setup ForecastSetup) int {
	user := setup.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	power := setup.PowerRepo.GetPowerByID(setup.PowerID)
	return power.AttackEffect.ToHitBonus + user.Offense.Aim
}

func (context *AttackerContext) calculateRawDamage(setup ForecastSetup) int {
	user := setup.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	powerToAttackWith := setup.PowerRepo.GetPowerByID(setup.PowerID)
	if powerToAttackWith.PowerType == power.Physical {
		return powerToAttackWith.AttackEffect.DamageBonus + user.Offense.Strength
	}

	if powerToAttackWith.PowerType == power.Spell {
		return powerToAttackWith.AttackEffect.DamageBonus + user.Offense.Mind
	}
	return 0
}

func (context *AttackerContext) calculateCriticalHit(setup ForecastSetup) {
	power := setup.PowerRepo.GetPowerByID(setup.PowerID)
	if power.AttackEffect.CriticalHitThreshold == 0 {
		context.CanCritical = false
		return
	}

	context.CanCritical = true
	context.CriticalHitThreshold = power.AttackEffect.CriticalHitThreshold
	context.CriticalHitDamage = context.RawDamage * 2
}