package registry_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ryantking/marina/pkg/v2/routes"
	"github.com/stretchr/testify/suite"

	"github.com/emicklei/go-restful"
)

const (
	apiVersionPath = "/v2"
)

type RegistryTestSuite struct {
	suite.Suite
	container *restful.Container
}

func (suite *RegistryTestSuite) SetupSuite() {
	suite.container = restful.NewContainer()
	suite.container.Add(routes.Registry())
}

func (suite *RegistryTestSuite) TestAPIVersion() {
	assert := suite.Assert()
	require := suite.Require()

	req := httptest.NewRequest(http.MethodGet, apiVersionPath, nil)
	rr := httptest.NewRecorder()
	suite.container.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	assert.EqualValues("true", b)
}

func TestRegistryTestSuite(t *testing.T) {
	tests := new(RegistryTestSuite)
	suite.Run(t, tests)
}
