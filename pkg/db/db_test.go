package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/ryantking/marina/pkg/config"
	"github.com/stretchr/testify/suite"

	// Sqlite driver
	_ "upper.io/db.v3/mysql"
)

type DBTestSuite struct {
	suite.Suite
}

func (suite *DBTestSuite) TearDownTest() {
	config.Destroy()
	Close()
}

func (suite *DBTestSuite) TestGetDB() {
	assert := suite.Assert()
	require := suite.Require()

	var db *gorm.DB
	require.NotPanics(func() {
		db = Get()
	})
	assert.NotNil(db)
	err := db.DB().Ping()
	require.NoError(err)
}

func TestDBTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tests := new(DBTestSuite)
	suite.Run(t, tests)
}
