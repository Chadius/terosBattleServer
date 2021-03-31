package squaddieclass_test

import (
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ClassRepositoryUnmarshalSuite struct {
	repo *squaddieclass.Repository
	jsonByteStream []byte
}

var _ = Suite(&ClassRepositoryUnmarshalSuite{})

func (suite *ClassRepositoryUnmarshalSuite) SetUpTest(checker *C) {
	suite.repo = squaddieclass.NewRepository()
	suite.jsonByteStream = []byte(`[{
		"ID": "aaaaaaaa",
		"name": "Mage"
	}]`)
}

func (suite *ClassRepositoryUnmarshalSuite) TestLoadClassesWithJSON(checker *C) {
	checker.Assert(suite.repo.GetNumberOfClasses(), Equals, 0)
	success, _ := suite.repo.AddJSONSource(suite.jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.repo.GetNumberOfClasses(), Equals, 1)
}

func (suite *ClassRepositoryUnmarshalSuite) TestLoadClassesDirectly(checker *C) {
	listOfClasses := []*squaddieclass.Class{
		{
			ID:                "class1",
			Name:              "Mage",
			BaseClassRequired: false,
		},
		{
			ID:                "class2",
			Name:              "Dimension Walker",
			BaseClassRequired: true,
		},
	}
	checker.Assert(suite.repo.GetNumberOfClasses(), Equals, 0)
	success, _ := suite.repo.AddListOfClasses(listOfClasses)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.repo.GetNumberOfClasses(), Equals, 2)
}

type ClassRepositoryRetrieveSuite struct {
	repo *squaddieclass.Repository
	mageClass *squaddieclass.Class
	dimensionWalkerClass *squaddieclass.Class
}

var _ = Suite(&ClassRepositoryRetrieveSuite{})

func (suite *ClassRepositoryRetrieveSuite) SetUpTest(checker *C) {
	suite.repo = squaddieclass.NewRepository()
	suite.mageClass = &squaddieclass.Class{
		ID:                "class0",
		Name:              "Mage",
		BaseClassRequired: false,
	}
	suite.dimensionWalkerClass = &squaddieclass.Class{
		ID:                "class1",
		Name:              "Dimension Walker",
		BaseClassRequired: true,
	}
	suite.repo.AddListOfClasses([]*squaddieclass.Class{suite.mageClass, suite.dimensionWalkerClass})
}

func (suite *ClassRepositoryRetrieveSuite) TestGetClassByID(checker *C) {
	foundClass, err := suite.repo.GetClassByID(suite.mageClass.ID)
	checker.Assert(err, IsNil)
	checker.Assert(foundClass.ID, Equals, suite.mageClass.ID)
	checker.Assert(foundClass.Name, Equals, suite.mageClass.Name)
}

func (suite *ClassRepositoryRetrieveSuite) TestRaiseErrorWhenClassDoesNotExist(checker *C) {
	_, err := suite.repo.GetClassByID("bad ID")
	checker.Assert(err, ErrorMatches,`class repository: No class found with ID: "bad ID"`)
}
