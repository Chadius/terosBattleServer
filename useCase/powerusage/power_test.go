package powerusage_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Power uses with other Entities", func() {

	Context("Calculate expected damage from using attacking powers", func() {
		var (
			teros *squaddie.Squaddie
			spear *power.Power
			blot  *power.Power
		)

		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			teros.Name = "Teros"

			spear = power.NewPower("Spear")
			spear.PowerType = power.Physical

			blot = power.NewPower("Blot")
			blot.PowerType = power.Spell
		})

		It("Calculates the To Hit Bonus", func() {
			teros.Aim = 2
			blot.AttackEffect.ToHitBonus = 1

			totalToHitBonus := powerusage.GetPowerToHitBonusWhenUsedBySquaddie(blot, teros, false)
			Expect(totalToHitBonus).To(Equal(3))
		})

		It("Calculates the To Hit Bonus when counterattacking", func() {
			teros.Aim = 2
			blot.AttackEffect.ToHitBonus = 1
			blot.AttackEffect.CounterattackToHitPenalty = -2

			totalToHitBonus := powerusage.GetPowerToHitBonusWhenUsedBySquaddie(blot, teros, true)
			Expect(totalToHitBonus).To(Equal(1))
		})

		Context("Calculate damage bonus", func() {
			BeforeEach(func() {
				teros.Strength = 2
				teros.Mind = 3

				spear.PowerType = power.Physical
				spear.AttackEffect.DamageBonus = 2

				blot.PowerType = power.Spell
				blot.AttackEffect.DamageBonus = 6
			})

			It("Calculates the Damage bonus of physical attacks", func() {
				totalDamageBonus := powerusage.GetPowerDamageBonusWhenUsedBySquaddie(spear, teros)
				Expect(totalDamageBonus).To(Equal(4))
			})

			It("Calculates the Damage bonus of spell attacks", func() {
				totalDamageBonus := powerusage.GetPowerDamageBonusWhenUsedBySquaddie(blot, teros)
				Expect(totalDamageBonus).To(Equal(9))
			})

			It("Calculates the Critical Damage bonus of physical attacks", func() {
				totalDamageBonus := powerusage.GetPowerCriticalDamageBonusWhenUsedBySquaddie(spear, teros)
				Expect(totalDamageBonus).To(Equal(8))
			})
		})

		Context("Calculate to hit penalties values against powers", func() {
			BeforeEach(func() {
				teros.Dodge = 2
				teros.Deflect = 9001

				spear.PowerType = power.Physical

				blot.PowerType = power.Spell
			})

			It("Calculates the to hit reduction against physical attacks", func() {
				toHitPenalty := powerusage.GetPowerToHitPenaltyAgainstSquaddie(spear, teros)
				Expect(toHitPenalty).To(Equal(2))
			})

			It("Calculates the to hit reduction against spell attacks", func() {
				toHitPenalty := powerusage.GetPowerToHitPenaltyAgainstSquaddie(blot, teros)
				Expect(toHitPenalty).To(Equal(9001))
			})
		})

		Context("Calculate damage if the attacker hits the target with the power", func() {
			var (
				bandit *squaddie.Squaddie
			)

			BeforeEach(func() {
				bandit = squaddie.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				teros.Strength = 1
				spear.AttackEffect.DamageBonus = 3
			})

			It("Does full damage against targets without armor or barrier", func() {
				totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(4))
			})

			It("Armor reduces damage against physical attacks", func() {
				bandit.Armor = 3
				totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(1))
			})

			It("Barrier absorbs damage against physical attacks and is depleted first", func() {
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 1
				totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(3))
				Expect(initialBarrierDamage).To(Equal(1))
			})

			It("Will deal no damage if barrier is strong enough", func() {
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 4
				totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(4))
			})

			It("May deal no damage if armor is strong enough", func() {
				bandit.Armor = 4
				totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(0))
			})
		})

		Context("Calculate damage if the attacker hits the target with the power", func() {
			var (
				bandit *squaddie.Squaddie
			)

			BeforeEach(func() {
				bandit = squaddie.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				teros.Mind = 2
				blot.AttackEffect.DamageBonus = 4
			})

			It("Does full damage against targets without armor or barrier", func() {
				totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(6))
			})

			It("Ignores Armor when using spell attacks", func() {
				bandit.Armor = 9001
				totalHealthDamage, _, _ := powerusage.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(6))
			})

			It("Barrier absorbs damage against spell attacks and is depleted first", func() {
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 1
				totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(5))
				Expect(initialBarrierDamage).To(Equal(1))
			})

			It("Will deal no damage if barrier is strong enough", func() {
				bandit.MaxBarrier = 9001
				bandit.CurrentBarrier = 9001
				totalHealthDamage, initialBarrierDamage, _ := powerusage.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(6))
			})

			It("Can deal extra Barrier damage if the barrier absorbs the attack", func() {
				bandit.MaxBarrier = 8
				bandit.CurrentBarrier = 8
				blot.AttackEffect.ExtraBarrierDamage = 2

				totalHealthDamage, initialBarrierDamage, extraBarrierDamage := powerusage.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(6))
				Expect(extraBarrierDamage).To(Equal(2))
			})

			It("Knows extra Barrier damage is reduced if the barrier is depleted", func() {
				bandit.MaxBarrier = 8
				bandit.CurrentBarrier = 7
				blot.AttackEffect.ExtraBarrierDamage = 2

				totalHealthDamage, initialBarrierDamage, extraBarrierDamage := powerusage.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(6))
				Expect(extraBarrierDamage).To(Equal(1))
			})
		})

		Context("Calculate expected damage summary", func() {
			var (
				bandit *squaddie.Squaddie
				bandit2 *squaddie.Squaddie
			)

			BeforeEach(func() {
				bandit = squaddie.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				bandit2 = squaddie.NewSquaddie("Bandit2")
				bandit2.Name = "Bandit2"
			})

			It("Give summary of the physical attack", func() {
				bandit.Armor = 1
				bandit.Dodge = 1
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 1

				teros.Strength = 1
				spear.AttackEffect.DamageBonus = 3

				teros.Mind = 2
				blot.AttackEffect.DamageBonus = 4

				attackingPowerSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
					Power:           spear,
					Attacker:        teros,
					Target:          bandit,
					IsCounterattack: false,
				})
				Expect(attackingPowerSummary.AttackingSquaddieID).To(Equal(teros.ID))
				Expect(attackingPowerSummary.PowerID).To(Equal(spear.ID))
				Expect(attackingPowerSummary.TargetSquaddieID).To(Equal(bandit.ID))
				Expect(attackingPowerSummary.IsACounterAttack).To(BeFalse())
				Expect(attackingPowerSummary.ChanceToHit).To(Equal(15))
				Expect(attackingPowerSummary.DamageTaken).To(Equal(2))
				Expect(attackingPowerSummary.ExpectedDamage).To(Equal(30))
				Expect(attackingPowerSummary.BarrierDamageTaken).To(Equal(1))
				Expect(attackingPowerSummary.ExpectedBarrierDamage).To(Equal(15))
			})

			It("Give summary of the spell attack with barrier burn", func() {
				bandit.Armor = 1
				bandit.Dodge = 1
				bandit.MaxBarrier = 10
				bandit.CurrentBarrier = 10

				teros.Aim = 3
				teros.Mind = 2
				blot.AttackEffect.DamageBonus = 4
				blot.AttackEffect.ExtraBarrierDamage = 3
				attackingPowerSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
					Power:           blot,
					Attacker:        teros,
					Target:          bandit,
					IsCounterattack: false,
				})
				Expect(attackingPowerSummary.ChanceToHit).To(Equal(33))
				Expect(attackingPowerSummary.DamageTaken).To(Equal(0))
				Expect(attackingPowerSummary.ExpectedDamage).To(Equal(0))
				Expect(attackingPowerSummary.BarrierDamageTaken).To(Equal(9))
				Expect(attackingPowerSummary.ExpectedBarrierDamage).To(Equal(9 * 33))
			})

			It("Produces an attack summary for each target", func() {
				powerSummary := powerusage.GetPowerSummary(spear, teros, []*squaddie.Squaddie{bandit, bandit2})
				Expect(powerSummary.UserSquaddieID).To(Equal(teros.ID))
				Expect(powerSummary.PowerID).To(Equal(spear.ID))
				Expect(powerSummary.AttackEffectSummary).To(HaveLen(2))
				Expect(powerSummary.AttackEffectSummary[0].TargetSquaddieID).To(Equal(bandit.ID))
				Expect(powerSummary.AttackEffectSummary[1].TargetSquaddieID).To(Equal(bandit2.ID))
			})
		})

		Context("Critical Hits", func() {
			var (
				bandit *squaddie.Squaddie
			)

			BeforeEach(func() {
				bandit = squaddie.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				teros.Strength = 1
				spear.AttackEffect.DamageBonus = 3
			})

			It("Adds the chance to crit to the attack summary", func() {
				spear.AttackEffect.CriticalHitThreshold = 4
				attackingPowerSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
					Power:           spear,
					Attacker:        teros,
					Target:          bandit,
					IsCounterattack: false,
				})
				Expect(attackingPowerSummary.ChanceToCritical).To(Equal(6))
			})

			It("Doubles the damage before applying armor and barrier to the attack summary", func() {
				bandit.Armor = 1
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 4
				spear.AttackEffect.CriticalHitThreshold = 4
				attackingPowerSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
					Power:           spear,
					Attacker:        teros,
					Target:          bandit,
					IsCounterattack: false,
				})
				Expect(attackingPowerSummary.CriticalDamageTaken).To(Equal(3))
				Expect(attackingPowerSummary.CriticalBarrierDamageTaken).To(Equal(4))
				Expect(attackingPowerSummary.CriticalExpectedDamage).To(Equal(3 * 21))
				Expect(attackingPowerSummary.CriticalExpectedBarrierDamage).To(Equal(4 * 21))
			})

			It("Does not factor critical effects if the attack cannot crit", func() {
				spear.AttackEffect.CriticalHitThreshold = 0
				attackingPowerSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
					Power:           spear,
					Attacker:        teros,
					Target:          bandit,
					IsCounterattack: false,
				})
				Expect(attackingPowerSummary.ChanceToCritical).To(Equal(0))
				Expect(attackingPowerSummary.CriticalDamageTaken).To(Equal(0))
				Expect(attackingPowerSummary.CriticalBarrierDamageTaken).To(Equal(0))
				Expect(attackingPowerSummary.CriticalExpectedDamage).To(Equal(0))
				Expect(attackingPowerSummary.CriticalExpectedBarrierDamage).To(Equal(0))
			})
		})
	})
	Context("Give squaddie powers", func() {
		var (
			teros *squaddie.Squaddie
			powerRepository *power.Repository
			spear *power.Power
		)
		BeforeEach(func() {
			powerRepository = power.NewPowerRepository()

			spear = power.NewPower("Spear")
			spear.PowerType = power.Physical
			spear.ID = "deadbeef"
			newPowers := []*power.Power{spear}
			powerRepository.AddSlicePowerSource(newPowers)

			teros = squaddie.NewSquaddie("Teros")
			teros.Name = "Teros"
		})
		It("Can give Squaddie innate Powers with a repository", func() {
			temporaryPowerReferences := []*power.Reference{{Name: "Spear", ID: spear.ID}}
			numberOfPowersAdded, err := powerusage.LoadAllOfSquaddieInnatePowers(teros, temporaryPowerReferences, powerRepository)
			Expect(numberOfPowersAdded).To(Equal(1))
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spear.ID))
		})
		It("Stop adding Powers to Squaddie if it doesn't exist", func() {
			scimitar := power.NewPower("Scimitar")
			scimitar.PowerType = power.Physical

			temporaryPowerReferences := []*power.Reference{{Name: "Scimitar", ID: scimitar.ID}}
			numberOfPowersAdded, err := powerusage.LoadAllOfSquaddieInnatePowers(teros, temporaryPowerReferences, powerRepository)
			Expect(numberOfPowersAdded).To(Equal(0))
			Expect(err.Error()).To(Equal("squaddie 'Teros' tried to add Power 'Scimitar' but it does not exist"))

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(0))
		})
	})
	Context("Create Power Reports when using Powers", func() {
		var (
			teros *squaddie.Squaddie
			bandit *squaddie.Squaddie
			bandit2 *squaddie.Squaddie
			blot *power.Power
		)

		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			teros.Name = "Teros"
			teros.Mind = 1

			bandit = squaddie.NewSquaddie("Bandit")
			bandit.Name = "Bandit"

			bandit2 = squaddie.NewSquaddie("Bandit")
			bandit2.Name = "Bandit"

			blot = power.NewPower("Blot")
			blot.PowerType = power.Spell
			blot.AttackEffect.DamageBonus = 1
		})

		It("Creates a Power Report saying it missed", func() {
			dieRoller := &testutility.AlwaysMissDieRoller{}

			powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				blot,
				teros,
				[]*squaddie.Squaddie{
					bandit,
				},
				dieRoller,
			)
			Expect(powerResult.AttackerID).To(Equal(teros.ID))
			Expect(powerResult.PowerID).To(Equal(blot.ID))

			Expect(powerResult.AttackingPowerResults).To(HaveLen(1))
			Expect(powerResult.AttackingPowerResults[0].WasAHit).To(BeFalse())
		})

		It("Creates a Power Report when it hits but does not crit", func() {
			dieRoller := &testutility.AlwaysHitDieRoller{}

			powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				blot,
				teros,
				[]*squaddie.Squaddie{
					bandit,
				},
				dieRoller,
			)
			Expect(powerResult.AttackerID).To(Equal(teros.ID))
			Expect(powerResult.PowerID).To(Equal(blot.ID))

			Expect(powerResult.AttackingPowerResults).To(HaveLen(1))
			Expect(powerResult.AttackingPowerResults[0].WasAHit).To(BeTrue())
			Expect(powerResult.AttackingPowerResults[0].WasACriticalHit).To(BeFalse())
			Expect(powerResult.AttackingPowerResults[0].DamageTaken).To(Equal(2))
			Expect(powerResult.AttackingPowerResults[0].BarrierDamage).To(Equal(0))
		})

		It("Creates a Power Report when it hits and crits", func() {
			dieRoller := &testutility.AlwaysHitDieRoller{}
			blot.AttackEffect.CriticalHitThreshold = 900

			powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				blot,
				teros,
				[]*squaddie.Squaddie{
					bandit,
				},
				dieRoller,
			)
			Expect(powerResult.AttackerID).To(Equal(teros.ID))
			Expect(powerResult.PowerID).To(Equal(blot.ID))

			Expect(powerResult.AttackingPowerResults).To(HaveLen(1))
			Expect(powerResult.AttackingPowerResults[0].WasAHit).To(BeTrue())
			Expect(powerResult.AttackingPowerResults[0].WasACriticalHit).To(BeTrue())
			Expect(powerResult.AttackingPowerResults[0].DamageTaken).To(Equal(4))
			Expect(powerResult.AttackingPowerResults[0].BarrierDamage).To(Equal(0))
		})

		It("Creates a Power Report against multiple targets", func() {
			dieRoller := &testutility.AlwaysMissDieRoller{}

			powerResult := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				blot,
				teros,
				[]*squaddie.Squaddie{
					bandit,
					bandit2,
				},
				dieRoller,
			)
			Expect(powerResult.AttackerID).To(Equal(teros.ID))
			Expect(powerResult.PowerID).To(Equal(blot.ID))

			Expect(powerResult.AttackingPowerResults).To(HaveLen(2))
			Expect(powerResult.AttackingPowerResults[0].TargetID).To(Equal(bandit.ID))
			Expect(powerResult.AttackingPowerResults[1].TargetID).To(Equal(bandit2.ID))
		})
	})
	Context("Squaddie commits to using an attack power", func() {
		var (
			teros *squaddie.Squaddie
			spear *power.Power
			scimitar *power.Power
			powerRepo *power.Repository
			bandit *squaddie.Squaddie
			blot *power.Power
			squaddieRepo *squaddie.Repository
		)
		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			spear = power.NewPower("Spear")
			spear.AttackEffect.CanBeEquipped = true

			scimitar = power.NewPower("scimitar the second")
			scimitar.AttackEffect.CanBeEquipped = true

			powerRepo = power.NewPowerRepository()
			powerRepo.AddSlicePowerSource([]*power.Power{
				spear,
				scimitar,
			})

			bandit = squaddie.NewSquaddie("Bandit")
			bandit.Name = "Bandit"

			blot = power.NewPower("Blot")
			blot.PowerType = power.Spell

			terosPowerReferences := []*power.Reference{
				spear.GetReference(),
				scimitar.GetReference(),
				blot.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)

			squaddieRepo = squaddie.NewSquaddieRepository()
			squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
				teros,
				bandit,
			})
		})

		It("Squaddies will equip a power when they Commit", func() {
			dieRoller := &testutility.AlwaysMissDieRoller{}

			powerReport := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				scimitar,
				teros,
				[]*squaddie.Squaddie{
					bandit,
				},
				dieRoller,
			)

			powerusage.CommitPowerUse(powerReport, squaddieRepo, powerRepo)
			Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(scimitar.ID))
		})
		It("Squaddies will keep their previous equipped power if they cannot equip the power they Commit with", func() {
			powerusage.SquaddieEquipPower(teros, scimitar.ID, powerRepo)

			dieRoller := &testutility.AlwaysMissDieRoller{}

			powerReport := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				blot,
				teros,
				[]*squaddie.Squaddie{
					bandit,
				},
				dieRoller,
			)

			powerusage.CommitPowerUse(powerReport, squaddieRepo, powerRepo)
			Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(scimitar.ID))
		})
		It("Squaddies will not equip powers if none exist, even after Committing", func() {
			mysticMage := squaddie.NewSquaddie("Mystic Mage")
			mysticMagePowerReferences := []*power.Reference{
				blot.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(mysticMage, mysticMagePowerReferences, powerRepo)

			squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
				mysticMage,
			})

			dieRoller := &testutility.AlwaysMissDieRoller{}

			powerReport := powerusage.UsePowerAgainstSquaddiesAndGetResults(
				blot,
				mysticMage,
				[]*squaddie.Squaddie{
					bandit,
				},
				dieRoller,
			)

			powerusage.CommitPowerUse(powerReport, squaddieRepo, powerRepo)
			Expect(powerusage.GetEquippedPower(mysticMage, powerRepo)).To(BeNil())
		})
	})
	Context("Target attempts to counter", func() {
		var (
			teros *squaddie.Squaddie
			spear *power.Power
			blot *power.Power

			bandit *squaddie.Squaddie
			axe *power.Power

			powerRepo *power.Repository
			squaddieRepo *squaddie.Repository
		)
		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			spear = power.NewPower("Spear")
			spear.AttackEffect.CanBeEquipped = true
			spear.AttackEffect.CanCounterattack = true
			spear.AttackEffect.CounterattackToHitPenalty = -2

			axe = power.NewPower("axe the second")
			axe.AttackEffect.CanBeEquipped = true
			axe.AttackEffect.CanCounterattack = true
			axe.AttackEffect.CounterattackToHitPenalty = -2

			powerRepo = power.NewPowerRepository()
			powerRepo.AddSlicePowerSource([]*power.Power{
				spear,
				axe,
			})

			blot = power.NewPower("Blot")
			blot.PowerType = power.Spell

			terosPowerReferences := []*power.Reference{
				spear.GetReference(),
				blot.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)

			bandit = squaddie.NewSquaddie("Bandit")
			bandit.Name = "Bandit"
			banditPowerReferences := []*power.Reference{
				axe.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(bandit, banditPowerReferences, powerRepo)

			squaddieRepo = squaddie.NewSquaddieRepository()
			squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
				teros,
				bandit,
			})
		})
		It("Target will counterattack if it has equipped a power that can counter", func() {
			powerusage.SquaddieEquipPower(teros, spear.ID, powerRepo)
			powerusage.SquaddieEquipPower(bandit, axe.ID, powerRepo)

			expectedTerosCounterattackSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
				Power:           spear,
				Attacker:        teros,
				Target:          bandit,
				IsCounterattack: false,
			})
			terosChanceToHit := expectedTerosCounterattackSummary.ChanceToHit

			banditAttackSummary := powerusage.GetExpectedDamage(&powerusage.AttackContext{
				Power:           axe,
				Attacker:        bandit,
				Target:          teros,
				IsCounterattack: false,
			})
			Expect(banditAttackSummary.CounterAttack).NotTo(BeNil())
			Expect(banditAttackSummary.CounterAttack.IsACounterAttack).To(BeTrue())
			Expect(banditAttackSummary.CounterAttack.AttackingSquaddieID).To(Equal(teros.ID))
			Expect(banditAttackSummary.CounterAttack.PowerID).To(Equal(spear.ID))
			Expect(banditAttackSummary.CounterAttack.TargetSquaddieID).To(Equal(bandit.ID))
			Expect(banditAttackSummary.CounterAttack.ChanceToHit).To(Equal(terosChanceToHit - 2))
		})
	})
})
