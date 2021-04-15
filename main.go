package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
	"github.com/cserrant/terosBattleServer/utility"
	"io/ioutil"
	"log"
)

func main() {
	squaddieRepo := loadSquaddieRepo()
	powerRepo := loadPowerRepo()

	attacker, target, power := loadActors(
		"squaddieTeros",
		"squaddieBandit",
		"powerSpear",
		squaddieRepo,
		powerRepo,
	)

	powerSummary := powerusage.GetPowerSummary(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      squaddieRepo,
			ActingSquaddieID:  attacker.ID,
			TargetSquaddieIDs: []string{target.ID},
			PowerID:           power.ID,
			PowerRepo:         powerRepo,
		},
		power, attacker, []*squaddie.Squaddie{target})
	attackingPowerSummary := powerSummary.AttackEffectSummary[0]
	println(attacker.Name, "will attack", target.Name, "with", power.Name)
	println("Chance to hit (out of 36) ", attackingPowerSummary.ChanceToHit)
	println("Damage taken              ", attackingPowerSummary.DamageTaken)
	println("Barrier damage            ", attackingPowerSummary.BarrierDamageTaken)
	println("---")
	println("Expected damage (36 = 1HP)", attackingPowerSummary.ExpectedDamage)
	println("Expected barrier damage   ", attackingPowerSummary.ExpectedBarrierDamage)

	println("---")
	dieRoller := &utility.RandomDieRoller{}
	attackResults := powerusage.UsePowerAgainstSquaddiesAndGetResults(&powerusagecontext.PowerUsageContext{
		SquaddieRepo:      squaddieRepo,
		ActingSquaddieID:  attacker.ID,
		TargetSquaddieIDs: []string{target.ID},
		PowerID:           power.ID,
		PowerRepo:         powerRepo,
	}, dieRoller)
	if !attackResults.AttackingPowerResults[0].WasAHit {
		println("Missed")
	} else if attackResults.AttackingPowerResults[0].WasACriticalHit {
		println("Critical Hit")
		println(attackResults.AttackingPowerResults[0].DamageTaken, "damage taken")
		println(attackResults.AttackingPowerResults[0].BarrierDamage, "barrier damage")
	} else {
		println("Hit")
		println(attackResults.AttackingPowerResults[0].DamageTaken, "damage taken")
		println(attackResults.AttackingPowerResults[0].BarrierDamage, "barrier damage")
	}
}

func loadSquaddieRepo() (repo *squaddie.Repository) {
	squaddieYamlData, err := ioutil.ReadFile("data/squaddieDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}

	squaddieRepo := squaddie.NewSquaddieRepository()
	squaddieRepo.AddYAMLSource(squaddieYamlData)
	return squaddieRepo
}

func loadPowerRepo() (repo *power.Repository) {
	powerYamlData, err := ioutil.ReadFile("data/powerDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}
	powerRepo := power.NewPowerRepository()
	powerRepo.AddYAMLSource(powerYamlData)
	return powerRepo
}

func loadActors (attackerID, targetID, powerID string, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	attacker := squaddieRepo.CloneSquaddieBasedOnSquaddieID(attackerID)
	attacker.SetBarrierToMax()

	target := squaddieRepo.CloneSquaddieBasedOnSquaddieID(targetID)
	target.SetBarrierToMax()

	powerusage.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerReferences, powerRepo)

	power := powerRepo.GetPowerByID(powerID)

	return attacker, target, power
}