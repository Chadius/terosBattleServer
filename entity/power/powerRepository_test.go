package power_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	. "gopkg.in/check.v1"
)

type PowerCreationSuite struct{
	spear *power.Power
	spear2 *power.Power
	repo *power.Repository
}

var _ = Suite(&PowerCreationSuite{})

func (suite *PowerCreationSuite) TestGenerateRandomIDForNewPowers(checker *C) {
	powerWithoutID := power.NewPower("New Attack")
	checker.Assert(powerWithoutID.ID, NotNil)
	checker.Assert(powerWithoutID.ID, Not(Equals), "")
}

func (suite *PowerCreationSuite) SetUpTest(checker *C) {
	suite.spear = power.NewPower("Spear")
	suite.spear.PowerType = power.Physical
	suite.spear.ID = "spearLevel1"
	suite.spear.AttackEffect = &power.AttackingEffect{
		ToHitBonus: 1,
	}

	suite.spear2 = power.NewPower("Spear")
	suite.spear2.PowerType = power.Physical
	suite.spear2.ID = "spearLevel2"
	suite.spear2.AttackEffect  = &power.AttackingEffect{
		ToHitBonus: 2,
	}

	newPowers := []*power.Power{suite.spear, suite.spear2}

	suite.repo = power.NewPowerRepository()
	suite.repo.AddSlicePowerSource(newPowers)
}

func (suite *PowerCreationSuite) TestAddPowersToNewRepository(checker *C) {
	newRepo := power.NewPowerRepository()
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 0)
	spear := power.NewPower("Spear")
	spear.PowerType = power.Physical
	newPowers := []*power.Power{spear}
	success, _ := newRepo.AddSlicePowerSource(newPowers)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)
}

func (suite *PowerCreationSuite) TestUsesIDToRetrievePowers(checker *C) {
	checker.Assert(suite.repo.GetNumberOfPowers(), Equals, 2)

	spearLevel1FromRepo := suite.repo.GetPowerByID(suite.spear.ID)
	checker.Assert(spearLevel1FromRepo.Name, Equals, "Spear")
	checker.Assert(spearLevel1FromRepo.ID, Equals, suite.spear.ID)
	checker.Assert(spearLevel1FromRepo.AttackEffect.ToHitBonus, Equals, suite.spear.AttackEffect.ToHitBonus)

	spearLevel2FromRepo := suite.repo.GetPowerByID(suite.spear2.ID)
	checker.Assert(spearLevel2FromRepo.Name, Equals, "Spear")
	checker.Assert(spearLevel2FromRepo.ID, Equals, suite.spear2.ID)
	checker.Assert(spearLevel2FromRepo.AttackEffect.ToHitBonus, Equals, suite.spear2.AttackEffect.ToHitBonus)
}

func (suite *PowerCreationSuite) TestReturnNilIfIDDoesNotExist(checker *C) {
	nonExistentPower := suite.repo.GetPowerByID("Nope")
	checker.Assert(nonExistentPower, IsNil)
}

func (suite *PowerCreationSuite) TestSearchForPowerByName(checker *C) {
	allSpearPowers := suite.repo.GetAllPowersByName("Spear")
	checker.Assert(allSpearPowers, HasLen, 2)

	hasSpearPower := false
	hasSpear2Power := false
	for _, power := range allSpearPowers {
		if power.ID == suite.spear.ID {
			hasSpearPower = true
		}
		if power.ID == suite.spear2.ID {
			hasSpear2Power = true
		}
	}

	checker.Assert(hasSpearPower, Equals, true)
	checker.Assert(hasSpear2Power, Equals, true)
}

func (suite *PowerCreationSuite) TestLoadPowersWithJSON(checker *C) {
	jsonByteStream := []byte(`[{
					"name": "Scimitar",
					"id": "deadbeef",
					"damage_bonus": 2,
					"power_type": "Physical",
					"attack_effect": {
						"can_counter_attack": true,
						"counter_attack_penalty": -2
					}
				}]`)
	newRepo := power.NewPowerRepository()
	success, _ := newRepo.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)
}
func (suite *PowerCreationSuite) TestStopsLoadingUponFirstInvalidPower(checker *C) {
	jsonByteStream := []byte(`[{
				"name": "Scimitar",
				"id": "deadbeef",
				"power_type": "Physical"
			},{
				"name": "Scimitar2",
				"id": "deadbeee",
				"power_type": "mystery"
			}]`)
	success, err := suite.repo.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, ErrorMatches, "AttackingPower 'Scimitar2' has unknown power_type: 'mystery'")
}
func (suite *PowerCreationSuite) TestLoadPowersWithYAML(checker *C) {
	yamlByteStream := []byte(`-
  name: Scimitar
  id: deadbeef
  power_type: Physical
  attack_effect:
    damage_bonus: 2
    can_counter_attack: true
    counter_attack_penalty: -2
`)
	newRepo := power.NewPowerRepository()
	success, _ := newRepo.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)

	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.Name, Equals, "Scimitar")
	checker.Assert(scimitar.ID, Equals, "deadbeef")
	checker.Assert(scimitar.AttackEffect.DamageBonus, Equals, 2)
	checker.Assert(scimitar.AttackEffect.CanCounterAttack, Equals, true)
	checker.Assert(scimitar.AttackEffect.CounterAttackToHitPenalty, Equals, -2)
}
