package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

var (
	levelRepo *levelupbenefit.Repository
	jsonByteStream []byte
	yamlByteStream []byte
	mageClass *squaddieclass.Class
	lotsOfSmallLevels []*levelupbenefit.LevelUpBenefit
	lotsOfBigLevels []*levelupbenefit.LevelUpBenefit
)

func Test(t *testing.T) { TestingT(t) }

type LevelUpBenefitRepositorySuite struct{}

var _ = Suite(&LevelUpBenefitRepositorySuite{})

func (s *LevelUpBenefitRepositorySuite) SetUpTest(c *C) {
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

	mageClass = &squaddieclass.Class{
		ID:                "class1",
		Name:              "Mage",
		BaseClassRequired: false,
	}

	lotsOfSmallLevels = (&testutility.LevelGenerator{
		Instructions: &testutility.LevelGeneratorInstruction{
			NumberOfLevels: 11,
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
}

func (s *LevelUpBenefitRepositorySuite) TearDownTest(c *C) {
	levelRepo = nil
	jsonByteStream = []byte(``)
	yamlByteStream = []byte(``)
	mageClass = nil
	lotsOfSmallLevels = []*levelupbenefit.LevelUpBenefit{}
	lotsOfBigLevels = []*levelupbenefit.LevelUpBenefit{}
}

func (s *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromJSON(checker *C) {
	levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
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
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := levelRepo.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 1)
}

func (s *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromYAML(checker *C) {
	levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
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
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := levelRepo.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 1)
}

func (s *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromASlice(checker *C) {
	levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            "class0",
			ID:                 "level0",
		},
		{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            "class0",
			ID:                 "level1",
		},
	})
	checker.Assert(success, Equals, true)
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 2)
}

func (s *LevelUpBenefitRepositorySuite) TestStopLoadingOnFirstInvalidLevelUpBenefit(checker *C) {
	levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
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
	success, err := levelRepo.AddJSONSource(byteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err.Error(), Equals, `unknown level up benefit type: "unknown"`)
}

func (s *LevelUpBenefitRepositorySuite) TestCanSearchLevelUpBenefits(checker *C) {
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
	levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	success, _ := levelRepo.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	benefits, err := levelRepo.GetLevelUpBenefitsByClassID("class0")
	checker.Assert(err, IsNil)
	checker.Assert(benefits, HasLen, 1)

	firstBenefit := benefits[0]
	checker.Assert(firstBenefit.LevelUpBenefitType, Equals, levelupbenefit.Small)
	checker.Assert(firstBenefit.ClassID, Equals, "class0")
	checker.Assert(firstBenefit.MaxHitPoints, Equals, 1)
	checker.Assert(firstBenefit.Aim, Equals, 0)
	checker.Assert(firstBenefit.Strength, Equals, 2)
	checker.Assert(firstBenefit.Mind, Equals, 3)
	checker.Assert(firstBenefit.Dodge, Equals, 4)
	checker.Assert(firstBenefit.Deflect, Equals, 5)
	checker.Assert(firstBenefit.MaxBarrier, Equals, 6)
	checker.Assert(firstBenefit.Armor, Equals, 7)

	checker.Assert(firstBenefit.PowerIDGained, HasLen, 1)
	checker.Assert(firstBenefit.PowerIDGained[0].Name, Equals, "Scimitar")
	checker.Assert(firstBenefit.PowerIDGained[0].ID, Equals, "deadbeef")
}

func (s *LevelUpBenefitRepositorySuite) TestRaisesAnErrorWithNonexistentClassID(checker *C) {
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
	levelRepo.AddJSONSource(jsonByteStream)

	benefits, err := levelRepo.GetLevelUpBenefitsByClassID("Class not found")
	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class ID: "Class not found"`)
	checker.Assert(benefits, HasLen, 0)
}

func (s *LevelUpBenefitRepositorySuite) TestGetBigAndSmallLevelsForAGivenClass(checker *C) {
	levelsByBenefitType, err := levelRepo.GetLevelUpBenefitsForClassByType(mageClass.ID)
	checker.Assert(err, IsNil)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 11)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 4)
}

func (s *LevelUpBenefitRepositorySuite) TestRaiseErrorIfClassDoesNotExist(checker *C) {
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
	levelRepo.AddJSONSource(jsonByteStream)
	levelsByBenefitType, err := levelRepo.GetLevelUpBenefitsForClassByType("bad ID")
	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class ID: "bad ID"`)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 0)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 0)
}
