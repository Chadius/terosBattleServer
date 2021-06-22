package powerusagescenario

// Setup is supplied upon creation to explain all of the relevant parts of this power.
type Setup struct {
	UserID          string
	PowerID         string
	Targets         []string
	IsCounterAttack bool
}

