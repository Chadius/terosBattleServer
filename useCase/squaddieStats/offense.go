package squaddiestats

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
)

func getSquaddieAndPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (*squaddie.Squaddie, *power.Power, error) {
	squaddie := repos.SquaddieRepo.GetOriginalSquaddieByID(squaddieID)
	if squaddie == nil {
		return nil, nil, fmt.Errorf("squaddie could not be found, ID: %s", squaddieID)
	}
	power := repos.PowerRepo.GetPowerByID(powerID)
	if power == nil {
		return nil, nil, fmt.Errorf("power could not be found, ID: %s", powerID)
	}
	if power.AttackEffect == nil {
		return nil, nil, fmt.Errorf("cannot attack with power, ID: %s", powerID)
	}
	return squaddie, power, nil
}

// GetSquaddieAimWithPower returns the to hit bonus against a target without dodge.
func GetSquaddieAimWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return squaddie.Offense.Aim + powerToMeasure.AttackEffect.ToHitBonus, nil
}

// GetSquaddieRawDamageWithPower returns the amount of damage that will be dealt to an unprotected target.
func GetSquaddieRawDamageWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	if powerToMeasure.PowerType == power.Physical {
		return squaddie.Offense.Strength + powerToMeasure.AttackEffect.DamageBonus, nil
	}
	return squaddie.Offense.Mind + powerToMeasure.AttackEffect.DamageBonus, nil
}

// GetSquaddieCriticalThresholdWithPower returns the critical hit threshold the squaddie needs to beat in order to crit.
func GetSquaddieCriticalThresholdWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	_, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	if powerToMeasure.CanCriticallyHit() != true {
		return 0, fmt.Errorf("cannot critical hit with power, ID: %s", powerID)
	}

	return powerToMeasure.AttackEffect.CriticalEffect.CriticalHitThreshold(), nil
}

// GetSquaddieCriticalRawDamageWithPower returns the total critical hit damage the squaddie deals (includes raw damage.)
func GetSquaddieCriticalRawDamageWithPower (squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	_, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	rawDamage, rawDamageErr := GetSquaddieRawDamageWithPower(squaddieID, powerID, repos)
	if rawDamageErr != nil {
		return 0, rawDamageErr
	}

	if powerToMeasure.CanCriticallyHit() != true {
		return 0, fmt.Errorf("cannot critical hit with power, ID: %s", powerID)
	}

	return rawDamage + powerToMeasure.AttackEffect.CriticalEffect.ExtraCriticalHitDamage(), nil
}

// GetSquaddieCanCounterAttackWithPower returns true if the squaddie can counter attack with this power.
func GetSquaddieCanCounterAttackWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (bool, error) {
	_, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return false, err
	}

	return powerToMeasure.AttackEffect.CanCounterAttack, nil
}

// GetSquaddieCounterAttackAimWithPower returns the to critical hit bonus against a target without dodge.
func GetSquaddieCounterAttackAimWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return squaddie.Offense.Aim + powerToMeasure.AttackEffect.ToHitBonus + powerToMeasure.AttackEffect.CounterAttackToHitPenalty, nil
}

// GetSquaddieExtraBarrierBurnWithPower returns the amount of extra barrier burn that will be dealt to a target with a barrier.
func GetSquaddieExtraBarrierBurnWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	_, powerToMeasure, err := getSquaddieAndPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return powerToMeasure.AttackEffect.ExtraBarrierBurn, nil
}
