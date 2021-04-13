package powerusage_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CalculateExpectedDamageFromAttackSuite struct {
	teros *squaddie.Squaddie
	bandit *squaddie.Squaddie
	bandit2 *squaddie.Squaddie
	spear *power.Power
	blot  *power.Power
}

var _ = Suite(&CalculateExpectedDamageFromAttackSuite{})

func (suite *CalculateExpectedDamageFromAttackSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("suite.teros")
	suite.teros.Name = "suite.teros"

	suite.spear = power.NewPower("suite.spear")
	suite.spear.PowerType = power.Physical

	suite.blot = power.NewPower("suite.blot")
	suite.blot.PowerType = power.Spell

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Name = "bandit"

	suite.bandit2 = squaddie.NewSquaddie("bandit2")
	suite.bandit2.Name = "bandit2"
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCalculateAttackerHitBonus(checker *C) {
	suite.teros.Aim = 2
	suite.blot.AttackEffect.ToHitBonus = 1

	totalToHitBonus := powerusage.GetPowerToHitBonusWhenUsedBySquaddie(suite.blot, suite.teros, false)
	checker.Assert(totalToHitBonus, Equals, 3)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCalculateAttackerHitBonusWhenCounterAttacking(checker *C) {
	suite.teros.Aim = 2
	suite.blot.AttackEffect.ToHitBonus = 1
	suite.blot.AttackEffect.CounterAttackToHitPenalty = -2

	totalToHitBonus := powerusage.GetPowerToHitBonusWhenUsedBySquaddie(suite.blot, suite.teros, true)
	checker.Assert(totalToHitBonus, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestPhysicalDamage(checker *C) {
	suite.teros.Strength = 2
	suite.teros.Mind = 3

	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.DamageBonus = 2

	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 6

	totalDamageBonus := powerusage.GetPowerDamageBonusWhenUsedBySquaddie(suite.spear, suite.teros)
	checker.Assert(totalDamageBonus, Equals, 4)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSpellDamage(checker *C) {
	suite.teros.Strength = 2
	suite.teros.Mind = 3

	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.DamageBonus = 2

	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 6

	totalDamageBonus := powerusage.GetPowerDamageBonusWhenUsedBySquaddie(suite.blot, suite.teros)
	checker.Assert(totalDamageBonus, Equals, 9)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCriticalPhysicalDamage(checker *C) {
	suite.teros.Strength = 2
	suite.teros.Mind = 3

	suite.spear.PowerType = power.Physical
	suite.spear.AttackEffect.DamageBonus = 2

	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 6

	totalDamageBonus := powerusage.GetPowerCriticalDamageBonusWhenUsedBySquaddie(suite.spear, suite.teros)
	checker.Assert(totalDamageBonus, Equals, 8)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestToHitReductionAgainstPhysical(checker *C) {
	suite.teros.Dodge = 2
	suite.teros.Deflect = 9001

	suite.spear.PowerType = power.Physical

	suite.blot.PowerType = power.Spell

	toHitPenalty := powerusage.GetPowerToHitPenaltyAgainstSquaddie(suite.spear, suite.teros)
	checker.Assert(toHitPenalty, Equals, 2)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestToHitReductionAgainstSpell(checker *C) {
	suite.teros.Dodge = 2
	suite.teros.Deflect = 9001

	suite.spear.PowerType = power.Physical

	suite.blot.PowerType = power.Spell

	toHitPenalty := powerusage.GetPowerToHitPenaltyAgainstSquaddie(suite.blot, suite.teros)
	checker.Assert(toHitPenalty, Equals, 9001)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestFullPhysicalDamageAgainstUnarmored(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 4)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSomePhysicalDamageAgainstSomeArmor(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.Armor = 3
	totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSomePhysicalDamageAgainstSomeBarrier(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.MaxBarrier = 4
	suite.bandit.CurrentBarrier = 1
	totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 3)
	checker.Assert(initialBarrierDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestNoPhysicalDamageAgainstStrongBarrier(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.MaxBarrier = 4
	suite.bandit.CurrentBarrier = 4
	totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 4)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestNoPhysicalDamageAgainstStrongArmor(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3
	suite.bandit.Armor = 4
	totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(suite.spear, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 0)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestFullSpellDamageAgainstUnarmored(checker *C) {
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestFullSpellDamageAgainstNoBarrier(checker *C) {
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.Armor = 9001
	totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestBarrierAbsorbsDamageBeforeHealth(checker *C) {
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.MaxBarrier = 4
	suite.bandit.CurrentBarrier = 1
	totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 5)
	checker.Assert(initialBarrierDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestNoSpellDamageAgainstStrongBarrier(checker *C) {
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.MaxBarrier = 9001
	suite.bandit.CurrentBarrier = 9001
	totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestPowerDealsExtraBarrierDamage(checker *C) {
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.MaxBarrier = 8
	suite.bandit.CurrentBarrier = 8
	suite.blot.AttackEffect.ExtraBarrierDamage = 2

	totalHealthDamage, initialBarrierDamage, extraBarrierDamage := powerusage.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 6)
	checker.Assert(extraBarrierDamage, Equals, 2)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryKnowsExtraBarrierDamageIsCappedIfBarrierIsDestroyed(checker *C) {
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	suite.bandit.MaxBarrier = 8
	suite.bandit.CurrentBarrier = 7
	suite.blot.AttackEffect.ExtraBarrierDamage = 2

	totalHealthDamage, initialBarrierDamage, extraBarrierDamage := powerusage.GetHowTargetDistributesDamage(suite.blot, suite.teros, suite.bandit)
	checker.Assert(totalHealthDamage, Equals, 0)
	checker.Assert(initialBarrierDamage, Equals, 6)
	checker.Assert(extraBarrierDamage, Equals, 1)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestPhysicalPowerSummary(checker *C) {
	suite.bandit.Armor = 1
	suite.bandit.Dodge = 1
	suite.bandit.MaxBarrier = 4
	suite.bandit.CurrentBarrier = 1

	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4

	attackingPowerSummary := powerusage.GetExpectedDamage(nil, &powerusagecontext.AttackContext{
		Power:           suite.spear,
		Attacker:        suite.teros,
		Target:          suite.bandit,
		IsCounterAttack: false,
	})
	checker.Assert(attackingPowerSummary.AttackingSquaddieID, Equals, suite.teros.ID)
	checker.Assert(attackingPowerSummary.PowerID, Equals, suite.spear.ID)
	checker.Assert(attackingPowerSummary.TargetSquaddieID, Equals, suite.bandit.ID)
	checker.Assert(attackingPowerSummary.IsACounterAttack, Equals, false)
	checker.Assert(attackingPowerSummary.ChanceToHit, Equals, 15)
	checker.Assert(attackingPowerSummary.DamageTaken, Equals, 2)
	checker.Assert(attackingPowerSummary.ExpectedDamage, Equals, 30)
	checker.Assert(attackingPowerSummary.BarrierDamageTaken, Equals, 1)
	checker.Assert(attackingPowerSummary.ExpectedBarrierDamage, Equals, 15)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryWithBarrierBurn(checker *C) {
	suite.bandit.Armor = 1
	suite.bandit.Dodge = 1
	suite.bandit.MaxBarrier = 10
	suite.bandit.CurrentBarrier = 10

	suite.teros.Aim = 3
	suite.teros.Mind = 2
	suite.blot.AttackEffect.DamageBonus = 4
	suite.blot.AttackEffect.ExtraBarrierDamage = 3
	attackingPowerSummary := powerusage.GetExpectedDamage(nil, &powerusagecontext.AttackContext{
		Power:           suite.blot,
		Attacker:        suite.teros,
		Target:          suite.bandit,
		IsCounterAttack: false,
	})
	checker.Assert(attackingPowerSummary.ChanceToHit, Equals, 33)
	checker.Assert(attackingPowerSummary.DamageTaken, Equals, 0)
	checker.Assert(attackingPowerSummary.ExpectedDamage, Equals, 0)
	checker.Assert(attackingPowerSummary.BarrierDamageTaken, Equals, 9)
	checker.Assert(attackingPowerSummary.ExpectedBarrierDamage, Equals, 9 * 33)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryPerTarget(checker *C) {
	powerSummary := powerusage.GetPowerSummary(suite.spear, suite.teros, []*squaddie.Squaddie{suite.bandit, suite.bandit2})
	checker.Assert(powerSummary.UserSquaddieID, Equals, suite.teros.ID)
	checker.Assert(powerSummary.PowerID, Equals, suite.spear.ID)
	checker.Assert(powerSummary.AttackEffectSummary, HasLen, 2)
	checker.Assert(powerSummary.AttackEffectSummary[0].TargetSquaddieID, Equals, suite.bandit.ID)
	checker.Assert(powerSummary.AttackEffectSummary[1].TargetSquaddieID, Equals, suite.bandit2.ID)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestChanceToCriticalHitOnTheSummary(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.spear.AttackEffect.CriticalHitThreshold = 4
	attackingPowerSummary := powerusage.GetExpectedDamage(nil, &powerusagecontext.AttackContext{
		Power:           suite.spear,
		Attacker:        suite.teros,
		Target:          suite.bandit,
		IsCounterAttack: false,
	})
	checker.Assert(attackingPowerSummary.ChanceToCritical, Equals, 6)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestCriticalHitDoublesDamageBeforeArmorAndBarrier(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.bandit.Armor = 1
	suite.bandit.MaxBarrier = 4
	suite.bandit.CurrentBarrier = 4
	suite.spear.AttackEffect.CriticalHitThreshold = 4
	attackingPowerSummary := powerusage.GetExpectedDamage(nil, &powerusagecontext.AttackContext{
		Power:           suite.spear,
		Attacker:        suite.teros,
		Target:          suite.bandit,
		IsCounterAttack: false,
	})
	checker.Assert(attackingPowerSummary.CriticalDamageTaken, Equals, 3)
	checker.Assert(attackingPowerSummary.CriticalBarrierDamageTaken, Equals, 4)
	checker.Assert(attackingPowerSummary.CriticalExpectedDamage, Equals, 3 * 21)
	checker.Assert(attackingPowerSummary.CriticalExpectedBarrierDamage, Equals, 4 * 21)
}

func (suite *CalculateExpectedDamageFromAttackSuite) TestSummaryIgnoresCriticalIfAttackCannotCritical(checker *C) {
	suite.teros.Strength = 1
	suite.spear.AttackEffect.DamageBonus = 3

	suite.spear.AttackEffect.CriticalHitThreshold = 0
	attackingPowerSummary := powerusage.GetExpectedDamage(nil, &powerusagecontext.AttackContext{
		Power:           suite.spear,
		Attacker:        suite.teros,
		Target:          suite.bandit,
		IsCounterAttack: false,
	})
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
	suite.teros.Name = "teros"
}

func (suite *SquaddieGainsPowerSuite) TestGiveSquaddieInnatePowersWithRepository(checker *C) {
	temporaryPowerReferences := []*power.Reference{{Name: "suite.spear", ID: suite.spear.ID}}
	numberOfPowersAdded, err := powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, temporaryPowerReferences, suite.powerRepository)
	checker.Assert(numberOfPowersAdded, Equals, 1)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.spear.ID)
}

func (suite *SquaddieGainsPowerSuite) TestStopAddingNonexistentPowers(checker *C) {
	scimitar := power.NewPower("Scimitar")
	scimitar.PowerType = power.Physical

	temporaryPowerReferences := []*power.Reference{{Name: "Scimitar", ID: scimitar.ID}}
	numberOfPowersAdded, err := powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, temporaryPowerReferences, suite.powerRepository)
	checker.Assert(numberOfPowersAdded, Equals, 0)
	checker.Assert(err.Error(), Equals, "squaddie 'teros' tried to add Power 'Scimitar' but it does not exist")

	attackIDNamePairs := suite.teros.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 0)
}

type CreatePowerReportSuite struct {
	teros *squaddie.Squaddie
	bandit *squaddie.Squaddie
	bandit2 *squaddie.Squaddie
	blot *power.Power
}

var _ = Suite(&CreatePowerReportSuite{})

func (suite *CreatePowerReportSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("suite.teros")
	suite.teros.Name = "suite.teros"
	suite.teros.Mind = 1

	suite.bandit = squaddie.NewSquaddie("suite.bandit")
	suite.bandit.Name = "suite.bandit"

	suite.bandit2 = squaddie.NewSquaddie("suite.bandit")
	suite.bandit2.Name = "suite.bandit"

	suite.blot = power.NewPower("suite.blot")
	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 1
}

func (suite *CreatePowerReportSuite) TestPowerReportWhenMissed(checker *C) {
	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.blot,
		suite.teros,
		[]*squaddie.Squaddie{
			suite.bandit,
		},
		dieRoller,
	)
	checker.Assert(powerResult.AttackerID, Equals, suite.teros.ID)
	checker.Assert(powerResult.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerResult.AttackingPowerResults, HasLen, 1)
	checker.Assert(powerResult.AttackingPowerResults[0].WasAHit, Equals, false)
}

func (suite *CreatePowerReportSuite) TestPowerReportWhenHitButNoCrit(checker *C) {
	dieRoller := &testutility.AlwaysHitDieRoller{}

	powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.blot,
		suite.teros,
		[]*squaddie.Squaddie{
			suite.bandit,
		},
		dieRoller,
	)
	checker.Assert(powerResult.AttackerID, Equals, suite.teros.ID)
	checker.Assert(powerResult.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerResult.AttackingPowerResults, HasLen, 1)
	checker.Assert(powerResult.AttackingPowerResults[0].WasAHit, Equals, true)
	checker.Assert(powerResult.AttackingPowerResults[0].WasACriticalHit, Equals, false)
	checker.Assert(powerResult.AttackingPowerResults[0].DamageTaken, Equals, 2)
	checker.Assert(powerResult.AttackingPowerResults[0].BarrierDamage, Equals, 0)
}

func (suite *CreatePowerReportSuite) TestPowerReportWhenCrits(checker *C) {
	dieRoller := &testutility.AlwaysHitDieRoller{}
	suite.blot.AttackEffect.CriticalHitThreshold = 900

	powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.blot,
		suite.teros,
		[]*squaddie.Squaddie{
			suite.bandit,
		},
		dieRoller,
	)
	checker.Assert(powerResult.AttackerID, Equals, suite.teros.ID)
	checker.Assert(powerResult.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerResult.AttackingPowerResults, HasLen, 1)
	checker.Assert(powerResult.AttackingPowerResults[0].WasAHit, Equals, true)
	checker.Assert(powerResult.AttackingPowerResults[0].WasACriticalHit, Equals, true)
	checker.Assert(powerResult.AttackingPowerResults[0].DamageTaken, Equals, 4)
	checker.Assert(powerResult.AttackingPowerResults[0].BarrierDamage, Equals, 0)
}

func (suite *CreatePowerReportSuite) TestReportPerTarget(checker *C) {
	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.blot,
		suite.teros,
		[]*squaddie.Squaddie{
			suite.bandit,
			suite.bandit2,
		},
		dieRoller,
	)
	checker.Assert(powerResult.AttackerID, Equals, suite.teros.ID)
	checker.Assert(powerResult.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerResult.AttackingPowerResults, HasLen, 2)
	checker.Assert(powerResult.AttackingPowerResults[0].TargetID, Equals, suite.bandit.ID)
	checker.Assert(powerResult.AttackingPowerResults[1].TargetID, Equals, suite.bandit2.ID)
}

type SquaddieCommitToPowerUsageSuite struct {
	teros *squaddie.Squaddie
	spear *power.Power
	scimitar *power.Power
	powerRepo *power.Repository
	bandit *squaddie.Squaddie
	blot *power.Power
	squaddieRepo *squaddie.Repository
}

var _ = Suite(&SquaddieCommitToPowerUsageSuite{})

func (suite *SquaddieCommitToPowerUsageSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("suite.teros")
	suite.spear = power.NewPower("suite.spear")
	suite.spear.AttackEffect.CanBeEquipped = true

	suite.scimitar = power.NewPower("scimitar the second")
	suite.scimitar.AttackEffect.CanBeEquipped = true

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.spear,
		suite.scimitar,
	})

	suite.bandit = squaddie.NewSquaddie("suite.bandit")
	suite.bandit.Name = "suite.bandit"

	suite.blot = power.NewPower("suite.blot")
	suite.blot.PowerType = power.Spell

	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
	})
}

func (suite *SquaddieCommitToPowerUsageSuite) TestSquaddiesEquipPowerUponCommit(checker *C) {
	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerReport := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.scimitar,
		suite.teros,
		[]*squaddie.Squaddie{
			suite.bandit,
		},
		dieRoller,
	)

	powerusage.CommitPowerUse(powerReport, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.scimitar.ID)
}

func (suite *SquaddieCommitToPowerUsageSuite) TestSquaddieWillKeepPreviousPowerIfCommitPowerIsUnequippable(checker *C) {
	powerusage.SquaddieEquipPower(suite.teros, suite.scimitar.ID, suite.powerRepo)

	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerReport := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.blot,
		suite.teros,
		[]*squaddie.Squaddie{
			suite.bandit,
		},
		dieRoller,
	)

	powerusage.CommitPowerUse(powerReport, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(powerusage.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.scimitar.ID)
}

func (suite *SquaddieCommitToPowerUsageSuite) TestSquaddieWillNotEquipPowerIfNoneExistAfterCommitting(checker *C) {
	mysticMage := squaddie.NewSquaddie("Mystic Mage")
	mysticMagePowerReferences := []*power.Reference{
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(mysticMage, mysticMagePowerReferences, suite.powerRepo)

	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		mysticMage,
	})

	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerReport := powerusage.UsePowerAgainstSquaddiesAndGetResults(
		nil,
		suite.blot,
		mysticMage,
		[]*squaddie.Squaddie{
			suite.bandit,
		},
		dieRoller,
	)

	powerusage.CommitPowerUse(powerReport, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(powerusage.GetEquippedPower(mysticMage, suite.powerRepo), IsNil)
}

type TargetAttemptsCounterSuite struct {
	teros *squaddie.Squaddie
	spear *power.Power
	blot *power.Power

	bandit *squaddie.Squaddie
	axe *power.Power

	powerRepo *power.Repository
	squaddieRepo *squaddie.Repository
}

var _ = Suite(&TargetAttemptsCounterSuite{})

func (suite *TargetAttemptsCounterSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("suite.teros")
	suite.spear = power.NewPower("suite.spear")
	suite.spear.AttackEffect.CanBeEquipped = true
	suite.spear.AttackEffect.CanCounterAttack = true
	suite.spear.AttackEffect.CounterAttackToHitPenalty = -2

	suite.axe = power.NewPower("axe the second")
	suite.axe.AttackEffect.CanBeEquipped = true
	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.CounterAttackToHitPenalty = -2

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.spear,
		suite.axe,
	})

	suite.blot = power.NewPower("suite.blot")
	suite.blot.PowerType = power.Spell

	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.blot.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)

	suite.bandit = squaddie.NewSquaddie("suite.bandit")
	suite.bandit.Name = "suite.bandit"
	banditPowerReferences := []*power.Reference{
		suite.axe.GetReference(),
	}
	powerusage.LoadAllOfSquaddieInnatePowers(suite.bandit, banditPowerReferences, suite.powerRepo)

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
	})
}

func (suite *TargetAttemptsCounterSuite) TestTargetWillCounterAttackWithEquippedCounterablePower(checker *C) {
	powerusage.SquaddieEquipPower(suite.teros, suite.spear.ID, suite.powerRepo)
	powerusage.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.powerRepo)

	expectedTerosCounterAttackSummary := powerusage.GetExpectedDamage(
		nil,
		&powerusagecontext.AttackContext{
		Power:				suite.spear,
		Attacker:			suite.teros,
		Target:				suite.bandit,
		IsCounterAttack:	false,
		PowerRepo:			suite.powerRepo,
	})
	terosHitRate := expectedTerosCounterAttackSummary.HitRate

	banditAttackSummary := powerusage.GetExpectedDamage(
		nil,
		&powerusagecontext.AttackContext{
		Power:				suite.axe,
		Attacker:			suite.bandit,
		Target:				suite.teros,
		IsCounterAttack:	false,
		PowerRepo:			suite.powerRepo,
	})
	checker.Assert(banditAttackSummary.CounterAttack, NotNil)
	checker.Assert(banditAttackSummary.CounterAttack.IsACounterAttack, Equals, true)
	checker.Assert(banditAttackSummary.CounterAttack.AttackingSquaddieID, Equals, suite.teros.ID)
	checker.Assert(banditAttackSummary.CounterAttack.PowerID, Equals, suite.spear.ID)
	checker.Assert(banditAttackSummary.CounterAttack.TargetSquaddieID, Equals, suite.bandit.ID)
	checker.Assert(banditAttackSummary.CounterAttack.HitRate, Equals, terosHitRate - 2)
}
