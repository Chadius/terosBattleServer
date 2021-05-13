package powerequip

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetEquippedPower returns the power the squaddie has equipped.
//   May return nil if the squaddie hasn't equipped anything.
func GetEquippedPower (squaddie *squaddie.Squaddie, repo *power.Repository) *power.Power {
	if squaddie.PowerCollection.CurrentlyEquippedPowerID != "" {
		return repo.GetPowerByID(squaddie.PowerCollection.CurrentlyEquippedPowerID)
	}

	for _, powerReference := range squaddie.PowerCollection.PowerReferences {
		powerToCheck := repo.GetPowerByID(powerReference.ID)
		if powerToCheck.AttackEffect != nil && powerToCheck.AttackEffect.CanBeEquipped == true {
			return powerToCheck
		}
	}
	return nil
}

// SquaddieEquipPower will make the Squaddie equip a different power.
//   returns true upon success
func SquaddieEquipPower(squaddie *squaddie.Squaddie, powerToEquipID string, repo *power.Repository) bool {
	if squaddie.PowerCollection.HasPowerWithID(powerToEquipID) == false {
		return false
	}

	powerToEquip := repo.GetPowerByID(powerToEquipID)
	if powerToEquip == nil {
		return false
	}
	if powerToEquip.AttackEffect.CanBeEquipped == false {
		return false
	}

	squaddie.PowerCollection.CurrentlyEquippedPowerID = powerToEquipID
	return true
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func LoadAllOfSquaddieInnatePowers(squaddie *squaddie.Squaddie, powerReferencesToLoad []*power.Reference, repo *power.Repository) (int, error) {
	numberOfPowersAdded := 0

	squaddie.PowerCollection.ClearInnatePowers()
	squaddie.PowerCollection.ClearTemporaryPowerReferences()

	for _, powerIDName := range powerReferencesToLoad {
		powerToAdd := repo.GetPowerByID(powerIDName.ID)
		if powerToAdd == nil {
			return numberOfPowersAdded, fmt.Errorf("squaddie '%s' tried to add Power '%s' but it does not exist", squaddie.Identification.Name, powerIDName.Name)
		}

		err := squaddie.PowerCollection.AddInnatePower(powerToAdd)
		if err == nil {
			numberOfPowersAdded = numberOfPowersAdded + 1
		}
	}

	return numberOfPowersAdded, nil
}