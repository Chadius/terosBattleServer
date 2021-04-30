package powerattackforecast

import "github.com/cserrant/terosBattleServer/entity/power"

// AttackerContext lists the attacker's relevant information when attacking
type AttackerContext struct {
	AttackerID		string
	TotalToHitBonus int
	RawDamage       int
	DamageType      power.Type
}

func (context *AttackerContext)calculate(setup ForecastSetup) {
	power := setup.PowerRepo.GetPowerByID(setup.PowerID)

	context.DamageType = power.PowerType
	context.RawDamage = context.calculateRawDamage(setup)
	context.TotalToHitBonus = context.calculateToHitBonus(setup)
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