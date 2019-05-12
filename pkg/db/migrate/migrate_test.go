package migrate

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ryantking/marina/pkg/config"
	"github.com/stretchr/testify/suite"

	// MySQL migrations driver
	_ "github.com/golang-migrate/migrate/database/mysql"
)

type MigrateTestSuite struct {
	suite.Suite
}

func (suite *MigrateTestSuite) TearDownTest() {
	config.Destroy()
}

func (suite *MigrateTestSuite) TestGetLatest() {
	require := suite.Require()

	wd, err := os.Getwd()
	require.NoError(err)
	migrationDir := fmt.Sprintf("%s/migrations-test", wd)
	err = os.Mkdir(migrationDir, 0777)
	require.NoError(err)
	config.Set("db.migrations_dir", migrationDir)
	defer os.RemoveAll(migrationDir)

	testData := []byte("test migration\n")
	for i := 1; i <= 5; i++ {
		fname := fmt.Sprintf("%s/%d_testmigration.up.sql", migrationDir, i)
		err = ioutil.WriteFile(fname, testData, 0777)
		require.NoError(err)
		fname = fmt.Sprintf("%s/%d_testmigration.down.sql", migrationDir, i)
		err = ioutil.WriteFile(fname, testData, 0777)
		require.NoError(err)
	}

	latest, err := GetLatest()
	require.NoError(err)
	require.Equal(5, latest)
}

func (suite *MigrateTestSuite) TestMigrations() {
	require := suite.Require()

	config.Set("db.migrations_dir", "../../../migrations")
	latest, err := GetLatest()
	require.NoError(err)

	cfg := config.Get()
	for i := 1; i <= latest; i++ {
		err := To(cfg.DB.Type, cfg.DB.DSN, uint(i))
		require.NoError(err)
		err = Down(cfg.DB.Type, cfg.DB.DSN)
		require.NoError(err)
	}
}

func TestMigrateTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tests := new(MigrateTestSuite)
	suite.Run(t, tests)
}
