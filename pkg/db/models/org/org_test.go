package org

import (
	"testing"

	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/mocks"
	"github.com/stretchr/testify/suite"
	udb "upper.io/db.v3"
)

type OrgTestSuite struct {
	suite.Suite
	col *mocks.Collection
	res *mocks.Result
}

func (suite *OrgTestSuite) TearDownSuite() {
	getCollection = db.GetCollection
}

func (suite *OrgTestSuite) SetupTest() {
	suite.col = new(mocks.Collection)
	suite.res = new(mocks.Result)

	getCollection = func(name string) (udb.Collection, error) {
		return suite.col, nil
	}
}

func (suite *OrgTestSuite) TearDownTest() {
	suite.col.AssertExpectations(suite.T())
	suite.res.AssertExpectations(suite.T())
}

func (suite *OrgTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()

	suite.col.On("Find", "name", "testOrg").Return(suite.res)
	suite.res.On("Exists").Return(true, nil)

	exists, err := Exists("testOrg")
	require.NoError(err)
	assert.True(exists)
}

func TestOrgTestSuite(t *testing.T) {
	tests := new(OrgTestSuite)
	suite.Run(t, tests)
}
