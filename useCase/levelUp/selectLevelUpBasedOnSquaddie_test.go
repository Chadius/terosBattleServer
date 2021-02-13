package levelup_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/usecase/levelup"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Squaddie choosing and gaining levels", func() {
	Context("Squaddie knows about big and small level ups", func() {
		var (
			teros *squaddie.Squaddie
			mageClass *squaddieclass.Class
			onlySmallLevelsClass *squaddieclass.Class
			classWithInitialLevel *squaddieclass.Class
			lotsOfSmallLevels []*levelupbenefit.LevelUpBenefit
			lotsOfBigLevels []*levelupbenefit.LevelUpBenefit
			classRepo *squaddieclass.Repository
			levelRepo *levelupbenefit.Repository
		)
		BeforeEach(func() {
			mageClass = &squaddieclass.Class{
				ID:                "class0",
				Name:              "Mage",
				BaseClassRequired: false,
			}

			onlySmallLevelsClass = &squaddieclass.Class{
				ID:                "class1",
				Name:              "SmallLevels",
				BaseClassRequired: false,
			}

			classWithInitialLevel = &squaddieclass.Class{
				ID:                "classWithInitialLevel",
				Name:              "Class wants big level first",
				BaseClassRequired: false,
				InitialBigLevelID: "classWithInitialLevelThisIsTakenFirst",
			}

			classRepo = squaddieclass.NewRepository()
			classRepo.AddListOfClasses([]*squaddieclass.Class{mageClass, onlySmallLevelsClass, classWithInitialLevel})

			lotsOfSmallLevels = (&testutility.LevelGenerator{
				Instructions: &testutility.LevelGeneratorInstruction{
					NumberOfLevels: 10,
					ClassID:        mageClass.ID,
					PrefixLevelID:  "lotsLevelsSmall",
					Type:           levelupbenefit.Small,
				},
			}).Build()

			lotsOfBigLevels = (&testutility.LevelGenerator{
				Instructions: &testutility.LevelGeneratorInstruction{
					NumberOfLevels: 4,
					ClassID:        mageClass.ID,
					PrefixLevelID:  "lotsLevelsBig",
					Type:           levelupbenefit.Big,
				},
			}).Build()

			levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
			levelRepo.AddLevels(lotsOfSmallLevels)
			levelRepo.AddLevels(lotsOfBigLevels)
			levelRepo.AddLevels((&testutility.LevelGenerator{
				Instructions: &testutility.LevelGeneratorInstruction{
					NumberOfLevels: 2,
					ClassID:        onlySmallLevelsClass.ID,
					PrefixLevelID:  "smallLevel",
					Type:           levelupbenefit.Small,
				},
			}).Build())

			levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelupbenefit.Small,
					ClassID:            classWithInitialLevel.ID,
					ID:                 "classWithInitialLevel0",
				},
				{
					LevelUpBenefitType: levelupbenefit.Small,
					ClassID:            classWithInitialLevel.ID,
					ID:                 "classWithInitialLevel1",
				},
				{
					LevelUpBenefitType: levelupbenefit.Small,
					ClassID:            classWithInitialLevel.ID,
					ID:                 "classWithInitialLevel2",
				},
				{
					LevelUpBenefitType: levelupbenefit.Big,
					ClassID:            classWithInitialLevel.ID,
					ID:                 "classWithInitialLevelThisIsTakenFirst",
				},
				{
					LevelUpBenefitType: levelupbenefit.Big,
					ClassID:            classWithInitialLevel.ID,
					ID:                 "classWithInitialLevelThisShouldNotBeTakenFirst",
				},
			})

			teros = squaddie.NewSquaddie("Teros")
			teros.AddClass(mageClass)
		})
		It("Only considers small levels when calculating its class level", func() {
			teros.AddClass(mageClass)
			for index, _ := range [5]int{} {
				levelup.ImproveSquaddie(lotsOfSmallLevels[index], teros, nil)
			}

			levelup.ImproveSquaddie(lotsOfBigLevels[0], teros, nil)
			levelup.ImproveSquaddie(lotsOfBigLevels[1], teros, nil)

			classLevels := levelup.GetSquaddieClassLevels(teros, levelRepo)
			Expect(classLevels[mageClass.ID]).To(Equal(5))
		})
		Context("Choosing levels when leveling up a squaddie", func() {
			It("gets one small and one big level if the squaddie has an odd class level", func() {
				teros.AddClass(mageClass)
				teros.SetClass(mageClass.ID)
				err := levelup.ImproveSquaddieBasedOnLevel(teros, lotsOfBigLevels[0].ID, levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelup.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[mageClass.ID]).To(Equal(1))
				Expect(teros.ClassLevelsConsumed[mageClass.ID].LevelsConsumed).To(HaveLen(2))

				hasSmallLevel := teros.ClassLevelsConsumed[mageClass.ID].AnyLevelsConsumed(func(consumedLevelID string) bool {
					return levelupbenefit.AnyLevelUpBenefits(lotsOfSmallLevels, func(level *levelupbenefit.LevelUpBenefit) bool {
						return level.ID == consumedLevelID
					})
				})
				Expect(hasSmallLevel).To(BeTrue())

				hasBigLevel := teros.ClassLevelsConsumed[mageClass.ID].AnyLevelsConsumed(func(consumedLevelID string) bool {
					return levelupbenefit.AnyLevelUpBenefits(lotsOfBigLevels, func(level *levelupbenefit.LevelUpBenefit) bool {
						return level.ID == consumedLevelID
					})
				})
				Expect(hasBigLevel).To(BeTrue())
			})
			It("raises an error if the class cannot be found", func() {
				err := levelup.ImproveSquaddieBasedOnLevel(teros, lotsOfBigLevels[0].ID, levelRepo, classRepo, nil)
				Expect(err.Error()).To(Equal(`class repository: No class found with ID: ""`))
			})
			It("does not add big levels if there are none available", func() {
				teros.AddClass(onlySmallLevelsClass)
				teros.SetClass(onlySmallLevelsClass.ID)
				err := levelup.ImproveSquaddieBasedOnLevel(teros, lotsOfBigLevels[0].ID, levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelup.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[onlySmallLevelsClass.ID]).To(Equal(1))
				Expect(teros.ClassLevelsConsumed[onlySmallLevelsClass.ID].LevelsConsumed).To(HaveLen(1))
			})
			It("chooses any small level at most once", func() {
				teros.AddClass(onlySmallLevelsClass)
				teros.SetClass(onlySmallLevelsClass.ID)
				levelup.ImproveSquaddieBasedOnLevel(teros, "", levelRepo, classRepo, nil)

				err := levelup.ImproveSquaddieBasedOnLevel(teros, "", levelRepo, classRepo, nil)
				Expect(err).To(BeNil())
				classLevels := levelup.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[onlySmallLevelsClass.ID]).To(Equal(2))
				Expect(teros.ClassLevelsConsumed[onlySmallLevelsClass.ID].LevelsConsumed).To(HaveLen(2))
			})
			It("does not add small levels if there are none available", func() {
				teros.AddClass(onlySmallLevelsClass)
				teros.SetClass(onlySmallLevelsClass.ID)
				levelup.ImproveSquaddieBasedOnLevel(teros, "", levelRepo, classRepo, nil)
				levelup.ImproveSquaddieBasedOnLevel(teros, "", levelRepo, classRepo, nil)
				err := levelup.ImproveSquaddieBasedOnLevel(teros, "", levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelup.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[onlySmallLevelsClass.ID]).To(Equal(2))
				Expect(teros.ClassLevelsConsumed[onlySmallLevelsClass.ID].LevelsConsumed).To(HaveLen(2))
			})
			It("squaddie must choose initial level first if it exists", func() {
				teros.AddClass(classWithInitialLevel)
				teros.SetClass(classWithInitialLevel.ID)
				err := levelup.ImproveSquaddieBasedOnLevel(teros, "classWithInitialLevelThisShouldNotBeTakenFirst", levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelup.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[classWithInitialLevel.ID]).To(Equal(1))
				Expect(teros.ClassLevelsConsumed[classWithInitialLevel.ID].LevelsConsumed).To(HaveLen(2))
				Expect(teros.ClassLevelsConsumed[classWithInitialLevel.ID].IsLevelAlreadyConsumed("classWithInitialLevelThisIsTakenFirst")).To(BeTrue())
				Expect(teros.ClassLevelsConsumed[classWithInitialLevel.ID].IsLevelAlreadyConsumed("classWithInitialLevelThisShouldNotBeTakenFirst")).To(BeFalse())

				levelup.ImproveSquaddieBasedOnLevel(teros, "classWithInitialLevelThisShouldNotBeTakenFirst", levelRepo, classRepo, nil)
				Expect(teros.ClassLevelsConsumed[classWithInitialLevel.ID].LevelsConsumed).To(HaveLen(3))

				levelup.ImproveSquaddieBasedOnLevel(teros, "classWithInitialLevelThisShouldNotBeTakenFirst", levelRepo, classRepo, nil)
				Expect(teros.ClassLevelsConsumed[classWithInitialLevel.ID].LevelsConsumed).To(HaveLen(5))
				Expect(teros.ClassLevelsConsumed[classWithInitialLevel.ID].IsLevelAlreadyConsumed("classWithInitialLevelThisShouldNotBeTakenFirst")).To(BeTrue())
			})
		})
	})
})