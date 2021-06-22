package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
)

// DefenderContext lists the target's relevant information when under attack
type DefenderContext struct {
	TargetID string

	TotalToHitPenalty int

	HitPoints int
	ArmorResistance int
	BarrierResistance int
}

func (context *DefenderContext) getPower(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) *power.Power {
	return repositories.PowerRepo.GetPowerByID(setup.PowerID)
}

func (context *DefenderContext) calculate(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) {
	context.TotalToHitPenalty = context.calculateTotalToHitPenalty(setup, repositories)
	context.ArmorResistance = context.calculateArmorResistance(setup, repositories)
	context.BarrierResistance = context.calculateBarrierResistance(repositories)
	context.HitPoints = context.calculateHitPoints(repositories)
}

func (context *DefenderContext) calculateTotalToHitPenalty(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) int {
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

func (context *DefenderContext) calculateArmorResistance(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) int {
	attackingPower := context.getPower(setup, repositories)
	target := repositories.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)

	if attackingPower.PowerType == power.Physical {
		return target.Defense.Armor
	}
	return 0
}

func (context *DefenderContext) calculateBarrierResistance(repositories *repositories.RepositoryCollection) int {
	target := repositories.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)
	return target.Defense.CurrentBarrier
}

func (context *DefenderContext) calculateHitPoints(repositories *repositories.RepositoryCollection) int {
	target := repositories.SquaddieRepo.GetOriginalSquaddieByID(context.TargetID)
	return target.Defense.CurrentHitPoints
}
