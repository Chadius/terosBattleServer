package power_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	. "gopkg.in/check.v1"
)

type HealingEffectLoadedFromData struct{
	healStaffJSON []byte
	healStaffYAML []byte
	repo *power.Repository
}

var _ = Suite(&HealingEffectLoadedFromData{})

func (suite *HealingEffectLoadedFromData) SetUpTest(checker *C) {
	suite.healStaffJSON = []byte(`[{
					"name": "Heal",
					"id": "power_heal",
					"power_type": "Spell",
					"healing_effect": {
						"hit_points_healed": 2
					}
				}]`)

	suite.healStaffYAML = []byte(`-
  name: Heal
  id: power_heal
  power_type: Spell
  healing_effect:
    hit_points_healed: 2
`)

	suite.repo = power.NewPowerRepository()
}

func (suite *HealingEffectLoadedFromData) TestLoadFromJSON(checker *C) {
	success, err := suite.repo.AddJSONSource(suite.healStaffJSON)
	checker.Assert(err, IsNil)
	checker.Assert(success, Equals, true)

	healStaff := suite.repo.GetPowerByID("power_heal")
	checker.Assert(healStaff.HealingEffect.HitPointsHealed, Equals, 2)
}

func (suite *HealingEffectLoadedFromData) TestLoadFromYAML(checker *C) {
	success, err := suite.repo.AddYAMLSource(suite.healStaffYAML)
	checker.Assert(err, IsNil)
	checker.Assert(success, Equals, true)

	healStaff := suite.repo.GetPowerByID("power_heal")
	checker.Assert(healStaff.HealingEffect.HitPointsHealed, Equals, 2)
}