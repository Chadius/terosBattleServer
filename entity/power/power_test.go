package power_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerChanceCheckSuite struct{}

var _ = Suite(&PowerChanceCheckSuite{})

func (suite *PowerChanceCheckSuite) TestGetChanceToHitBasedOnHitRate(checker *C) {
	checker.Assert(power.GetChanceToHitBasedOnHitRate(9001), Equals,36)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(5), Equals,36)

	checker.Assert(power.GetChanceToHitBasedOnHitRate(-6), Equals,0)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(-9001), Equals,0)

	checker.Assert(power.GetChanceToHitBasedOnHitRate(4), Equals,35)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(3), Equals,33)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(2), Equals,30)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(1), Equals,26)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(0), Equals,21)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(-1), Equals,15)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(-2), Equals,10)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(-3), Equals,6)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(-4), Equals,3)
	checker.Assert(power.GetChanceToHitBasedOnHitRate(-5), Equals,1)
}
