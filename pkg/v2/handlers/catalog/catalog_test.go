package catalog

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/web"
	"github.com/stretchr/testify/suite"
)

type CatalogTestSuite struct {
	suite.Suite
	r http.Handler
}

func (suite *CatalogTestSuite) SetupSuite() {
	e := echo.New()
	e.HTTPErrorHandler = web.ErrorHandler
	e.Binder = new(docker.Binder)
	e.GET("/v2/_catalog", Get)
	suite.r = e
}

func (suite *CatalogTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	getNames = func() ([]string, error) {
		return []string{"org1/repo1", "org1/repo2", "org2/repo3"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v2/_catalog", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	assert.JSONEq(`{"repositories": ["org1/repo1", "org1/repo2", "org2/repo3"]}`, string(b))
}

func (suite *CatalogTestSuite) TestGetPaginated() {
	assert := suite.Assert()
	require := suite.Require()

	getNamesPaginated = func(n uint, last string) ([]string, string, error) {
		assert.EqualValues(3, n)
		assert.Equal("org1/repo1", last)
		return []string{"org1/repo2", "org1/repo3", "org2/repo4"}, "org2/repo4", nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v2/_catalog?n=3&last=org1/repo1", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	assert.JSONEq(`{"repositories": ["org1/repo2", "org1/repo3", "org2/repo4"]}`, string(b))
	assert.Equal("/v2/_catalog?n=3&last=org2/repo4", rr.Header().Get(headerLink))
}

func TestCatalogTestSuite(t *testing.T) {
	tests := new(CatalogTestSuite)
	suite.Run(t, tests)
}
