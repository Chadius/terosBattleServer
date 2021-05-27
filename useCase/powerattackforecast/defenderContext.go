package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
)

// DefenderContext lists the target's relevant information when under attack
type DefenderContext struct {
	TargetID string

	TotalToHitPenalty int

	ArmorResistance int
	BarrierResistance int
}

func (context *DefenderContext) getPower(setup *powerusagescenario.Setup) *power.Power {
	return setup.PowerRepo.GetPowerByID(setup.PowerID)
}

func (context *DefenderContext) calculate(setup *powerusagescenario.Setup) {
	context.TotalToHitPenalty = context.calculateTotalToHitPenalty(setup)
	context.ArmorResistance = context.calculateArmorResistance(setup)
	context.BarrierResistance = context.calculateBarrierResistance(setup)
}

func (context *DefenderContext) calculateTotalToHitPenalty(setup *powerusagescenario.Setup) int {
	attackingPower := context.getPower(setup)
	target := setup.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Dodge
	}

	if attackingPower.PowerType == power.Spell {
		return target.Defense.Deflect
	}
	return 0
}

func (context *DefenderContext) calculateArmorResistance(setup *powerusagescenario.Setup) int {
	attackingPower := context.getPower(setup)
	target := setup.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Armor
	}
	return 0
}

func (context *DefenderContext) calculateBarrierResistance(setup *powerusagescenario.Setup) int {
	target := setup.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)
	return target.Defense.CurrentBarrier
}
