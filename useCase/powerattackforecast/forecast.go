package powerattackforecast

import (
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
)

// Forecast will store the information needed to explain what will happen when a squaddie
//  uses a given power. It can be asked multiple questions.
type Forecast struct {
	Setup                     powerusagescenario.Setup
	Repositories *repositories.RepositoryCollection
	ForecastedResultPerTarget []Calculation
}

// Calculation holds the results of the forecast.
type Calculation struct {
	Repositories *repositories.RepositoryCollection
	Setup *powerusagescenario.Setup
	Attack	*AttackForecast
	CounterAttackSetup *powerusagescenario.Setup
	CounterAttack *AttackForecast
}

// AttackForecast shows what will happen if the power used is offensive.
type AttackForecast struct {
	AttackerContext AttackerContext
	DefenderContext DefenderContext
	VersusContext VersusContext
}

// CalculateForecast gives a numerical prediction of the power's effect.
func (forecast *Forecast) CalculateForecast() {
	for _, targetID := range forecast.Setup.Targets {
		attack := forecast.CalculateAttackForecast(targetID)
		var counterAttack *AttackForecast
		var counterAttackSetup *powerusagescenario.Setup
		if forecast.isCounterattackPossible(targetID) {
			counterAttackSetup, counterAttack = forecast.createCounterAttackForecast(targetID)
		}

		calculation := Calculation{
			Setup: &powerusagescenario.Setup{
				UserID:          forecast.Setup.UserID,
				PowerID:         forecast.Setup.PowerID,
				Targets:         []string{targetID},
				IsCounterAttack: false,
			},
			Repositories: &repositories.RepositoryCollection{
				SquaddieRepo:    forecast.Repositories.SquaddieRepo,
				PowerRepo:       forecast.Repositories.PowerRepo,
			},
			Attack: attack,
			CounterAttackSetup: counterAttackSetup,
			CounterAttack: counterAttack,
		}
		forecast.ForecastedResultPerTarget = append(forecast.ForecastedResultPerTarget, calculation)
	}
}

func (forecast *Forecast) isCounterattackPossible(targetID string) bool {
	squaddieThatWantsToCounter := forecast.Repositories.SquaddieRepo.GetOriginalSquaddieByID(targetID)
	if forecast.Setup.IsCounterAttack == false && powerequip.CanSquaddieCounterWithEquippedWeapon(squaddieThatWantsToCounter, forecast.Repositories.PowerRepo) {
		return true
	}
	return false
}

func (forecast *Forecast) createCounterAttackForecast(counterAttackingSquaddieID string) (*powerusagescenario.Setup, *AttackForecast) {
	counterAttackingSquaddie := forecast.Repositories.SquaddieRepo.GetOriginalSquaddieByID(counterAttackingSquaddieID)
	counterAttackingPowerID := counterAttackingSquaddie.PowerCollection.CurrentlyEquippedPowerID
	counterAttackingTargetID := forecast.Setup.UserID

	counterForecastSetup := powerusagescenario.Setup{
		UserID:          counterAttackingSquaddieID,
		PowerID:         counterAttackingPowerID,
		Targets:         []string{counterAttackingTargetID},
		IsCounterAttack: true,
	}

	counterAttackForecast := Forecast{
		Setup:                     counterForecastSetup,
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    forecast.Repositories.SquaddieRepo,
			PowerRepo:       forecast.Repositories.PowerRepo,
		},
	}

	counterAttackForecast.CalculateForecast()

	return &counterForecastSetup, counterAttackForecast.CalculateAttackForecast(counterAttackingTargetID)
}

// CalculateAttackForecast figures out what will happen when this attack power is used.
func (forecast *Forecast) CalculateAttackForecast(targetID string) *AttackForecast {
	attackerContext := AttackerContext{}
	attackerContext.calculate(forecast.Setup, forecast.Repositories)

	defenderContext := DefenderContext{TargetID: targetID}
	defenderContext.calculate(&forecast.Setup, forecast.Repositories)

	versusContext := VersusContext{}
	versusContext.calculate(attackerContext, defenderContext)

	return &AttackForecast{
		AttackerContext: attackerContext,
		DefenderContext: defenderContext,
		VersusContext: versusContext,
	}
}
