package squaddie_test

import (
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ClassProgressSuite struct{}

var _ = Suite(&ClassProgressSuite{})

func (suite *ClassProgressSuite) TestIfLevelIsConsumed(checker *C) {
	progress := &squaddie.ClassProgress{
		ClassID:        "007",
		ClassName:      "Superspy",
		LevelsConsumed: []string{"1", "2", "3", "4", "5"},
	}
	checker.Assert(progress.IsLevelAlreadyConsumed("1"), Equals, true)
	checker.Assert(progress.IsLevelAlreadyConsumed("10"), Equals, false)
}
