package power_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	. "gopkg.in/check.v1"
)

type AttackingEffectCounterAttackPenaltyTest struct {}

var _ = Suite(&AttackingEffectCounterAttackPenaltyTest{})

func (suite *AttackingEffectCounterAttackPenaltyTest) SetUpTest (checker *C) {}

func (suite *AttackingEffectCounterAttackPenaltyTest) TestDefaultPenalty(checker *C) {
	counterAttackingPower := &power.Power{
		Reference:    power.Reference{
			Name: "Static",
			ID:   "power0",
		},
		PowerType:    power.Physical,
		AttackEffect: &power.AttackingEffect{
			ToHitBonus:                    0,
			DamageBonus:                   0,
			CanCounterAttack:              true,
			CounterAttackPenaltyReduction: 0,
			CriticalEffect:                nil,
		},
	}

	counterAttackPenalty, err := counterAttackingPower.AttackEffect.CounterAttackPenalty()
	checker.Assert(err, IsNil)
	checker.Assert(counterAttackPenalty, Equals, -2)
}

func (suite *AttackingEffectCounterAttackPenaltyTest) TestRaisesErrorIfPowerCannotCounterAttack(checker *C) {
	cannotCounterWithThisPower := &power.Power{
		Reference:    power.Reference{
			Name: "Static",
			ID:   "power0",
		},
		PowerType:    power.Physical,
		AttackEffect: &power.AttackingEffect{
			ToHitBonus:                    0,
			DamageBonus:                   0,
			CanCounterAttack:              false,
			CounterAttackPenaltyReduction: 0,
			CriticalEffect:                nil,
		},
	}

	_, err := cannotCounterWithThisPower.AttackEffect.CounterAttackPenalty()
	checker.Assert(err, ErrorMatches, "power cannot counter, cannot calculate penalty")
}

func (suite *AttackingEffectCounterAttackPenaltyTest) TestAppliesPenaltyReduction(checker *C) {
	counterAttackingPower := &power.Power{
		Reference:    power.Reference{
			Name: "Static",
			ID:   "power0",
		},
		PowerType:    power.Physical,
		AttackEffect: &power.AttackingEffect{
			ToHitBonus:                    0,
			DamageBonus:                   0,
			CanCounterAttack:              true,
			CounterAttackPenaltyReduction: 2,
			CriticalEffect:                nil,
		},
	}

	counterAttackPenalty, err := counterAttackingPower.AttackEffect.CounterAttackPenalty()
	checker.Assert(err, IsNil)
	checker.Assert(counterAttackPenalty, Equals, 0)
}