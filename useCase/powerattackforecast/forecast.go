package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powercounter"
)

// Forecast will store the information needed to explain what will happen when a squaddie
//  uses a given power. It can be asked multiple questions.
type Forecast struct {
	Setup	ForecastSetup
	ForecastedResultPerTarget []Calculation
}

// ForecastSetup is supplied upon creation to explain all of the relevant parts of this power.
type ForecastSetup struct {
	UserID          string
	PowerID         string
	Targets        []string
	SquaddieRepo    *squaddie.Repository
	PowerRepo       *power.Repository
	IsCounterAttack bool
}

// Calculation holds the results of the forecast.
type Calculation struct {
	Attack	*AttackForecast
	CounterAttack *AttackForecast
}

// AttackForecast shows what will happen if the power used is offensive.
type AttackForecast struct {
	AttackerContext AttackerContext
	DefenderContext DefenderContext
	VersusContext VersusContext
}

// CalculateForecast uses the given Setup and anticipates what will happen
//  when the power is used.
func (forecast *Forecast) CalculateForecast() {
	for _, targetID := range forecast.Setup.Targets {
		attack := forecast.CalculateAttackForecast(targetID)
		var counterAttack *AttackForecast
		if forecast.isCounterattackPossible(targetID) {
			counterAttack = forecast.createCounterAttackForecast(targetID)
		}

		calculation := Calculation{
			Attack: attack,
			CounterAttack: counterAttack,
		}
		forecast.ForecastedResultPerTarget = append(forecast.ForecastedResultPerTarget, calculation)
	}
}

func (forecast *Forecast) isCounterattackPossible(targetID string) bool {
	squaddieThatWantsToCounter := forecast.Setup.SquaddieRepo.GetOriginalSquaddieByID(targetID)
	if forecast.Setup.IsCounterAttack == false && powercounter.CanSquaddieCounterWithEquippedWeapon(squaddieThatWantsToCounter, forecast.Setup.PowerRepo) {
		return true
	}
	return false
}

func (forecast *Forecast) createCounterAttackForecast(targetID string) *AttackForecast {
	counterAttackingSquaddie := forecast.Setup.SquaddieRepo.GetOriginalSquaddieByID(targetID)
	counterAttackingPowerID := counterAttackingSquaddie.PowerCollection.CurrentlyEquippedPowerID

	counterForecastSetup := ForecastSetup{
		UserID:          targetID,
		PowerID:         counterAttackingPowerID,
		Targets:         []string{forecast.Setup.UserID},
		SquaddieRepo:    forecast.Setup.SquaddieRepo,
		PowerRepo:       forecast.Setup.PowerRepo,
		IsCounterAttack: true,
	}

	counterAttackForecast := Forecast{
		Setup:                     counterForecastSetup,
	}

	counterAttackForecast.CalculateForecast()

	return counterAttackForecast.CalculateAttackForecast(targetID)
}

// CalculateAttackForecast figures out what will happen when this attack power is used.
func (forecast *Forecast) CalculateAttackForecast(targetID string) *AttackForecast {
	attackerContext := AttackerContext{AttackerID: forecast.Setup.UserID}
	attackerContext.calculate(forecast.Setup)

	defenderContext := DefenderContext{TargetID: targetID}
	defenderContext.calculate(&forecast.Setup)

	versusContext := VersusContext{}
	versusContext.calculate(attackerContext, defenderContext)

	return &AttackForecast{
		AttackerContext: attackerContext,
		DefenderContext: defenderContext,
		VersusContext: versusContext,
	}
}
