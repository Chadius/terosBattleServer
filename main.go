package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/utility"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	squaddieRepo := loadSquaddieRepo()
	powerRepo := loadPowerRepo()
	repos := repositories.RepositoryCollection{
		PowerRepo: powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	attacker, target, power := loadActors(
		"squaddieTeros",
		"squaddieBandit",
		"powerSpear",
		&repos,
	)

	powerSetup := powerusagescenario.Setup{
		UserID:          attacker.Identification.ID,
		PowerID:         power.ID,
		Targets:         []string{target.Identification.ID},
		IsCounterAttack: false,
	}

	powerForecast := &powerattackforecast.Forecast{
		Setup: powerSetup,
		Repositories: &repositories.RepositoryCollection{
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
		println()
	}
}

func printAttackForecast(forecast *powerattackforecast.Calculation) {
	printPartOfAttackForecast(forecast.Attack, forecast.Setup, forecast.Repositories)
	if forecast.CounterAttack != nil {
		println("")
		println("then Counterattack:")
		printPartOfAttackForecast(forecast.CounterAttack, forecast.CounterAttackSetup, forecast.Repositories)
	}
}

func printPartOfAttackForecast(forecast *powerattackforecast.AttackForecast, setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(setup.Targets[0])
	attackingPower := powerRepo.GetPowerByID(setup.PowerID)

	//hitChance := power.GetChanceToHitBasedOnHitRate(forecast.VersusContext.ToHitBonus)
	println(attacker.Identification.Name, "will attack", target.Identification.Name, "with", attackingPower.Name)
	println("Attacker ToHit bonus", forecast.VersusContext.ToHit.ToHitBonus)

	if forecast.VersusContext.NormalDamage.IsFatalToTarget {
		println("will kill if it hits")
	}

	//println("Chance to hit (out of 36) ", hitChance)
	println("Forecasted Damage              ", forecast.VersusContext.NormalDamage.RawDamageDealt)
	//println("Forecasted Barrier damage            ", forecast.VersusContext.NormalDamage.TotalRawBarrierBurnt)
}

func printAttackReport(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(result.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)
	attackingPower := powerRepo.GetPowerByID(result.PowerID)

	println(attacker.Identification.Name, "attacks", target.Identification.Name, "with", attackingPower.Name)

	println(attacker.Identification.Name, "attacks with a", result.Attack.AttackRoll, "+", result.Attack.AttackerToHitBonus, "=", result.Attack.AttackerTotal)
	println(target.Identification.Name, "defends with a", result.Attack.DefendRoll, "+", result.Attack.DefenderToHitPenalty, "=", result.Attack.DefenderTotal)
	if !result.Attack.HitTarget {
		println("Missed")
		return
	}

	if result.Attack.CriticallyHitTarget {
		println("Critical Hit")
	} else {
		println("Hit")
	}
	damageTaken := "  deals " + strconv.Itoa(result.Attack.Damage.RawDamageDealt)
	if result.Attack.Damage.TotalRawBarrierBurnt > 0 {
		damageTaken += " damage, " + strconv.Itoa(result.Attack.Damage.TotalRawBarrierBurnt) + " barrier burn"
	}
	println(damageTaken)

	healthStatus := target.Identification.Name + " HP: " + strconv.Itoa(target.Defense.CurrentHitPoints) + "/" + strconv.Itoa(target.Defense.MaxHitPoints)
	if target.Defense.CurrentBarrier > 0  {
		healthStatus += "Barrier" + strconv.Itoa(target.Defense.CurrentBarrier)
	}
	println(healthStatus)

	if target.Defense.IsDead() {
		println(target.Identification.Name, "falls!")
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

func loadActors (attackerID, targetID, powerID string, repositories *repositories.RepositoryCollection) (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo
	attacker := squaddieRepo.GetOriginalSquaddieByID(attackerID)
	attacker.Defense.SetBarrierToMax()
	powerequip.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerCollection.PowerReferences, repositories)
	powerequip.EquipDefaultPower(attacker, repositories)

	target := squaddieRepo.GetOriginalSquaddieByID(targetID)
	target.Defense.SetBarrierToMax()
	powerequip.LoadAllOfSquaddieInnatePowers(target, target.PowerCollection.PowerReferences, repositories)
	powerequip.EquipDefaultPower(target, repositories)

	power := powerRepo.GetPowerByID(powerID)

	return attacker, target, power
}