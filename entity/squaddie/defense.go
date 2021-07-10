package squaddie

import (
	"github.com/cserrant/terosBattleServer/entity/damagedistribution"
)

// Defense holds everything needed to prevent the squaddie from getting hindered.
type Defense struct {
	CurrentHitPoints    		int                            `json:"current_hit_points" yaml:"current_hit_points"`
	MaxHitPoints        		int                            `json:"max_hit_points" yaml:"max_hit_points"`
	Dodge               		int                               `json:"dodge" yaml:"dodge"`
	Deflect             		int                             `json:"deflect" yaml:"deflect"`
	CurrentBarrier      		int                              `json:"current_barrier" yaml:"current_barrier"`
	MaxBarrier          		int                              `json:"max_barrier" yaml:"max_barrier"`
	Armor               		int                               `json:"armor" yaml:"armor"`
}

// SetHPToMax restores the Squaddie's HitPoints.
func (defense *Defense) SetHPToMax() {
	defense.CurrentHitPoints = defense.MaxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (defense *Defense) SetBarrierToMax() {
	defense.CurrentBarrier = defense.MaxBarrier
}

// ReduceHitPoints reduces the squaddie's HP, possibly killing them.
//   Hit Points cannot e reduced below 0.
func (defense *Defense) ReduceHitPoints(damage int) int {
	actualDamageTaken := damage
	if defense.CurrentHitPoints < damage {
		actualDamageTaken = defense.CurrentHitPoints
	}

	defense.CurrentHitPoints -= damage

	if defense.CurrentHitPoints < 0 {
		defense.CurrentHitPoints = 0
	}
	return actualDamageTaken
}

// ReduceBarrier reduces the squaddie's Barrier.
//   Barrier cannot e reduced below 0.
func (defense *Defense) ReduceBarrier(burn int) int {
	actualBarrierBurn := burn
	if defense.CurrentBarrier < burn {
		actualBarrierBurn = defense.CurrentBarrier
	}

	defense.CurrentBarrier -= burn

	if defense.CurrentBarrier < 0 {
		defense.CurrentBarrier = 0
	}

	return actualBarrierBurn
}

// TakeDamageDistribution reduces the Squaddie's Barrier and Hit Points based on the distribution.
func (defense *Defense) TakeDamageDistribution(distribution *damagedistribution.DamageDistribution) {
	actualBarrierBurn := defense.ReduceBarrier(distribution.DamageAbsorbedByBarrier)
	actualDamageTaken := defense.ReduceHitPoints(distribution.RawDamageDealt)

	distribution.ActualBarrierBurn = actualBarrierBurn
	distribution.ActualDamageTaken = actualDamageTaken
}

// IsDead returns true if the squaddie has died
func (defense *Defense) IsDead() bool {
	return defense.CurrentHitPoints <= 0
}

// GainHitPoints heals the squaddie and returns the number of hit points healed.
func (defense *Defense) GainHitPoints(hitPoints int) int {
	actualHealingReceived := hitPoints
	if defense.CurrentHitPoints + actualHealingReceived >= defense.MaxHitPoints {
		actualHealingReceived = defense.MaxHitPoints - defense.CurrentHitPoints
	}
	defense.CurrentHitPoints += actualHealingReceived
	return actualHealingReceived
}
