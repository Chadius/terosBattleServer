package power

// HealingEffect is a power designed to restore hit points and cure ailments.
type HealingEffect struct {
	HealingAdjustmentBasedOnUserMind HealingAdjustmentBasedOnUserMind `json:"healing_adjustment_based_on_user_mind" yaml:"healing_adjustment_based_on_user_mind"`
	HitPointsHealed                   int             `json:"hit_points_healed" yaml:"hit_points_healed"`
}

// HealingAdjustmentBasedOnUserMind indicates how much the user's Mind should be adjusted.
type HealingAdjustmentBasedOnUserMind string
const (
	Full HealingAdjustmentBasedOnUserMind = "full"
	Half HealingAdjustmentBasedOnUserMind = "half"
	Zero HealingAdjustmentBasedOnUserMind = "zero"
)
