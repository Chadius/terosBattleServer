package levelup

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
)

// improveSquaddieStats improves the Squaddie by using the LevelUpBenefit.
func improveSquaddieStats(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Defense != nil {
		squaddieToImprove.Defense.MaxHitPoints = squaddieToImprove.Defense.MaxHitPoints + benefit.Defense.MaxHitPoints
		squaddieToImprove.Defense.Dodge = squaddieToImprove.Defense.Dodge + benefit.Defense.Dodge
		squaddieToImprove.Defense.Deflect = squaddieToImprove.Defense.Deflect + benefit.Defense.Deflect
		squaddieToImprove.Defense.MaxBarrier = squaddieToImprove.Defense.MaxBarrier + benefit.Defense.MaxBarrier
		squaddieToImprove.Defense.Armor = squaddieToImprove.Defense.Armor + benefit.Defense.Armor
	}

	if benefit.Offense != nil {
		squaddieToImprove.Offense.Aim = squaddieToImprove.Offense.Aim + benefit.Offense.Aim
		squaddieToImprove.Offense.Strength = squaddieToImprove.Offense.Strength + benefit.Offense.Strength
		squaddieToImprove.Offense.Mind = squaddieToImprove.Offense.Mind + benefit.Offense.Mind
	}
}

// ImproveSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, powerRepo *power.Repository) error {
	if squaddieToImprove.ClassProgress.HasAddedClass(benefit.Identification.ClassID) == false {
		return fmt.Errorf(`squaddie "%s" cannot add levels to unknown class "%s"`, squaddieToImprove.Identification.Name, benefit.Identification.ClassID)
	}
	if squaddieToImprove.ClassProgress.IsClassLevelAlreadyUsed(benefit.Identification.ID) {
		return fmt.Errorf(`%s already consumed LevelUpBenefit - class:"%s" id:"%s"`, squaddieToImprove.Identification.Name, benefit.Identification.ClassID, benefit.Identification.ID)
	}

	improveSquaddieStats(benefit, squaddieToImprove)
	err := refreshSquaddiePowers(benefit, squaddieToImprove, powerRepo)
	if err != nil {
		return err
	}

	improveSquaddieMovement(benefit, squaddieToImprove)

	squaddieToImprove.ClassProgress.SetBaseClassIfNoBaseClass(benefit.Identification.ClassID)
	squaddieToImprove.ClassProgress.MarkLevelUpBenefitAsConsumed(benefit.Identification.ClassID, benefit.Identification.ID)
	return nil
}

func refreshSquaddiePowers(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, powerRepo *power.Repository) error {
	initialSquaddiePowerReferences := squaddieToImprove.PowerCollection.PowerReferences
	if initialSquaddiePowerReferences == nil || len(initialSquaddiePowerReferences) == 0 {
		initialSquaddiePowerReferences = squaddieToImprove.PowerCollection.GetInnatePowerIDNames()
	}

	powerReferencesToKeep := squaddie.FilterPowerID(initialSquaddiePowerReferences, func(existingPower *power.Reference) bool {
		return squaddie.ContainsPowerID(benefit.PowerIDLost, existingPower.ID) == false
	})
	powerReferencesToLoad := append(powerReferencesToKeep, benefit.PowerIDGained...)

	_, err := powerequip.LoadAllOfSquaddieInnatePowers(squaddieToImprove, powerReferencesToLoad, powerRepo)
	return err
}

func improveSquaddieMovement(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Movement == nil {
		return
	}

	squaddieToImprove.Movement.Distance = squaddieToImprove.Movement.Distance + benefit.Movement.Distance

	if squaddie.MovementValueByType[squaddieToImprove.Movement.Type] < squaddie.MovementValueByType[benefit.Movement.Type] {
		squaddieToImprove.Movement.Type = benefit.Movement.Type
	}

	if benefit.Movement.HitAndRun {
		squaddieToImprove.Movement.HitAndRun = benefit.Movement.HitAndRun
	}
}
