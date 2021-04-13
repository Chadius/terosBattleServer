package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
	"github.com/cserrant/terosBattleServer/utility"
	"io/ioutil"
	"log"
)

func main() {
	attacker, target, power := loadActors(
		"squaddieTeros",
		"squaddieBandit",
		"powerSpear",
	)

	powerSummary := powerusage.GetPowerSummary(power, attacker, []*squaddie.Squaddie{target})
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
	attackResults := powerusage.UsePowerAgainstSquaddiesAndGetResults(nil, power, attacker, []*squaddie.Squaddie{target}, dieRoller)
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

func loadActors (attackerID, targetID, powerID string) (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	squaddieYamlData, err := ioutil.ReadFile("data/squaddieDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}

	squaddieRepo := squaddie.NewSquaddieRepository()
	squaddieRepo.AddYAMLSource(squaddieYamlData)
	attacker := squaddieRepo.CloneSquaddieBasedOnSquaddieID(attackerID)
	attacker.SetBarrierToMax()

	target := squaddieRepo.CloneSquaddieBasedOnSquaddieID(targetID)
	target.SetBarrierToMax()

	powerYamlData, err := ioutil.ReadFile("data/powerDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}
	powerRepo := power.NewPowerRepository()
	powerRepo.AddYAMLSource(powerYamlData)

	powerusage.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerReferences, powerRepo)

	power := powerRepo.GetPowerByID(powerID)

	return attacker, target, power
}