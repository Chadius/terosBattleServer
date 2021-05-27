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

func (context *DefenderContext) getPower(setup *powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) *power.Power {
	return repositories.PowerRepo.GetPowerByID(setup.PowerID)
}

func (context *DefenderContext) calculate(setup *powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) {
	context.TotalToHitPenalty = context.calculateTotalToHitPenalty(setup, repositories)
	context.ArmorResistance = context.calculateArmorResistance(setup, repositories)
	context.BarrierResistance = context.calculateBarrierResistance(repositories)
}

func (context *DefenderContext) calculateTotalToHitPenalty(setup *powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) int {
	attackingPower := context.getPower(setup, repositories)
	target := repositories.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Dodge
	}

	if attackingPower.PowerType == power.Spell {
		return target.Defense.Deflect
	}
	return 0
}

func (context *DefenderContext) calculateArmorResistance(setup *powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) int {
	attackingPower := context.getPower(setup, repositories)
	target := repositories.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Armor
	}
	return 0
}

func (context *DefenderContext) calculateBarrierResistance(repositories *powerusagescenario.RepositoryCollection) int {
	target := repositories.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)
	return target.Defense.CurrentBarrier
}
