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
	CanBeEquipped             bool `json:"can_be_equipped" yaml:"can_be_equipped"`
	CanCounterAttack          bool `json:"can_counter_attack" yaml:"can_counter_attack"`
	CounterAttackToHitPenalty int  `json:"counter_attack_penalty" yaml:"counter_attack_penalty"`
	CriticalEffect *CriticalEffect `json:"critical_effect" yaml:"critical_effect"`
}

// NewPower generates a Power with default values.
func NewPower(name string) *Power {
	newAttackingPower := Power{
		Reference: Reference{
			Name: name,
			ID:   utility.StringWithCharset(8, "abcdefgh0123456789"),
		},
		PowerType: Physical,
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

// CanCriticallyHit returns true if this power is capable of critically hitting for additional effects.
func (p *Power) CanCriticallyHit() bool {
	return p.AttackEffect != nil && p.AttackEffect.CriticalEffect != nil
}