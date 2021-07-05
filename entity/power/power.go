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
	HealingEffect *HealingEffect  `json:"healing_effect" yaml:"healing_effect"`
}

// GetReference returns a new PowerReference.
func (p Power) GetReference() *Reference {
	return &Reference{
		Name: p.Name,
		ID:   p.ID,
	}
}

// NewPower generates a Power with default values.
func NewPower(name string) *Power {
	newAttackingPower := Power{
		Reference: Reference{
			Name: name,
			ID:   "power_" + utility.StringWithCharset(8, "abcdefgh0123456789"),
		},
		PowerType: Physical,
	}
	return &newAttackingPower
}

// CheckPowerForErrors verifies the Power's fields and raises an error if it's invalid.
func CheckPowerForErrors(newPower *Power) (newError error) {
	if newPower.PowerType != Physical &&
		newPower.PowerType != Spell {
		return fmt.Errorf("AttackingPower '%s' has unknown power_type: '%s'", newPower.Name, newPower.PowerType)
	}

	return nil
}
