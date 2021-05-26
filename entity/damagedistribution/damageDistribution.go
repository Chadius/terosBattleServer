package damagedistribution

// DamageDistribution tracks how damage is distributed.
type DamageDistribution struct {
	DamageAbsorbedByArmor   int
	DamageAbsorbedByBarrier int
	DamageDealt             int
	ExtraBarrierBurnt       int
	TotalBarrierBurnt       int
}
