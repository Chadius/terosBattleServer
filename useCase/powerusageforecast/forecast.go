package powerattackusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
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
}

// AttackForecast shows what will happen if the power used is offensive.
type AttackForecast struct {
	TotalToHitBonus	int
}

// CalculateForecast uses the given Setup and anticipates what will happen
//  when the power is used.
func (forecast *Forecast) CalculateForecast() {
	for _, targetID := range forecast.Setup.Targets {
		calculation := Calculation{
			Attack: forecast.CalculateAttackForecast(targetID),
		}
		forecast.ForecastedResultPerTarget = append(forecast.ForecastedResultPerTarget, calculation)
	}
}

// CalculateAttackForecast figures out what will happen when this attack power is used.
func (forecast *Forecast) CalculateAttackForecast(targetID string) *AttackForecast {
	//user := forecast.Setup.SquaddieRepo.GetOriginalSquaddieByID(forecast.Setup.UserID)
	//power := forecast.Setup.PowerRepo.GetPowerByID(forecast.Setup.PowerID)
	//target := forecast.Setup.SquaddieRepo.GetOriginalSquaddieByID(targetID)

	return &AttackForecast{
		TotalToHitBonus: forecast.calculateToHitBonus(),
	}
}

func (forecast *Forecast) calculateToHitBonus() int {
	user := forecast.Setup.SquaddieRepo.GetOriginalSquaddieByID(forecast.Setup.UserID)
	power := forecast.Setup.PowerRepo.GetPowerByID(forecast.Setup.PowerID)
	return power.AttackEffect.ToHitBonus + user.Offense.Aim
}