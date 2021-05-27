package powerusagescenario

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// Setup is supplied upon creation to explain all of the relevant parts of this power.
type Setup struct {
	UserID          string
	PowerID         string
	Targets         []string
	IsCounterAttack bool
}

// RepositoryCollection holds all of the repositories used in the setup.
type RepositoryCollection struct {
	SquaddieRepo    *squaddie.Repository
	PowerRepo       *power.Repository
}