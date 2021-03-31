package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	"github.com/kindrid/gotest"
	"github.com/kindrid/gotest/should"
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

func setUp() {
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

func tearDown() {
	levelRepo = nil
	jsonByteStream = []byte(``)
	yamlByteStream = []byte(``)
	mageClass = nil
	lotsOfSmallLevels = []*levelupbenefit.LevelUpBenefit{}
	lotsOfBigLevels = []*levelupbenefit.LevelUpBenefit{}
}

func TestCreateLevelUpBenefitsFromJSON(t *testing.T) {
	setUp()
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
	gotest.Assert(t, levelRepo.GetNumberOfLevelUpBenefits(), should.Equal, 0)
	success, _ := levelRepo.AddJSONSource(jsonByteStream)
	gotest.Assert(t, success, should.BeTrue)
	gotest.Assert(t, levelRepo.GetNumberOfLevelUpBenefits(), should.Equal, 1)
	tearDown()
}

func TestCreateLevelUpBenefitsFromYAML(t *testing.T) {
	setUp()
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
	gotest.Assert(t, levelRepo.GetNumberOfLevelUpBenefits(), should.Equal, 0)
	success, _ := levelRepo.AddYAMLSource(yamlByteStream)
	gotest.Assert(t, success, should.BeTrue)
	gotest.Assert(t, levelRepo.GetNumberOfLevelUpBenefits(), should.Equal, 1)
	tearDown()
}

func TestCreateLevelUpBenefitsFromASlice(t *testing.T) {
	setUp()
	levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	gotest.Assert(t, levelRepo.GetNumberOfLevelUpBenefits(), should.Equal, 0)
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
	gotest.Assert(t, success, should.BeTrue)
	gotest.Assert(t, levelRepo.GetNumberOfLevelUpBenefits(), should.Equal, 2)
	tearDown()
}

func TestStopLoadingOnFirstInvalidLevelUpBenefit(t *testing.T) {
	setUp()
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
	gotest.Assert(t, success, should.BeFalse)
	gotest.Assert(t, err.Error(), should.Equal, `unknown level up benefit type: "unknown"`)
	tearDown()
}

func TestCanSearchLevelUpBenefits(t *testing.T) {
	setUp()
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
	gotest.Assert(t, success, should.BeTrue)

	benefits, err := levelRepo.GetLevelUpBenefitsByClassID("class0")
	gotest.Assert(t, err, should.BeNil)
	gotest.Assert(t, benefits, should.HaveLength, 1)

	firstBenefit := benefits[0]
	gotest.Assert(t, firstBenefit.LevelUpBenefitType, should.Equal, levelupbenefit.Small)
	gotest.Assert(t, firstBenefit.ClassID, should.Equal, "class0")
	gotest.Assert(t, firstBenefit.MaxHitPoints, should.Equal, 1)
	gotest.Assert(t, firstBenefit.Aim, should.Equal, 0)
	gotest.Assert(t, firstBenefit.Strength, should.Equal, 2)
	gotest.Assert(t, firstBenefit.Mind, should.Equal, 3)
	gotest.Assert(t, firstBenefit.Dodge, should.Equal, 4)
	gotest.Assert(t, firstBenefit.Deflect, should.Equal, 5)
	gotest.Assert(t, firstBenefit.MaxBarrier, should.Equal, 6)
	gotest.Assert(t, firstBenefit.Armor, should.Equal, 7)

	gotest.Assert(t, firstBenefit.PowerIDGained, should.HaveLength, 1)
	gotest.Assert(t, firstBenefit.PowerIDGained[0].Name, should.Equal, "Scimitar")
	gotest.Assert(t, firstBenefit.PowerIDGained[0].ID, should.Equal, "deadbeef")
	tearDown()
}

func TestRaisesAnErrorWithNonexistentClassID(t *testing.T) {
	setUp()
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
	gotest.Assert(t, err.Error(), should.Equal, `no LevelUpBenefits for this class ID: "Class not found"`)
	gotest.Assert(t, benefits, should.HaveLength, 0)
	tearDown()
}

func TestGetBigAndSmallLevelsForAGivenClass(t *testing.T) {
	setUp()
	levelsByBenefitType, err := levelRepo.GetLevelUpBenefitsForClassByType(mageClass.ID)
	gotest.Assert(t, err, should.BeNil)
	gotest.Assert(t, levelsByBenefitType[levelupbenefit.Small], should.HaveLength, 11)
	gotest.Assert(t, levelsByBenefitType[levelupbenefit.Big], should.HaveLength, 4)
	tearDown()
}

func TestRaiseErrorIfClassDoesNotExist(t *testing.T) {
	setUp()
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
	gotest.Assert(t, err.Error(), should.Equal, `no LevelUpBenefits for this class ID: "bad ID"`)
	gotest.Assert(t, levelsByBenefitType[levelupbenefit.Small], should.HaveLength, 0)
	gotest.Assert(t, levelsByBenefitType[levelupbenefit.Big], should.HaveLength, 0)
	tearDown()
}
