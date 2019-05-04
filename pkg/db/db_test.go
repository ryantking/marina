package db

import (
	"testing"

	"github.com/ryantking/marina/pkg/config"
	"github.com/stretchr/testify/suite"

	// Sqlite driver
	_ "upper.io/db.v3/mysql"
)

type DBTestSuite struct {
	suite.Suite
}

func (suite *DBTestSuite) TearDownTest() {
	require := suite.Require()

	config.Destroy()
	err := Close()
	require.NoError(err)
}

func (suite *DBTestSuite) TestGetDB() {
	require := suite.Require()

	db, err := Get()
	require.NoError(err)
	err = db.Ping()
	require.NoError(err)
}

func (suite *DBTestSuite) TestDBTimeout() {
	assert := suite.Assert()
	require := suite.Require()

	config.Set("db.dsn", "foo:bar@tcp(badurl:3306)/db")
	config.Set("db.timeout", "1ms")

	db, err := Get()
	require.Error(err)
	assert.Nil(db)
}

func TestDBTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tests := new(DBTestSuite)
	suite.Run(t, tests)
}
