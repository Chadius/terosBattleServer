package power

import (
	"fmt"

	"github.com/cserrant/terosBattleServer/utility"
)

// Reference is used to identify a power and is used to quickly identify a power.
type Reference struct {
	Name string `json:"name" yaml:"name"`
	ID   string `json:"id" yaml:"id"`
}

// Type defines the expected sources the power could be conjured from.
type Type string

const (
	// Physical powers use martial training and cunning. Examples: Swords, Bows, Pushing
	Physical = "Physical"
	// Spell powers are magical in nature and conjured without tools. Examples: Fireball, Mindread
	Spell = "Spell"
)

// Power are the abilities every Squaddie can use. These range from dealing damage, to opening doors, to healing.
type Power struct {
	Reference `yaml:",inline"`
	PowerType Type `json:"power_type" yaml:"power_type"`
	AttackEffect *AttackingEffect  `json:"attack_effect" yaml:"attack_effect"`
}

// GetReference returns a new PowerReference.
func (p Power) GetReference() *Reference {
	return &Reference{
		Name: p.Name,
		ID:   p.ID,
	}
}

// AttackingEffect is a power designed to deal damage.
type AttackingEffect struct {
	ToHitBonus                int  `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	DamageBonus               int  `json:"damage_bonus" yaml:"damage_bonus"`
	ExtraBarrierBurn          int  `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	CriticalHitThreshold      int  `json:"critical_hit_threshold" yaml:"critical_hit_threshold"`
	CanBeEquipped             bool `json:"can_be_equipped" yaml:"can_be_equipped"`
	CanCounterAttack          bool `json:"can_counter_attack" yaml:"can_counter_attack"`
	CounterAttackToHitPenalty int  `json:"counter_attack_penalty" yaml:"counter_attack_penalty"`
}

// NewPower generates a Power with default values.
func NewPower(name string) *Power {
	newAttackingPower := Power{
		Reference: Reference{
			Name: name,
			ID:   utility.StringWithCharset(8, "abcdefgh0123456789"),
		},
		PowerType: Physical,
		AttackEffect: &AttackingEffect{
			ToHitBonus:           0,
			DamageBonus:          0,
			ExtraBarrierBurn:     0,
			CriticalHitThreshold: 0,
		},
	}
	return &newAttackingPower
}

// CheckPowerForErrors verifies the Power's fields and raises an error if it's invalid.
func CheckPowerForErrors(newPower *Power) (newError error) {
	return checkAttackingEffectForErrors(newPower)
}

func checkAttackingEffectForErrors(newAttackingPower *Power) (newError error) {
	if newAttackingPower.PowerType != Physical &&
		newAttackingPower.PowerType != Spell {
		return fmt.Errorf("AttackingPower '%s' has unknown power_type: '%s'", newAttackingPower.Name, newAttackingPower.PowerType)
	}

	return nil
}

// GetChanceToHitBasedOnHitRate is a smarter look up table.
func GetChanceToHitBasedOnHitRate(toHitBonus int) (chanceOutOf36 int) {
	if toHitBonus > 4 {
		return 36
	}

	if toHitBonus < -5 {
		return 0
	}

	toHitChanceReference := map[int]int{
		4:  35,
		3:  33,
		2:  30,
		1:  26,
		0:  21,
		-1: 15,
		-2: 10,
		-3: 6,
		-4: 3,
		-5: 1,
	}

	return toHitChanceReference[toHitBonus]
}

// GetChanceToCriticalBasedOnThreshold is a smarter look up table.
func GetChanceToCriticalBasedOnThreshold(criticalThreshold int) (chanceOutOf36 int) {
	if criticalThreshold > 11 {
		return 36
	}

	if criticalThreshold < 2 {
		return 0
	}

	criticalChanceByCriticalThreshold := map[int]int{
		11: 35,
		10: 33,
		9:  30,
		8:  26,
		7:  21,
		6:  15,
		5:  10,
		4:  6,
		3:  3,
		2:  1,
	}

	return criticalChanceByCriticalThreshold[criticalThreshold]
}
