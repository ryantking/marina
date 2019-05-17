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
	require := suite.Require()
	cfg := config.Get().DB
	db := db.Get()

	db.AutoMigrate(&Model{})
	suite.cleaner = dbcleaner.New()
	mysql := engine.NewMySQLEngine(cfg.DSN)
	suite.cleaner.SetEngine(mysql)
	p := polluter.New(polluter.MySQLEngine(db.DB()))
	seed, err := os.Open("seed.yml")
	require.NoError(err)
	err = p.Pollute(seed)
	require.NoError(err)
}

func (suite *OrgTestSuite) TearDownSuite() {
	db.Get().DropTable(&Model{})
}

func (suite *OrgTestSuite) SetupTest() {
	suite.cleaner.Acquire("organization")
}

func (suite *OrgTestSuite) TearDownTest() {
	suite.cleaner.Clean("organization")
}

func (suite *OrgTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		name   string
		exists bool
	}{{"library", true}, {"mysql", true}, {"minio", false}}

	for _, tt := range tests {
		exists, err := Exists(tt.name)
		require.NoError(err, "For: %s", tt.name)
		assert.Equal(tt.exists, exists, "For: %s", tt.name)
	}
}

func TestOrgTestSuite(t *testing.T) {
	tests := new(OrgTestSuite)
	suite.Run(t, tests)
}
