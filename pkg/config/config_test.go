package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	env map[string]string
}

func (suite *ConfigTestSuite) SetupSuite() {
	vars := []string{
		"ROOT_URL", "ENVIRONMENT", "VERSION", "SERVER_ADDRESS", "SERVER_WRITE_TIMEOUT", "SERVER_READ_TIMEOUT",
		"DB_DSN", "DB_TYPE", "DB_TIMEOUT", "DB_MAX_CONNS", "DB_MIGRATIONS_DIR", "SWAGGER_URL", "SWAGGER_DOCS_PATH",
		"S3_ENDPOINT", "S3_ACCESS_KEY_ID", "S3_SECRET_ACCESS_KEY",
	}

	suite.env = map[string]string{}
	for _, key := range vars {
		suite.env[key] = os.Getenv(key)
	}

	Version = "dev"
}

func (suite *ConfigTestSuite) TearDownTest() {
	Version = "dev"
	for key, value := range suite.env {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
	Destroy()
}

func (suite *ConfigTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	require.NotPanics(func() {
		Get()
	})
	cfg1 := Get()
	cfg2 := Get()
	assert.Equal(cfg1, cfg2)
}

func (suite *ConfigTestSuite) TestGetDefaults() {
	assert := suite.Assert()

	cfg := Get()
	assert.Equal("http://localhost:8080", cfg.RootURL)
	assert.Equal("DEVELOPMENT", cfg.Environment)
	assert.Equal(Version, cfg.Version)
	assert.Equal("0.0.0.0:8080", cfg.Server.Address)
	assert.Equal(time.Minute*5, cfg.Server.WriteTimeout)
	assert.Equal(time.Minute*5, cfg.Server.ReadTimeout)
	assert.Equal("marina:marina@tcp(localhost:3306)/marinatest", cfg.DB.DSN)
	assert.Equal("mysql", cfg.DB.Type)
	assert.Equal(time.Minute, cfg.DB.Timeout)
	assert.Equal(200, cfg.DB.MaxConns)
	assert.Equal("migrations", cfg.DB.MigrationsDir)
	assert.Equal("https://swagger.cogolo.net", cfg.Swagger.URL)
	assert.Equal("/api/docs", cfg.Swagger.DocsPath)
	assert.Equal("localhost:9000", cfg.S3.Endpoint)
	assert.Equal("minio", cfg.S3.AccessKeyID)
	assert.Equal("minio123", cfg.S3.SecretAccessKey)
}

func (suite *ConfigTestSuite) TestGetEnvironment() {
	assert := suite.Assert()

	os.Setenv("ROOT_URL", "test")
	cfg := Get()
	assert.Equal("test", cfg.RootURL)
}

func (suite *ConfigTestSuite) TestSet() {
	assert := suite.Assert()

	Set("root_url", "test")
	cfg := Get()
	assert.Equal("test", cfg.RootURL)
}

func TestConfigTestSuite(t *testing.T) {
	tests := new(ConfigTestSuite)
	suite.Run(t, tests)
}
