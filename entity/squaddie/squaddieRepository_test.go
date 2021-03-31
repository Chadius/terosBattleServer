package squaddie_test

import (
	"bytes"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieRepositorySuite struct {
	squaddieFactory *squaddie.Repository
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieRepositorySuite{})

func (suite *SquaddieRepositorySuite) SetUpTest(checker *C) {
	suite.squaddieFactory = squaddie.NewSquaddieRepository()
	suite.teros = squaddie.NewSquaddie("Teros")
	suite.teros.Armor = 2
	suite.teros.Dodge = 3
	suite.teros.Deflect = 4
	suite.teros.MaxBarrier = 1
}

func (suite *SquaddieRepositorySuite) TestUseJSONSource(checker *C) {
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 0)
	jsonByteStream := []byte(`[{
				"name": "teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestCloneSquaddies(checker *C) {
	jsonByteStream := []byte(`[{
				"id": "squaddieID",
				"name": "teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("squaddieID")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Name, Equals, "teros")
	checker.Assert(suite.teros.Aim, Equals, 5)

	missingno := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestGetExistingSquaddieUsingID(checker *C) {
	jsonByteStream := []byte(`[{
				"ID": "12345",
				"name": "teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("12345")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Name, Equals, "teros")
	checker.Assert(suite.teros.Aim, Equals, 5)

	missingno := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestClonedSquaddiesHaveDifferentID(checker *C) {
	jsonByteStream := []byte(`[{
				"ID": "terosID",
				"name": "teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	teros0 := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("terosID")
	teros1 := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("terosID")
	checker.Assert(teros0, Not(Equals), teros1)
}

func (suite *SquaddieRepositorySuite) TestLoadClassLevels(checker *C) {
	jsonByteStream := []byte(`[{
				"id": "terosSquaddieID",
				"name": "teros",
				"aim": 5,
				"affiliation": "Player",
				"class_levels": {
                  "1": {
                    "id": "1",
                    "name": "Mage",
                    "levels_gained": ["123"]
                  },
                  "2": {
                    "id": "2",
                    "name": "Dimension Walker"
                  }
               }
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("terosSquaddieID")
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
}

func (suite *SquaddieRepositorySuite) TestStopLoadingSquaddiesWhenInvalid(checker *C) {
	jsonByteStream := []byte(`[{
				"Name": "teros",
				"Affiliation": "Player"
			},{
				"Name": "teros2",
				"Affiliation": "Unknown Affiliation"
			}]`)
	success, err := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, ErrorMatches, "Squaddie has unknown affiliation: 'Unknown Affiliation'")
}

func (suite *SquaddieRepositorySuite) TestCreateSquaddiesWithMovement(checker *C) {
	jsonByteStream := []byte(`[
			{
				"id": "Soldier",
				"name": "Soldier",
				"affiliation": "Player",
				"movement": { "distance": 5, "type": "foot"}
			},
			{
				"id": "Scout",
				"name": "Scout",
				"affiliation": "Player",
				"movement": { "distance": 4, "type": "light"}
			},
			{
				"id": "Bird",
				"name": "Bird",
				"affiliation": "Player",
				"movement": { "distance": 3, "type": "fly", "hit_and_run": true}
			},
			{
				"id": "Teleporter",
				"name": "Teleporter",
				"affiliation": "Player",
				"movement": { "distance": 2, "type": "teleport"}
			}
			]`)

	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Soldier")
	checker.Assert(soldier.Name, Equals, "Soldier")
	checker.Assert(soldier.GetMovementDistancePerRound(), Equals, 5)
	checker.Assert(soldier.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
	checker.Assert(soldier.CanHitAndRun(), Equals, false)

	scout := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Scout")
	checker.Assert(scout.Name, Equals, "Scout")
	checker.Assert(scout.GetMovementDistancePerRound(), Equals, 4)
	checker.Assert(scout.GetMovementType(), Equals, squaddie.MovementType(squaddie.Light))
	checker.Assert(scout.CanHitAndRun(), Equals, false)

	bird := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Bird")
	checker.Assert(bird.Name, Equals, "Bird")
	checker.Assert(bird.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(bird.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(bird.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Teleporter")
	checker.Assert(teleporter.Name, Equals, "Teleporter")
	checker.Assert(teleporter.GetMovementDistancePerRound(), Equals, 2)
	checker.Assert(teleporter.GetMovementType(), Equals, squaddie.MovementType(squaddie.Teleport))
	checker.Assert(teleporter.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestCanGetExisitingSquaddies(checker *C) {
	originalSquaddie := squaddie.NewSquaddie("Original")
	suite.squaddieFactory.AddSquaddies([]*squaddie.Squaddie{originalSquaddie})
	referencedSquaddie := suite.squaddieFactory.GetOriginalSquaddieByID(originalSquaddie.ID)
	checker.Assert(referencedSquaddie, Equals, originalSquaddie)
}

func (suite *SquaddieRepositorySuite) TestMarshallIntoJSON(checker *C) {
	byteStream, err := suite.squaddieFactory.MarshalSquaddieIntoJSON(suite.teros)
	checker.Assert(err, IsNil)

	hasIDNameAffiliationInJSON := bytes.Contains(byteStream, []byte(fmt.Sprintf(`"id":"%s","name":"Teros","affiliation":"Player"`, suite.teros.ID)))
	checker.Assert(hasIDNameAffiliationInJSON, Equals, true)

	hasDefaultHitPointsInJSON := bytes.Contains(byteStream, []byte(`"current_hit_points":5,"max_hit_points":5`))
	checker.Assert(hasDefaultHitPointsInJSON, Equals, true)

	hasOffensiveStatsInJSON := bytes.Contains(byteStream, []byte(`"aim":0,"strength":0,"mind":0`))
	checker.Assert(hasOffensiveStatsInJSON, Equals, true)

	hasDefensiveStatsInJSON := bytes.Contains(byteStream, []byte(`"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2`))
	checker.Assert(hasDefensiveStatsInJSON, Equals, true)

	hasDefaultMovementInJSON := bytes.Contains(byteStream, []byte(`"movement":{"distance":3,"type":"foot","hit_and_run":false}`))
	checker.Assert(hasDefaultMovementInJSON, Equals, true)

	hasNoPowersInJSON := bytes.Contains(byteStream, []byte(`"powers":[]`))
	checker.Assert(hasNoPowersInJSON, Equals, true)

	hasNoClassesInJSON := bytes.Contains(byteStream, []byte(`"base_class":"","current_class":"","class_levels":{}`))
	checker.Assert(hasNoClassesInJSON, Equals, true)
}

func (suite *SquaddieRepositorySuite) TestMarshallSquaddieMovement(checker *C) {
	suite.teros.Movement.Type = squaddie.Teleport
	suite.teros.Movement.Distance = 8
	suite.teros.Movement.HitAndRun = true

	byteStream, err := suite.squaddieFactory.MarshalSquaddieIntoJSON(suite.teros)
	checker.Assert(err, IsNil)
	movementJSON := `"movement":{"distance":8,"type":"teleport","hit_and_run":true}`
	containsPowersJson := bytes.Contains(byteStream, []byte(movementJSON))
	checker.Assert(containsPowersJson, Equals, true)
}

func (suite *SquaddieRepositorySuite) TestMarshallSquaddiePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.AddInnatePower(attackA)
	byteStream, err := suite.squaddieFactory.MarshalSquaddieIntoJSON(suite.teros)
	checker.Assert(err, IsNil)

	powersJSON := fmt.Sprintf(`"powers":[{"name":"Attack Formation A","id":"%s"}]`, attackA.ID)
	containsPowersJson := bytes.Contains(byteStream, []byte(powersJSON))
	checker.Assert(containsPowersJson, Equals, true)
}

func (suite *SquaddieRepositorySuite) TestLoadSquaddieByYAML(checker *C) {
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 0)
	yamlByteStream := []byte(`-
  name: teros
  aim: 5
  max_barrier: 3
  affiliation: Player
`)
	suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestLoadClassLevelsYAML(checker *C) {
	yamlByteStream := []byte(`-
  name: teros
  id: terosSquaddieID
  aim: 5
  max_barrier: 3
  affiliation: Player
  class_levels:
    '1':
      id: '1'
      name: Mage
      levels_gained:
      - '123'
    '2':
      id: '2'
      name: Dimension Walker
`)
	success, _ := suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("terosSquaddieID")
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
}

func (suite *SquaddieRepositorySuite) TestStopLoadingSquaddiesUponFirstInvalid(checker *C) {
	yamlByteStream := []byte(`-
  name: teros
  affiliation: Player
-
  name: teros2
  affiliation: Unknown Affiliation`)
	success, err := suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, NotNil)
	checker.Assert(err.Error(), Equals, "Squaddie has unknown affiliation: 'Unknown Affiliation'")
}

func (suite *SquaddieRepositorySuite) TestLoadSquaddiesWithDifferentMovementYAML(checker *C) {
	yamlByteStream := []byte(`-
  id: Soldier
  name: Soldier
  affiliation: Player
  movement:
    distance: 5
    type: foot
-
  id: Scout
  name: Scout
  affiliation: Player
  movement:
    distance: 4
    type: light
-
  id: Bird
  name: Bird
  affiliation: Player
  movement:
    distance: 3
    type: fly
    hit_and_run: true
-
  id: Teleporter
  name: Teleporter
  affiliation: Player
  movement:
    distance: 2
    type: teleport
    hit_and_run: false
`)

	success, _ := suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Soldier")
	checker.Assert(soldier.Name, Equals, "Soldier")
	checker.Assert(soldier.GetMovementDistancePerRound(), Equals, 5)
	checker.Assert(soldier.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
	checker.Assert(soldier.CanHitAndRun(), Equals, false)

	scout := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Scout")
	checker.Assert(scout.Name, Equals, "Scout")
	checker.Assert(scout.GetMovementDistancePerRound(), Equals, 4)
	checker.Assert(scout.GetMovementType(), Equals, squaddie.MovementType(squaddie.Light))
	checker.Assert(scout.CanHitAndRun(), Equals, false)

	bird := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Bird")
	checker.Assert(bird.Name, Equals, "Bird")
	checker.Assert(bird.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(bird.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(bird.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieFactory.CloneSquaddieBasedOnSquaddieID("Teleporter")
	checker.Assert(teleporter.Name, Equals, "Teleporter")
	checker.Assert(teleporter.GetMovementDistancePerRound(), Equals, 2)
	checker.Assert(teleporter.GetMovementType(), Equals, squaddie.MovementType(squaddie.Teleport))
	checker.Assert(teleporter.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestAddSquaddieDirectly(checker *C) {
	success, err := suite.squaddieFactory.AddSquaddies([]*squaddie.Squaddie{squaddie.NewSquaddie("Generic")})
	checker.Assert(success, Equals, true)
	checker.Assert(err, IsNil)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestCloneSquaddie(checker *C) {
	
}

type SquaddieCloneSuite struct {
	squaddieFactory *squaddie.Repository
	base *squaddie.Squaddie
}

var _ = Suite(&SquaddieCloneSuite{})

func (suite *SquaddieCloneSuite) SetUpTest(checker *C) {
	suite.base = squaddie.NewSquaddie("Base")
	suite.squaddieFactory = squaddie.NewSquaddieRepository()
	suite.squaddieFactory.AddSquaddies([]*squaddie.Squaddie{suite.base})
}

func (suite *SquaddieCloneSuite) TestCloneHasAffiliationAndNameNotID(checker *C) {
	suite.base.Affiliation = squaddie.Enemy
	clone, err := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(err, IsNil)
	checker.Assert(clone.Name, Equals, suite.base.Name)
	checker.Assert(clone.Affiliation, Equals, suite.base.Affiliation)
	checker.Assert(clone.ID, Not(Equals), suite.base.ID)
}

func (suite *SquaddieCloneSuite) TestCloneUsesGivenID(checker *C) {
	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "12345")
	checker.Assert(clone.ID, Equals, "12345")
}

func (suite *SquaddieCloneSuite) TestCloneUsesName(checker *C) {
	clone, err := suite.squaddieFactory.CloneAndRenameSquaddie(suite.base, "ClonedSquaddie", "12345")
	checker.Assert(err, IsNil)
	checker.Assert(suite.base.Name, Equals, "Base")
	checker.Assert(clone.Name, Equals, "ClonedSquaddie")
	checker.Assert(clone.ID, Not(Equals), suite.base.ID)
}

func (suite *SquaddieCloneSuite) TestRenameCloneNeedsAName(checker *C) {
	clone, err := suite.squaddieFactory.CloneAndRenameSquaddie(suite.base, "", "12345")
	checker.Assert(err, ErrorMatches, `cannot clone squaddie "Base" without a name`)
	checker.Assert(clone, IsNil)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesBasicStats(checker *C) {
	suite.base.CurrentHitPoints = 1
	suite.base.MaxHitPoints += 5
	suite.base.CurrentBarrier = 2
	suite.base.MaxBarrier += 5

	suite.base.Aim = 2
	suite.base.Strength = 3
	suite.base.Mind = 4
	suite.base.Dodge = 5
	suite.base.Deflect = 6
	suite.base.Armor = 7

	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.CurrentHitPoints, Equals, suite.base.CurrentHitPoints)
	checker.Assert(clone.MaxHitPoints, Equals, suite.base.MaxHitPoints)
	checker.Assert(clone.Aim, Equals, suite.base.Aim)
	checker.Assert(clone.Strength, Equals, suite.base.Strength)
	checker.Assert(clone.Mind, Equals, suite.base.Mind)
	checker.Assert(clone.Dodge, Equals, suite.base.Dodge)
	checker.Assert(clone.Deflect, Equals, suite.base.Deflect)
	checker.Assert(clone.CurrentBarrier, Equals, suite.base.CurrentBarrier)
	checker.Assert(clone.MaxBarrier, Equals, suite.base.MaxBarrier)
	checker.Assert(clone.Armor, Equals, suite.base.Armor)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesMovement(checker *C) {
	suite.base.Movement = squaddie.Movement{
		Distance:  suite.base.Movement.Distance + 2,
		Type:      squaddie.Fly,
		HitAndRun: true,
	}

	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.Movement.Distance, Equals, suite.base.Movement.Distance)
	checker.Assert(clone.Movement.Type, Equals, suite.base.Movement.Type)
	checker.Assert(clone.Movement.HitAndRun, Equals, suite.base.Movement.HitAndRun)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesPowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.base.AddInnatePower(attackA)
	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")

	attackIDNamePairs := clone.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].ID, Equals, attackA.ID)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesClasses(checker *C) {
	initialClass := &squaddieclass.Class{
		ID: "initial",
		Name: "Initial Class",
		BaseClassRequired: false,
		InitialBigLevelID: "",
	}
	advancedClass := &squaddieclass.Class{
		ID: "advanced",
		Name: "Advanced Class",
		BaseClassRequired: true,
		InitialBigLevelID: "advanceLevel0",
	}

	suite.base.AddClass(initialClass)
	suite.base.AddClass(advancedClass)
	suite.base.SetBaseClassIfNoBaseClass(initialClass.ID)
	suite.base.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel0")
	suite.base.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel1")
	suite.base.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel2")

	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.BaseClassID, Equals, suite.base.BaseClassID)
	checker.Assert(clone.CurrentClass, Equals, suite.base.CurrentClass)
	for classID, levelsConsumed := range suite.base.ClassLevelsConsumed {
		checker.Assert(clone.ClassLevelsConsumed[classID], NotNil)

		cloneLevelsConsumed := clone.ClassLevelsConsumed[classID]
		checker.Assert(cloneLevelsConsumed, Not(Equals), levelsConsumed)
		checker.Assert(cloneLevelsConsumed, DeepEquals, levelsConsumed)
	}
}
