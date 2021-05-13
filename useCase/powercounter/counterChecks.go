package powercounter

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
)

// CanSquaddieCounterWithEquippedWeapon returns true if the squaddie can use the currently equipped
//   weapon for counter attacks.
func CanSquaddieCounterWithEquippedWeapon(squaddie *squaddie.Squaddie, repo *power.Repository) bool {
	currentlyEquippedPower := powerequip.GetEquippedPower(squaddie, repo)
	if currentlyEquippedPower == nil {
		return false
	}
	return currentlyEquippedPower.AttackEffect.CanCounterAttack
}

// TODO delete this
// CanTargetSquaddieCounterAttack returns true if the target can counterAttack the attacker.
func CanTargetSquaddieCounterAttack(context *powerusagecontext.PowerUsageContext, attackContext *powerusagecontext.AttackContext) bool {
	target := context.SquaddieRepo.GetOriginalSquaddieByID(attackContext.TargetID)
	return CanSquaddieCounterWithEquippedWeapon(target, context.PowerRepo)
}

