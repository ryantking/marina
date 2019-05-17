package org

import (
	"os"
	"testing"

	"github.com/romanyx/polluter"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/db"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

type OrgTestSuite struct {
	suite.Suite
	cleaner dbcleaner.DbCleaner
}

func (suite *OrgTestSuite) SetupSuite() {
	cfg := config.Get().DB
	db := db.Get()

	db.AutoMigrate(&Org{}, &Repo{})
	suite.cleaner = dbcleaner.New()
	mysql := engine.NewMySQLEngine(cfg.DSN)
	suite.cleaner.SetEngine(mysql)
}

func (suite *OrgTestSuite) TearDownSuite() {
	db.Get().DropTable(&Repo{}, &Org{})
}

func (suite *OrgTestSuite) SetupTest() {
	require := suite.Require()
	suite.cleaner.Acquire("organizations")
	suite.cleaner.Acquire("repositories")

	p := polluter.New(polluter.MySQLEngine(db.Get().DB()))
	seed, err := os.Open("seed.yml")
	require.NoError(err)
	err = p.Pollute(seed)
	require.NoError(err)
}

func (suite *OrgTestSuite) TearDownTest() {
	suite.cleaner.Clean("repositories")
	suite.cleaner.Clean("organizations")
}

func (suite *OrgTestSuite) TestOrgExists() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		name   string
		exists bool
	}{{"library", true}, {"mysql", true}, {"minio", false}}

	for _, tt := range tests {
		exists, err := OrgExists(tt.name)
		require.NoError(err, "For: %s", tt.name)
		assert.Equal(tt.exists, exists, "For: %s", tt.name)
	}
}

func (suite *OrgTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		repoName string
		orgName  string
		exists   bool
	}{{"alpine", "library", true}, {"mysql", "mysql", true}, {"mysql", "library", false}}

	for _, tt := range tests {
		exists, err := Exists(tt.repoName, tt.orgName)
		require.NoError(err, "For: %s/%s", tt.orgName, tt.repoName)
		assert.Equal(tt.exists, exists, "For: %s/%s", tt.orgName, tt.repoName)
	}
}

func TestOrgTestSuite(t *testing.T) {
	tests := new(OrgTestSuite)
	suite.Run(t, tests)
}
