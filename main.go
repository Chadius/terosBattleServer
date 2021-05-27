package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
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

	powerSetup := powerusagescenario.Setup{
		UserID:          attacker.Identification.ID,
		PowerID:         power.ID,
		Targets:         []string{target.Identification.ID},
		IsCounterAttack: false,
	}

	powerForecast := &powerattackforecast.Forecast{
		Setup: powerSetup,
		Repositories: &powerusagescenario.RepositoryCollection{
			SquaddieRepo:    squaddieRepo,
			PowerRepo:       powerRepo,
		},
	}
	powerForecast.CalculateForecast()

	for _, forecast := range powerForecast.ForecastedResultPerTarget {
		printAttackForecast(&forecast)
	}

	println("---")
	powerResult := &powercommit.Result{
		Forecast: powerForecast,
		DieRoller: &utility.RandomDieRoller{},
	}
	powerResult.Commit()

	for _, attackReport := range powerResult.ResultPerTarget {
		printAttackReport(attackReport, powerForecast.Repositories)
	}
}

func printAttackForecast(forecast *powerattackforecast.Calculation) {
	printPartOfAttackForecast(forecast.Attack, forecast.Setup, forecast.Repositories)
	if forecast.CounterAttack != nil {
		println("")
		println("Counterattack:")
		printPartOfAttackForecast(forecast.CounterAttack, forecast.CounterAttackSetup, forecast.Repositories)
	}
}

func printPartOfAttackForecast(forecast *powerattackforecast.AttackForecast, setup *powerusagescenario.Setup, repositories *powerusagescenario.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(setup.Targets[0])
	attackingPower := powerRepo.GetPowerByID(setup.PowerID)

	hitChance := power.GetChanceToHitBasedOnHitRate(forecast.VersusContext.ToHitBonus)
	println(attacker.Identification.Name, "will attack", target.Identification.Name, "with", attackingPower.Name)
	println("Attacker ToHit bonus", forecast.VersusContext.ToHitBonus)
	println("Chance to hit (out of 36) ", hitChance)
	println("Damage taken              ", forecast.VersusContext.NormalDamage.DamageDealt)
	println("Barrier damage            ", forecast.VersusContext.NormalDamage.TotalBarrierBurnt)
	println("---")
	println("Expected damage (36 = 1HP)", forecast.VersusContext.NormalDamage.DamageDealt * hitChance)
	println("Expected barrier damage   ", forecast.VersusContext.NormalDamage.TotalBarrierBurnt * hitChance)
}

func printAttackReport(result *powercommit.ResultPerTarget, repositories *powerusagescenario.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(result.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)
	attackingPower := powerRepo.GetPowerByID(result.PowerID)

	println(attacker.Identification.Name, "attacks with a", result.AttackRoll)
	println(target.Identification.Name, "defends with a", result.DefendRoll)
	if !result.Attack.HitTarget {
		println("Missed")
		return
	}

	if result.Attack.CriticallyHitTarget {
		println("Critical Hit")
	} else {
		println("Hit")
	}
	println(attackingPower.Name, "deals")
	println(result.Attack.Damage.DamageDealt, "damage taken")
	println(result.Attack.Damage.TotalBarrierBurnt, "barrier damage")

	println(target.Identification.Name, "HP:", target.Defense.CurrentHitPoints,"/",target.Defense.MaxHitPoints,"Barrier",target.Defense.CurrentBarrier,"/",target.Defense.MaxBarrier)
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
	attacker := squaddieRepo.GetOriginalSquaddieByID(attackerID)
	attacker.Defense.SetBarrierToMax()
	powerequip.EquipDefaultPower(attacker, powerRepo)
	powerequip.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerCollection.PowerReferences, powerRepo)

	target := squaddieRepo.GetOriginalSquaddieByID(targetID)
	target.Defense.SetBarrierToMax()
	powerequip.EquipDefaultPower(target, powerRepo)
	powerequip.LoadAllOfSquaddieInnatePowers(target, target.PowerCollection.PowerReferences, powerRepo)

	power := powerRepo.GetPowerByID(powerID)

	return attacker, target, power
}