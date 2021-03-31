package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type LevelUpBenefitRepositorySuite struct{
	levelRepo *levelupbenefit.Repository
	jsonByteStream []byte
	yamlByteStream []byte
	mageClass *squaddieclass.Class
	lotsOfSmallLevels []*levelupbenefit.LevelUpBenefit
	lotsOfBigLevels []*levelupbenefit.LevelUpBenefit
}

var _ = Suite(&LevelUpBenefitRepositorySuite{})

func (suite *LevelUpBenefitRepositorySuite) SetUpTest(c *C) {
	suite.jsonByteStream = []byte(`[
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

	suite.mageClass = &squaddieclass.Class{
		ID:                "class1",
		Name:              "Mage",
		BaseClassRequired: false,
	}

	suite.lotsOfSmallLevels = (&testutility.LevelGenerator{
		Instructions: &testutility.LevelGeneratorInstruction{
			NumberOfLevels: 11,
			ClassID:        suite.mageClass.ID,
			PrefixLevelID:  "lotsLevelsSmall",
			Type:           levelupbenefit.Small,
		},
	}).Build()

	suite.lotsOfBigLevels = (&testutility.LevelGenerator{
		Instructions: &testutility.LevelGeneratorInstruction{
			NumberOfLevels: 4,
			ClassID:        suite.mageClass.ID,
			PrefixLevelID:  "lotsLevelsBig",
			Type:           levelupbenefit.Big,
		},
	}).Build()

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddLevels(suite.lotsOfSmallLevels)
	suite.levelRepo.AddLevels(suite.lotsOfBigLevels)
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromJSON(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.jsonByteStream = []byte(`[
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
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := suite.levelRepo.AddJSONSource(suite.jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 1)
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromYAML(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.yamlByteStream = []byte(
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
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := suite.levelRepo.AddYAMLSource(suite.yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 1)
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromASlice(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
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
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 2)
}

func (suite *LevelUpBenefitRepositorySuite) TestStopLoadingOnFirstInvalidLevelUpBenefit(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
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
	success, err := suite.levelRepo.AddJSONSource(byteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err.Error(), Equals, `unknown level up benefit type: "unknown"`)
}

func (suite *LevelUpBenefitRepositorySuite) TestCanSearchLevelUpBenefits(checker *C) {
	suite.jsonByteStream = []byte(`[
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
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	success, _ := suite.levelRepo.AddJSONSource(suite.jsonByteStream)
	checker.Assert(success, Equals, true)

	benefits, err := suite.levelRepo.GetLevelUpBenefitsByClassID("class0")
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

func (suite *LevelUpBenefitRepositorySuite) TestRaisesAnErrorWithNonexistentClassID(checker *C) {
	suite.jsonByteStream = []byte(`[
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
	suite.levelRepo.AddJSONSource(suite.jsonByteStream)

	benefits, err := suite.levelRepo.GetLevelUpBenefitsByClassID("Class not found")
	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class ID: "Class not found"`)
	checker.Assert(benefits, HasLen, 0)
}

func (suite *LevelUpBenefitRepositorySuite) TestGetBigAndSmallLevelsForAGivenClass(checker *C) {
	levelsByBenefitType, err := suite.levelRepo.GetLevelUpBenefitsForClassByType(suite.mageClass.ID)
	checker.Assert(err, IsNil)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 11)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 4)
}

func (suite *LevelUpBenefitRepositorySuite) TestRaiseErrorIfClassDoesNotExist(checker *C) {
	suite.jsonByteStream = []byte(`[
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
	suite.levelRepo.AddJSONSource(suite.jsonByteStream)
	levelsByBenefitType, err := suite.levelRepo.GetLevelUpBenefitsForClassByType("bad ID")
	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class ID: "bad ID"`)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 0)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 0)
}
