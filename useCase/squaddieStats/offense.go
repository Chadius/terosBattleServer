package squaddiestats

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

func getSquaddieAndPower(squaddieID, powerID string, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) (*squaddie.Squaddie, *power.Power, error) {
	squaddie := squaddieRepo.GetOriginalSquaddieByID(squaddieID)
	if squaddie == nil {
		return nil, nil, fmt.Errorf("squaddie could not be found, ID: %s", squaddieID)
	}
	power := powerRepo.GetPowerByID(powerID)
	if power == nil {
		return nil, nil, fmt.Errorf("power could not be found, ID: %s", powerID)
	}
	if power.AttackEffect == nil {
		return nil, nil, fmt.Errorf("cannot attack with power, ID: %s", powerID)
	}
	return squaddie, power, nil
}

// GetSquaddieAimWithPower returns the to hit bonus against a target without dodge.
func GetSquaddieAimWithPower(squaddieID, powerID string, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) (int, error) {
	squaddie, power, err := getSquaddieAndPower(squaddieID, powerID, squaddieRepo, powerRepo)
	if err != nil {
		return 0, err
	}

	return squaddie.Offense.Aim + power.AttackEffect.ToHitBonus, nil
}

// GetSquaddieRawDamageWithPower returns the amount of damage that will be dealt to an unprotected target.
func GetSquaddieRawDamageWithPower(squaddieID, powerID string, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, squaddieRepo, powerRepo)
	if err != nil {
		return 0, err
	}

	if powerToMeasure.PowerType == power.Physical {
		return squaddie.Offense.Strength + powerToMeasure.AttackEffect.DamageBonus, nil
	}
	return squaddie.Offense.Mind + powerToMeasure.AttackEffect.DamageBonus, nil
}
