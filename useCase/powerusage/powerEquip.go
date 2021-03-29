package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetEquippedPower returns the power the squaddie has equipped.
//   May return nil if the squaddie hasn't equipped anything.
func GetEquippedPower (squaddie *squaddie.Squaddie, repo *power.Repository) *power.Power {
	if squaddie.CurrentlyEquippedPowerID != "" {
		return repo.GetPowerByID(squaddie.CurrentlyEquippedPowerID)
	}

	for _, powerReference := range squaddie.PowerReferences {
		powerToCheck := repo.GetPowerByID(powerReference.ID)
		if powerToCheck.AttackEffect != nil && powerToCheck.AttackEffect.CanBeEquipped == true {
			return powerToCheck
		}
	}
	return nil
}

// SquaddieEquipPower will make the Squaddie equip a different power.
//   If the power is invalid, will return nil
func SquaddieEquipPower(squaddie *squaddie.Squaddie, powerToEquipID string, repo *power.Repository) bool {
	if squaddie.HasPowerWithID(powerToEquipID) == false {
		return false
	}

	powerToEquip := repo.GetPowerByID(powerToEquipID)
	if powerToEquip == nil {
		return false
	}
	if powerToEquip.AttackEffect.CanBeEquipped == false {
		return false
	}

	squaddie.CurrentlyEquippedPowerID = powerToEquipID
	return true
}

// CanSquaddieCounterWithEquippedWeapon returns true if the squaddie can use the currently equipped
//   weapon for counter attacks.
func CanSquaddieCounterWithEquippedWeapon(squaddie *squaddie.Squaddie, repo *power.Repository) bool {
	currentlyEquippedPower := GetEquippedPower(squaddie, repo)
	if currentlyEquippedPower == nil {
		return false
	}
	return currentlyEquippedPower.AttackEffect.CanCounterattack
}