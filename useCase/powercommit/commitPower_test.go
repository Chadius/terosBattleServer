package powercommit_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
)

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
	suite.teros = squaddie.NewSquaddie("teros")
	suite.spear = power.NewPower("spear")
	suite.spear.AttackEffect.CanBeEquipped = true

	suite.scimitar = power.NewPower("scimitar the second")
	suite.scimitar.AttackEffect.CanBeEquipped = true

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.spear,
		suite.scimitar,
		suite.blot,
	})

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
	})
}

func (suite *SquaddieCommitToPowerUsageSuite) TestSquaddiesEquipPowerUponCommit(checker *C) {
	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerReport := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.scimitar.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)

	powercommit.CommitPowerUse(powerReport, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(powerequip.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.scimitar.ID)
}

func (suite *SquaddieCommitToPowerUsageSuite) TestSquaddieWillKeepPreviousPowerIfCommitPowerIsUnequippable(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.scimitar.ID, suite.powerRepo)

	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerReport := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)

	powercommit.CommitPowerUse(powerReport, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(powerequip.GetEquippedPower(suite.teros, suite.powerRepo).ID, Equals, suite.scimitar.ID)
}

func (suite *SquaddieCommitToPowerUsageSuite) TestSquaddieWillNotEquipPowerIfNoneExistAfterCommitting(checker *C) {
	mysticMage := squaddie.NewSquaddie("Mystic Mage")
	mysticMagePowerReferences := []*power.Reference{
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(mysticMage, mysticMagePowerReferences, suite.powerRepo)

	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		mysticMage,
	})

	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerReport := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  mysticMage.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)

	powercommit.CommitPowerUse(powerReport, suite.squaddieRepo, suite.powerRepo)
	checker.Assert(powerequip.GetEquippedPower(mysticMage, suite.powerRepo), IsNil)
}

type CreatePowerReportSuite struct {
	teros 			*squaddie.Squaddie
	bandit 			*squaddie.Squaddie
	bandit2 		*squaddie.Squaddie
	blot 			*power.Power
	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository
}

var _ = Suite(&CreatePowerReportSuite{})

func (suite *CreatePowerReportSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
	suite.teros.Offense.Mind = 1

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"

	suite.bandit2 = squaddie.NewSquaddie("bandit")
	suite.bandit2.Identification.Name = "bandit"

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect.DamageBonus = 1

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.blot,
	})

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.bandit2})
}

func (suite *CreatePowerReportSuite) TestPowerReportWhenMissed(checker *C) {
	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerCommit := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)
	checker.Assert(powerCommit.AttackerID, Equals, suite.teros.Identification.ID)
	checker.Assert(powerCommit.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerCommit.AttackingPowerReports, HasLen, 1)
	checker.Assert(powerCommit.AttackingPowerReports[0].WasAHit, Equals, false)
}

func (suite *CreatePowerReportSuite) TestPowerReportWhenHitButNoCrit(checker *C) {
	dieRoller := &testutility.AlwaysHitDieRoller{}

	powerCommit := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)
	checker.Assert(powerCommit.AttackerID, Equals, suite.teros.Identification.ID)
	checker.Assert(powerCommit.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerCommit.AttackingPowerReports, HasLen, 1)
	checker.Assert(powerCommit.AttackingPowerReports[0].AttackerID, Equals, suite.teros.Identification.ID)
	checker.Assert(powerCommit.AttackingPowerReports[0].TargetID, Equals, suite.bandit.Identification.ID)
	checker.Assert(powerCommit.AttackingPowerReports[0].PowerID, Equals, suite.blot.ID)

	checker.Assert(powerCommit.AttackingPowerReports[0].WasAHit, Equals, true)
	checker.Assert(powerCommit.AttackingPowerReports[0].WasACriticalHit, Equals, false)
	checker.Assert(powerCommit.AttackingPowerReports[0].DamageTaken, Equals, 2)
	checker.Assert(powerCommit.AttackingPowerReports[0].BarrierDamage, Equals, 0)
}

func (suite *CreatePowerReportSuite) TestPowerReportWhenCrits(checker *C) {
	dieRoller := &testutility.AlwaysHitDieRoller{}
	suite.blot.AttackEffect.CriticalHitThreshold = 900

	powerCommit := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)
	checker.Assert(powerCommit.AttackerID, Equals, suite.teros.Identification.ID)
	checker.Assert(powerCommit.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerCommit.AttackingPowerReports, HasLen, 1)
	checker.Assert(powerCommit.AttackingPowerReports[0].WasAHit, Equals, true)
	checker.Assert(powerCommit.AttackingPowerReports[0].WasACriticalHit, Equals, true)
	checker.Assert(powerCommit.AttackingPowerReports[0].DamageTaken, Equals, 4)
	checker.Assert(powerCommit.AttackingPowerReports[0].BarrierDamage, Equals, 0)
}

func (suite *CreatePowerReportSuite) TestReportPerTarget(checker *C) {
	dieRoller := &testutility.AlwaysMissDieRoller{}

	powerCommit := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID, suite.bandit2.Identification.ID},
			PowerID:           suite.blot.ID,
			PowerRepo:         suite.powerRepo,
		},
		dieRoller,
	)
	checker.Assert(powerCommit.AttackerID, Equals, suite.teros.Identification.ID)
	checker.Assert(powerCommit.PowerID, Equals, suite.blot.ID)

	checker.Assert(powerCommit.AttackingPowerReports, HasLen, 2)
	checker.Assert(powerCommit.AttackingPowerReports[0].TargetID, Equals, suite.bandit.Identification.ID)
	checker.Assert(powerCommit.AttackingPowerReports[1].TargetID, Equals, suite.bandit2.Identification.ID)
}
