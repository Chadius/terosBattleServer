package powerusage_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Squaddies can equip powers", func() {
	var (
		teros *squaddie.Squaddie
		spear *power.Power
		scimitar *power.Power
		blot *power.Power
		powerRepo *power.Repository
	)
	BeforeEach(func() {
		teros = squaddie.NewSquaddie("Teros")
		spear = power.NewPower("Spear")
		spear.AttackEffect.CanBeEquipped = true
		spear.AttackEffect.CanCounterattack = true

		scimitar = power.NewPower("scimitar the second")
		scimitar.AttackEffect.CanBeEquipped = true

		blot = power.NewPower("Magic spell Blot")

		powerRepo = power.NewPowerRepository()
		powerRepo.AddSlicePowerSource([]*power.Power{
			spear,
			scimitar,
			blot,
		})
	})
	It("Squaddie equips the first power by default", func() {
		terosPowerReferences := []*power.Reference{
			spear.GetReference(),
			scimitar.GetReference(),
			blot.GetReference(),
		}
		powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
		Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(spear.ID))
	})
	It("Squaddie equips the first power it can equip by default", func() {
		terosPowerReferences := []*power.Reference{
			blot.GetReference(),
			spear.GetReference(),
			scimitar.GetReference(),
		}
		powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
		Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(spear.ID))
	})
	It("If no powers can be equipped, Squaddie cannot equip any of them", func() {
		terosPowerReferences := []*power.Reference{
			blot.GetReference(),
		}
		powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
		Expect(powerusage.GetEquippedPower(teros, powerRepo)).To(BeNil())
	})
	It("Squaddie can equip different powers upon command", func() {
		terosPowerReferences := []*power.Reference{
			spear.GetReference(),
			scimitar.GetReference(),
			blot.GetReference(),
		}
		powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
		success := powerusage.SquaddieEquipPower(teros, scimitar.ID, powerRepo)
		Expect(success).To(BeTrue())
		Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(scimitar.ID))
	})
	It("Squaddie cannot equip powers that cannot be equipped", func() {
		terosPowerReferences := []*power.Reference{
			spear.GetReference(),
			scimitar.GetReference(),
			blot.GetReference(),
		}
		powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
		success := powerusage.SquaddieEquipPower(teros, blot.ID, powerRepo)
		Expect(success).To(BeFalse())
		Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(spear.ID))
	})
	It("Squaddie cannot equip powers that do not exist", func() {
		success := powerusage.SquaddieEquipPower(teros, "kwyjibo", powerRepo)
		Expect(success).To(BeFalse())
		Expect(powerusage.GetEquippedPower(teros, powerRepo)).To(BeNil())
	})
	It("Squaddie cannot equip powers they do not own", func() {
		notTerosPower := power.NewPower("Does not belong to Teros")
		notTerosPower.AttackEffect.CanBeEquipped = true
		powerRepo.AddSlicePowerSource([]*power.Power{
			notTerosPower,
		})

		terosPowerReferences := []*power.Reference{
			spear.GetReference(),
			scimitar.GetReference(),
			blot.GetReference(),
		}
		powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
		success := powerusage.SquaddieEquipPower(teros, notTerosPower.ID, powerRepo)
		Expect(success).To(BeFalse())
		Expect(powerusage.GetEquippedPower(teros, powerRepo).ID).To(Equal(spear.ID))
	})

	Context("Powers can be used to counter attack", func() {
		It("Squaddie can counter because it equipped a weapon that can counter", func() {
			terosPowerReferences := []*power.Reference{
				spear.GetReference(),
				scimitar.GetReference(),
				blot.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
			Expect(powerusage.CanSquaddieCounterWithEquippedWeapon(teros, powerRepo)).To(BeTrue())
		})
		It("Squaddie cannot counter because it equipped a weapon that cannot counter", func() {
			terosPowerReferences := []*power.Reference{
				spear.GetReference(),
				scimitar.GetReference(),
				blot.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
			powerusage.SquaddieEquipPower(teros, scimitar.ID, powerRepo)
			Expect(powerusage.CanSquaddieCounterWithEquippedWeapon(teros, powerRepo)).To(BeFalse())
		})
		It("Squaddie cannot counter because it does not have an equippable power", func() {
			terosPowerReferences := []*power.Reference{
				blot.GetReference(),
			}
			powerusage.LoadAllOfSquaddieInnatePowers(teros, terosPowerReferences, powerRepo)
			Expect(powerusage.CanSquaddieCounterWithEquippedWeapon(teros, powerRepo)).To(BeFalse())
		})
	})
})
