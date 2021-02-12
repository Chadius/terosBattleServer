package levelUpBenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddieClass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CRUD LevelUpBenefits", func() {
	var (
		repo *levelUpBenefit.Repository
		jsonByteStream []byte
		yamlByteStream []byte
	)
	BeforeEach(func() {
		repo = levelUpBenefit.NewLevelUpBenefitRepository()
		jsonByteStream = []byte(`[
          {
            "id":"abcdefg0",
            "level_up_benefit_type": "small",
            "class_id": "class0",
            "max_hit_points": 1,
            "aim": 0,
            "strength": 2,
            "mind": 3,
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "armor": 7,
            "powers": [
              {
                "name": "Scimitar",
                "id": "deadbeef"
              }
            ],
            "movement": {
              "distance": 1,
              "type": "teleport",
              "hit_and_run": true
            }
      }
]`)
		yamlByteStream = []byte(
`
- id: abcdefg0
  class_id: class0
  level_up_benefit_type: small
  max_hit_points: 1
  aim: 0
  strength: 2
  mind: 3
  dodge: 4
  deflect: 5
  max_barrier: 6
  armor: 7
  powers:
  - name: Scimitar
    id: deadbeef
  movement:
    distance: 1,
    type: teleport
    hit_and_run": true
`)
	})
	Context("Create LevelUpBenefit objects from different sources", func() {
		It("Can create LevelUpBenefits from JSON", func() {
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(0))
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(1))
		})
		It("Can create LevelUpBenefits from YAML", func() {
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(0))
			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(1))
		})
		It("Can add LevelUpBenefits directly", func() {
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(0))
			success, _ := repo.AddLevels([]*levelUpBenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID: "class0",
					ID: "level0",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID: "class0",
					ID: "level1",
				},
			})
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(2))
		})
	})
	Context("Can search and retrieve LevelUpBenefit objects using descriptors", func() {
		var (
			mageClass *squaddieClass.Class
			lotsOfSmallLevels []*levelUpBenefit.LevelUpBenefit
			lotsOfBigLevels []*levelUpBenefit.LevelUpBenefit
			levelRepo *levelUpBenefit.Repository
		)
		BeforeEach(func() {
			mageClass = &squaddieClass.Class{
				ID:                "class1",
				Name:              "Mage",
				BaseClassRequired: false,
			}

			lotsOfSmallLevels = []*levelUpBenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall0",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall1",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall2",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall3",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall4",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall5",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall6",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall7",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall8",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall9",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall10",
				},
			}

			lotsOfBigLevels = []*levelUpBenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig0",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig1",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig2",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig3",
				},
			}

			levelRepo = levelUpBenefit.NewLevelUpBenefitRepository()
			levelRepo.AddLevels(lotsOfSmallLevels)
			levelRepo.AddLevels(lotsOfBigLevels)
		})
		It("Can get LevelUpBenefits by squaddie and class", func() {
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			benefits, err := repo.GetLevelUpBenefitsByClassID("class0")
			Expect(err).To(BeNil())
			Expect(len(benefits)).To(Equal(1))

			firstBenefit := benefits[0]
			Expect(firstBenefit.LevelUpBenefitType).To(Equal(levelUpBenefit.Small))
			Expect(firstBenefit.ClassID).To(Equal("class0"))
			Expect(firstBenefit.MaxHitPoints).To(Equal(1))
			Expect(firstBenefit.Aim).To(Equal(0))
			Expect(firstBenefit.Strength).To(Equal(2))
			Expect(firstBenefit.Mind).To(Equal(3))
			Expect(firstBenefit.Dodge).To(Equal(4))
			Expect(firstBenefit.Deflect).To(Equal(5))
			Expect(firstBenefit.MaxBarrier).To(Equal(6))
			Expect(firstBenefit.Armor).To(Equal(7))

			Expect(firstBenefit.PowerIDGained).To(HaveLen(1))
			Expect(firstBenefit.PowerIDGained[0].Name).To(Equal("Scimitar"))
			Expect(firstBenefit.PowerIDGained[0].ID).To(Equal("deadbeef"))
		})
		It("Raises an error if you search for wrong LevelUpBenefits", func() {
			repo.AddJSONSource(jsonByteStream)

			benefits, err := repo.GetLevelUpBenefitsByClassID("Class not found")
			Expect(err.Error()).To(Equal(`no LevelUpBenefits for this class ID: "Class not found"`))
			Expect(benefits).To(HaveLen(0))
		})
		It("can give you big and small levels for a given class", func() {
			levelsByBenefitType, err := levelRepo.GetLevelUpBenefitsForClassByType(mageClass.ID)
			Expect(err).To(BeNil())
			Expect(levelsByBenefitType[levelUpBenefit.Small]).To(HaveLen(11))
			Expect(levelsByBenefitType[levelUpBenefit.Big]).To(HaveLen(4))
		})
		It("raises an error if the class does not exist", func() {
			levelsByBenefitType, err := levelRepo.GetLevelUpBenefitsForClassByType("bad ID")
			Expect(err.Error()).To(Equal(`no LevelUpBenefits for this class ID: "bad ID"`))
			Expect(levelsByBenefitType[levelUpBenefit.Small]).To(HaveLen(0))
			Expect(levelsByBenefitType[levelUpBenefit.Big]).To(HaveLen(0))
		})
	})
	It("Stops loading LevelUpBenefits upon validating the first invalid LevelUpBenefit", func() {
		byteStream := []byte(`[
          {
            "id":"abcdefg0",
            "class_id": "class0",
            "level_up_benefit_type": "small",
            "max_hit_points": 1,
            "aim": 0,
            "strength": 2,
            "mind": 3,
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "armor": 7,
            "powers": [
              {
                "name": "Scimitar",
                "id": "deadbeef"
              }
            ]
          },
		  {
				"level_up_benefit_type": "unknown",
                "class_id": "class0",
				"max_hit_points": 1,
				"aim": 0,
				"strength": 2,
				"mind": 3,
				"dodge": 4,
				"deflect": 5,
				"max_barrier": 6,
				"armor": 7,
				"powers": [{"name": "Scimitar", "id": "deadbeef"}]
          }
]`)
		success, err := repo.AddJSONSource(byteStream)
		Expect(success).To(BeFalse())
		Expect(err.Error()).To(Equal(`unknown level up benefit type: "unknown"`))
	})
})