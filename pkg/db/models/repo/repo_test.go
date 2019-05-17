package repo

import (
	"os"
	"testing"

	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/engine"
	"github.com/romanyx/polluter"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/db"
	"github.com/ryantking/marina/pkg/db/models/org"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
	cleaner  dbcleaner.DbCleaner
	polluter *polluter.Polluter
}

func (suite *RepoTestSuite) SetupSuite() {
	cfg := config.Get().DB
	db := db.Get()

	db.AutoMigrate(&org.Model{}, &Model{})
	suite.cleaner = dbcleaner.New()
	mysql := engine.NewMySQLEngine(cfg.DSN)
	suite.cleaner.SetEngine(mysql)
	suite.polluter = polluter.New(polluter.MySQLEngine(db.DB()))
}

func (suite *RepoTestSuite) TearDownSuite() {
	db.Get().DropTable(&Model{})
}

func (suite *RepoTestSuite) SetupTest() {
	require := suite.Require()

	suite.cleaner.Acquire("organizations")
	suite.cleaner.Acquire("repositoriess")
	seed, err := os.Open("seed.yml")
	require.NoError(err)
	err = suite.polluter.Pollute(seed)
	require.NoError(err)
}

func (suite *RepoTestSuite) TearDownTest() {
	suite.cleaner.Clean("repositoriess")
	suite.cleaner.Clean("organizations")
}

func (suite *RepoTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()

	// suite.col.On("Find", "name", "testRepo").Return(suite.res)
	// suite.res.On("And", "org_name", "testOrg").Return(suite.res)
	// suite.res.On("Exists").Return(true, nil)

	exists, err := Exists("alpine", "library")
	require.NoError(err)
	assert.True(exists)
}

func TestRepoTestSuite(t *testing.T) {
	tests := new(RepoTestSuite)
	suite.Run(t, tests)
}
