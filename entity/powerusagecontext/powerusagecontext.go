package powerusagecontext

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/report"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// PowerUsageContext contains all of the information needed to calculate the effect of using a power.
type PowerUsageContext struct {
	SquaddieRepo		*squaddie.Repository
	ActingSquaddieID	string
	TargetSquaddieIDs	[]string

	PowerID				string
	PowerRepo			*power.Repository

	PowerForecast		*PowerForecast
	PowerReport			*report.PowerReport
}

// AttackContext holds the information needed to calculate expected damage.
type AttackContext struct {
	PowerID				string
	AttackerID			string
	TargetID			string
	IsCounterAttack		bool
}

// Clone returns a duplicate of the AttackContext.
func (context *AttackContext) Clone() *AttackContext {
	return &AttackContext{
		PowerID:           context.PowerID,
		AttackerID:        context.AttackerID,
		TargetID:          context.TargetID,
		IsCounterAttack: context.IsCounterAttack,
	}
}

// PowerForecast showcases the expected results of using a given power.
type PowerForecast struct {
	UserSquaddieID 		string
	PowerID        		string
	AttackEffectSummary []*AttackingPowerForecast
}

// AttackingPowerForecast gives a summary of the chance to hit and damage dealt by attacks.
type AttackingPowerForecast struct {
	AttackingSquaddieID				string
	PowerID							string
	TargetSquaddieID				string

	CriticalHitThreshold			int
	ChanceToHit						int
	DamageTaken						int
	HitRate							int
	BarrierDamageTaken				int

	//  Expected damage counts the number of 36ths so we can use integers for fractional math.
	ExpectedDamage					int
	ExpectedBarrierDamage			int
	ChanceToCritical				int
	CriticalExpectedDamage			int
	CriticalExpectedBarrierDamage	int

	CriticalDamageTaken				int
	CriticalBarrierDamageTaken		int

	IsACounterAttack				bool
	CounterAttack                 	*AttackingPowerForecast
}
