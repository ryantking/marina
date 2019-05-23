package catalog

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/ryantking/marina/pkg/web"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

type CatalogTestSuite struct {
	suite.Suite
	r       http.Handler
	cleaner dbcleaner.DbCleaner
}

func (suite *CatalogTestSuite) SetupSuite() {
	e := echo.New()
	e.HTTPErrorHandler = web.ErrorHandler
	e.Binder = new(docker.Binder)
	e.GET("/v2/_catalog", Get)
	suite.r = e
	mysql := engine.NewMySQLEngine(config.Get().DB.DSN)
	suite.cleaner = dbcleaner.New()
	suite.cleaner.SetEngine(mysql)
}

func (suite *CatalogTestSuite) SetupTest() {
	suite.cleaner.Acquire("Repository", "Organization")
	testutil.Clear(context.Background())
	testutil.Seed(context.Background())
}

func (suite *CatalogTestSuite) TearDownTest() {
	suite.cleaner.Clean("Repository", "Organization")
}

func (suite *CatalogTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	req := httptest.NewRequest(http.MethodGet, "/v2/_catalog", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	expected := `{"repositories": ["library/alpine", "library/nginx", "library/redis", "mysql/mysql", "mysql/mysql-client"]}`
	assert.JSONEq(expected, string(b))
}

func (suite *CatalogTestSuite) TestGetPaginated() {
	assert := suite.Assert()
	require := suite.Require()

	req := httptest.NewRequest(http.MethodGet, "/v2/_catalog?n=2&last=library/nginx", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	assert.JSONEq(`{"repositories": ["library/redis", "mysql/mysql"]}`, string(b))
	assert.Equal("/v2/_catalog?n=2&last=mysql/mysql", rr.Header().Get(headerLink))
}

func TestCatalogTestSuite(t *testing.T) {
	tests := new(CatalogTestSuite)
	suite.Run(t, tests)
}
