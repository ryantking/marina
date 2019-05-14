package repo

import (
	"testing"

	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/mocks"
	"github.com/stretchr/testify/suite"
	udb "upper.io/db.v3"
)

type RepoTestSuite struct {
	suite.Suite
	col *mocks.Collection
	res *mocks.Result
}

func (suite *RepoTestSuite) TearDownSuite() {
	getCollection = db.GetCollection
}

func (suite *RepoTestSuite) SetupTest() {
	suite.col = new(mocks.Collection)
	suite.res = new(mocks.Result)

	getCollection = func(name string) (udb.Collection, error) {
		return suite.col, nil
	}
}

func (suite *RepoTestSuite) TearDownTest() {
	suite.col.AssertExpectations(suite.T())
	suite.res.AssertExpectations(suite.T())
}

func (suite *RepoTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()

	suite.col.On("Find", "name", "testRepo").Return(suite.res)
	suite.res.On("And", "org_name", "testOrg").Return(suite.res)
	suite.res.On("Exists").Return(true, nil)

	exists, err := Exists("testRepo", "testOrg")
	require.NoError(err)
	assert.True(exists)
}

func TestRepoTestSuite(t *testing.T) {
	tests := new(RepoTestSuite)
	suite.Run(t, tests)
}
