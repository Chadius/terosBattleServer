package damagedistribution

// ToHitComparison describes how much to hit bonus and penalty are involved in an attack.
type ToHitComparison struct {
	ToHitBonus int
	AttackerToHitBonus int
	DefenderToHitPenalty int
}

// DamageDistribution tracks how a Squaddie handles damage across their Barrier, Armor and Health.
type DamageDistribution struct {
	DamageAbsorbedByArmor   int
	DamageAbsorbedByBarrier int
	RawDamageDealt          int
	ExtraBarrierBurnt       int
	TotalRawBarrierBurnt    int
	IsFatalToTarget         bool
	ActualBarrierBurn       int
	ActualDamageTaken       int
}
