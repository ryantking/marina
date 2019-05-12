package blob

import (
	"testing"

	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/mocks"
	"github.com/stretchr/testify/suite"
	udb "upper.io/db.v3"
)

type BlobTestSuite struct {
	suite.Suite
	col *mocks.Collection
	res *mocks.Result
}

func (suite *BlobTestSuite) TearDownSuite() {
	getCollection = db.GetCollection
}

func (suite *BlobTestSuite) SetupTest() {
	suite.col = new(mocks.Collection)
	suite.res = new(mocks.Result)

	getCollection = func(name string) (udb.Collection, error) {
		return suite.col, nil
	}
}

func (suite *BlobTestSuite) TearDownTest() {
	suite.col.AssertExpectations(suite.T())
	suite.res.AssertExpectations(suite.T())
}

func (suite *BlobTestSuite) TestCreate() {
	assert := suite.Assert()

	expected := Model{
		Digest:   "testDigest",
		RepoName: "testRepo",
		OrgName:  "testOrg",
	}
	suite.col.On("Insert", &expected).Return(nil, nil)

	err := Create("testDigest", "testRepo", "testOrg")
	assert.NoError(err)
}

func (suite *BlobTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()

	suite.col.On("Find", "digest", "testDigest").Return(suite.res)
	suite.res.On("Exists").Return(true, nil)

	exists, err := Exists("testDigest")
	require.NoError(err)
	assert.True(exists)
}

func (suite *BlobTestSuite) TestDelete() {
	assert := suite.Assert()

	suite.col.On("Find", "digest", "testDigest").Return(suite.res)
	suite.res.On("Delete").Return(nil)

	err := Delete("testDigest")
	assert.NoError(err)
}

func TestBlobTestSuite(t *testing.T) {
	tests := new(BlobTestSuite)
	suite.Run(t, tests)
}
