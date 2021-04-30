package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// DefenderContext lists the target's relevant information when under attack
type DefenderContext struct {
	TargetID string

	TotalToHitPenalty int

	ArmorResistance int
	BarrierResistance int
}

func (context *DefenderContext) getPower(setup ForecastSetup) *power.Power {
	return setup.PowerRepo.GetPowerByID(setup.PowerID)
}

func (context *DefenderContext) getTarget(setup ForecastSetup) *squaddie.Squaddie {
	return setup.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)
}

func (context *DefenderContext) calculate(setup ForecastSetup) {
	context.TotalToHitPenalty = context.calculateTotalToHitPenalty(setup)
	context.ArmorResistance = context.calculateArmorResistance(setup)
	context.BarrierResistance = context.calculateBarrierResistance(setup)
}

func (context *DefenderContext) calculateTotalToHitPenalty(setup ForecastSetup) int {
	attackingPower := context.getPower(setup)
	target := context.getTarget(setup)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Dodge
	}

	if attackingPower.PowerType == power.Spell {
		return target.Defense.Deflect
	}
	return 0
}

func (context *DefenderContext) calculateArmorResistance(setup ForecastSetup) int {
	attackingPower := context.getPower(setup)
	target := context.getTarget(setup)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Armor
	}
	return 0
}

func (context *DefenderContext) calculateBarrierResistance(setup ForecastSetup) int {
	target := context.getTarget(setup)
	return target.Defense.CurrentBarrier
}
