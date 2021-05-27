package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
)

// AttackerContext lists the attacker's relevant information when attacking
type AttackerContext struct {
	IsCounterAttack bool
	CounterAttackPenalty int
	TotalToHitBonus int

	RawDamage       int
	DamageType      power.Type

	ExtraBarrierBurn int

	CanCritical bool
	CriticalHitThreshold int
	CriticalHitDamage int
}

func (context *AttackerContext)calculate(setup powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) {
	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)

	context.DamageType = power.PowerType
	context.ExtraBarrierBurn = power.AttackEffect.ExtraBarrierBurn

	context.RawDamage = context.calculateRawDamage(setup, repositories)
	context.calculateToHitBonus(setup, repositories)

	context.calculateCriticalHit(setup, repositories)
}

func (context *AttackerContext) calculateToHitBonus(setup powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) {
	user := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)

	context.IsCounterAttack = setup.IsCounterAttack
	context.TotalToHitBonus = power.AttackEffect.ToHitBonus + user.Offense.Aim
	if context.IsCounterAttack {
		context.CounterAttackPenalty = power.AttackEffect.CounterAttackToHitPenalty
		context.TotalToHitBonus -= context.CounterAttackPenalty
	}
}

func (context *AttackerContext) calculateRawDamage(setup powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) int {
	user := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	powerToAttackWith := repositories.PowerRepo.GetPowerByID(setup.PowerID)
	if powerToAttackWith.PowerType == power.Physical {
		return powerToAttackWith.AttackEffect.DamageBonus + user.Offense.Strength
	}

	if powerToAttackWith.PowerType == power.Spell {
		return powerToAttackWith.AttackEffect.DamageBonus + user.Offense.Mind
	}
	return 0
}

func (context *AttackerContext) calculateCriticalHit(setup powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) {
	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)
	if power.AttackEffect.CriticalHitThreshold == 0 {
		context.CanCritical = false
		return
	}

	context.CanCritical = true
	context.CriticalHitThreshold = power.AttackEffect.CriticalHitThreshold
	context.CriticalHitDamage = context.RawDamage * 2
}