package powerforecast_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/powerforecast"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CalculateExpectedDamageFromAttackSuite struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	bandit2			*squaddie.Squaddie
	spear			*power.Power
	blot			*power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository
}

var _ = Suite(&CalculateExpectedDamageFromAttackSuite{})

func (suite *CalculateExpectedDamageFromAttackSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	suite.bandit2 = squaddie.NewSquaddie("bandit2")
	suite.bandit2.Identification.Name = "bandit2"

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.bandit2})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCalculateAttackerHitBonus(checker *C) {
	suite.teros.Offense.Aim = 2
	suite.blot.AttackEffect.ToHitBonus = 1

	totalToHitBonus := powerforecast.GetPowerToHitBonusWhenUsedBySquaddie(suite.blot, suite.teros, false)
	checker.Assert(totalToHitBonus, Equals, 3)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCalculateAttackerHitBonusWhenCounterAttacking(checker *C) {
	suite.teros.Offense.Aim = 2
	suite.blot.AttackEffect.ToHitBonus = 1
	suite.blot.AttackEffect.CounterAttackToHitPenalty = -2

	totalToHitBonus := powerforecast.GetPowerToHitBonusWhenUsedBySquaddie(suite.blot, suite.teros, true)
	checker.Assert(totalToHitBonus, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestPhysicalDamage(checker *C) {
	suite.teros.Offense.Strength = 2
	suite.teros.Offense.Mind = 3

	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.DamageBonus = 2

	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 6

	totalDamageBonus := powerforecast.GetPowerDamageBonusWhenUsedBySquaddie(suite.spear, suite.teros)
	checker.Assert(totalDamageBonus, Equals, 4)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSpellDamage(checker *C) {
	suite.teros.Offense.Strength = 2
	suite.teros.Offense.Mind = 3

	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.DamageBonus = 2

	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 6

	totalDamageBonus := powerforecast.GetPowerDamageBonusWhenUsedBySquaddie(suite.blot, suite.teros)
	checker.Assert(totalDamageBonus, Equals, 9)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCriticalPhysicalDamage(checker *C) {
	suite.teros.Offense.Strength = 2
	suite.teros.Offense.Mind = 3

	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.DamageBonus = 2

	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 6

	totalDamageBonus := powerforecast.GetPowerCriticalDamageBonusWhenUsedBySquaddie(suite.spear, suite.teros)
	checker.Assert(totalDamageBonus, Equals, 8)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestToHitReductionAgainstPhysical(checker *C) {
	suite.teros.Defense.Dodge = 2
	suite.teros.Defense.Deflect = 9001

	suite.spear.PowerType = power.Physical

	suite.blot.PowerType = power.Spell

	toHitPenalty := powerforecast.GetPowerToHitPenaltyAgainstSquaddie(suite.spear, suite.teros)
	checker.Assert(toHitPenalty, Equals, 2)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestToHitReductionAgainstSpell(checker *C) {
	suite.teros.Defense.Dodge = 2
	suite.teros.Defense.Deflect = 9001

	suite.spear.PowerType = power.Physical

	suite.blot.PowerType = power.Spell

	toHitPenalty := powerforecast.GetPowerToHitPenaltyAgainstSquaddie(suite.blot, suite.teros)
	checker.Assert(toHitPenalty, Equals, 9001)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestFullPhysicalDamageAgainstUnarmored(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	totalHealthDamage, _, _ := powerforecast.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 4)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSomePhysicalDamageAgainstSomeArmor(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.Defense.Armor = 3
	totalHealthDamage, _, _ := powerforecast.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSomePhysicalDamageAgainstSomeBarrier(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.Defense.MaxBarrier = 4
	suite.bandit.Defense.CurrentBarrier = 1
	totalHealthDamage, initialBarrierDamage, _ := powerforecast.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 3)
	checker.Assert(initialBarrierDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestNoPhysicalDamageAgainstStrongBarrier(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.Defense.MaxBarrier = 4
	suite.bandit.Defense.CurrentBarrier = 4
	totalHealthDamage, initialBarrierDamage, _ := powerforecast.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 4)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestNoPhysicalDamageAgainstStrongArmor(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.Defense.Armor = 4
	totalHealthDamage, initialBarrierDamage, _ := powerforecast.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 0)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestFullSpellDamageAgainstUnarmored(checker *C) {
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	totalHealthDamage, _, _ := powerforecast.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestFullSpellDamageAgainstNoBarrier(checker *C) {
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.Defense.Armor = 9001
	totalHealthDamage, _, _ := powerforecast.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestBarrierAbsorbsDamageBeforeHealth(checker *C) {
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.Defense.MaxBarrier = 4
	suite.bandit.Defense.CurrentBarrier = 1
	totalHealthDamage, initialBarrierDamage, _ := powerforecast.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 5)
	checker.Assert(initialBarrierDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestNoSpellDamageAgainstStrongBarrier(checker *C) {
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.Defense.MaxBarrier = 9001
	suite.bandit.Defense.CurrentBarrier = 9001
	totalHealthDamage, initialBarrierDamage, _ := powerforecast.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestPowerDealsExtraBarrierDamage(checker *C) {
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.Defense.MaxBarrier = 8
	suite.bandit.Defense.CurrentBarrier = 8
	suite.blot.AttackEffect.ExtraBarrierBurn = 2

	totalHealthDamage, initialBarrierDamage, extraBarrierDamage := powerforecast.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 6)
	checker.Assert(extraBarrierDamage, Equals, 2)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryKnowsExtraBarrierDamageIsCappedIfBarrierIsDestroyed(checker *C) {
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.Defense.MaxBarrier = 8
	suite.bandit.Defense.CurrentBarrier = 7
	suite.blot.AttackEffect.ExtraBarrierBurn = 2

	totalHealthDamage, initialBarrierDamage, extraBarrierDamage := powerforecast.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 6)
	checker.Assert(extraBarrierDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestPhysicalPowerSummary(checker *C) {
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.Dodge = 1
	suite.bandit.Defense.MaxBarrier = 4
	suite.bandit.Defense.CurrentBarrier = 1

	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	attackingPowerSummary := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.spear.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:           suite.spear.ID,
			AttackerID:        suite.teros.Identification.ID,
			TargetID:          suite.bandit.Identification.ID,
			IsCounterAttack: false,
		},
	)
	checker.Assert(attackingPowerSummary.AttackingSquaddieID, Equals, suite.teros.Identification.ID)
	checker.Assert(attackingPowerSummary.PowerID, Equals, suite.spear.ID)
	checker.Assert(attackingPowerSummary.TargetSquaddieID, Equals, suite.bandit.Identification.ID)
	checker.Assert(attackingPowerSummary.IsACounterAttack, Equals, false)
	checker.Assert(attackingPowerSummary.ChanceToHit, Equals, 15)
	checker.Assert(attackingPowerSummary.DamageTaken, Equals, 2)
	checker.Assert(attackingPowerSummary.ExpectedDamage, Equals, 30)
	checker.Assert(attackingPowerSummary.BarrierDamageTaken, Equals, 1)
	checker.Assert(attackingPowerSummary.ExpectedBarrierDamage, Equals, 15)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryWithBarrierBurn(checker *C) {
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.Dodge = 1
	suite.bandit.Defense.MaxBarrier = 10
	suite.bandit.Defense.CurrentBarrier = 10

	suite.teros.Offense.Aim = 3
	suite.teros.Offense.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4
	suite.blot.AttackEffect.ExtraBarrierBurn = 3
	attackingPowerSummary := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:           suite.blot.ID,
			AttackerID:        suite.teros.Identification.ID,
			TargetID:          suite.bandit.Identification.ID,
			IsCounterAttack: false,
		},
	)
	checker.Assert(attackingPowerSummary.ChanceToHit, Equals, 33)
	checker.Assert(attackingPowerSummary.DamageTaken, Equals, 0)
	checker.Assert(attackingPowerSummary.ExpectedDamage, Equals, 0)
	checker.Assert(attackingPowerSummary.BarrierDamageTaken, Equals, 9)
	checker.Assert(attackingPowerSummary.ExpectedBarrierDamage, Equals, 9 * 33)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryPerTarget(checker *C) {
	powerSummary := powerforecast.CalculatePowerForecast(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID, suite.bandit2.Identification.ID},
			PowerID:           suite.spear.ID,
			PowerRepo:         suite.powerRepo,
		},
	)
	checker.Assert(powerSummary.UserSquaddieID, Equals, suite.teros.Identification.ID)
	checker.Assert(powerSummary.PowerID, Equals, suite.spear.ID)
	checker.Assert(powerSummary.AttackPowerForecast, HasLen, 2)
	checker.Assert(powerSummary.AttackPowerForecast[0].TargetSquaddieID, Equals, suite.bandit.Identification.ID)
	checker.Assert(powerSummary.AttackPowerForecast[1].TargetSquaddieID, Equals, suite.bandit2.Identification.ID)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestChanceToCriticalHitOnTheSummary(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.spear.AttackEffect.CriticalHitThreshold = 4
	attackingPowerSummary := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.spear.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:           suite.spear.ID,
			AttackerID:        suite.teros.Identification.ID,
			TargetID:          suite.bandit.Identification.ID,
			IsCounterAttack: false,
		},
	)
	checker.Assert(attackingPowerSummary.ChanceToCritical, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCriticalHitDoublesDamageBeforeArmorAndBarrier(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.MaxBarrier = 4
	suite.bandit.Defense.CurrentBarrier = 4
	suite.spear.AttackEffect.CriticalHitThreshold = 4
	attackingPowerSummary := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.spear.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:           suite.spear.ID,
			AttackerID:        suite.teros.Identification.ID,
			TargetID:          suite.bandit.Identification.ID,
			IsCounterAttack: false,
		},
	)
	checker.Assert(attackingPowerSummary.CriticalDamageTaken, Equals, 3)
	checker.Assert(attackingPowerSummary.CriticalBarrierDamageTaken, Equals, 4)
	checker.Assert(attackingPowerSummary.CriticalExpectedDamage, Equals, 3 * 21)
	checker.Assert(attackingPowerSummary.CriticalExpectedBarrierDamage, Equals, 4 * 21)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryIgnoresCriticalIfAttackCannotCritical(checker *C) {
	suite.teros.Offense.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.spear.AttackEffect.CriticalHitThreshold = 0
	attackingPowerSummary := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.spear.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:           suite.spear.ID,
			AttackerID:        suite.teros.Identification.ID,
			TargetID:          suite.bandit.Identification.ID,
			IsCounterAttack: false,
		},
	)
	checker.Assert(attackingPowerSummary.ChanceToCritical, Equals, 0)
	checker.Assert(attackingPowerSummary.CriticalDamageTaken, Equals, 0)
	checker.Assert(attackingPowerSummary.CriticalBarrierDamageTaken, Equals, 0)
	checker.Assert(attackingPowerSummary.CriticalExpectedDamage, Equals, 0)
	checker.Assert(attackingPowerSummary.CriticalExpectedBarrierDamage, Equals, 0)
}

type SquaddieGainsPowerSuite struct {
	teros *squaddie.Squaddie
	powerRepository *power.Repository
	spear *power.Power
}

var _ = Suite(&SquaddieGainsPowerSuite{})

func (suite *SquaddieGainsPowerSuite) SetUpTest(checker *C) {
	suite.powerRepository = power.NewPowerRepository()

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.ID = "deadbeef"
	newPowers := []*power.Power{suite.spear}
	suite.powerRepository.AddSlicePowerSource(newPowers)

	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
}

func (suite *SquaddieGainsPowerSuite) TestGiveSquaddieInnatePowersWithRepository(checker *C) {
	temporaryPowerReferences := []*power.Reference{{Name: "spear", ID: suite.spear.ID}}
	numberOfPowersAdded, err := powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, temporaryPowerReferences, suite.powerRepository)
	checker.Assert(numberOfPowersAdded, Equals, 1)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.spear.ID)
}

func (suite *SquaddieGainsPowerSuite) TestStopAddingNonexistentPowers(checker *C) {
	scimitar := power.NewPower("Scimitar")
	scimitar.PowerType = power.Physical

	temporaryPowerReferences := []*power.Reference{{Name: "Scimitar", ID: scimitar.ID}}
	numberOfPowersAdded, err := powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, temporaryPowerReferences, suite.powerRepository)
	checker.Assert(numberOfPowersAdded, Equals, 0)
	checker.Assert(err.Error(), Equals, "squaddie 'teros' tried to add Power 'Scimitar' but it does not exist")

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 0)
}
