package base

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
	r http.Handler
}

func (suite *BaseTestSuite) SetupSuite() {
	e := echo.New()
	e.GET("/v2/", Get)
	suite.r = e
}

func (suite *BaseTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	req := httptest.NewRequest(http.MethodGet, "/v2/", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	assert.EqualValues("true", b)
}

func TestBaseTestSuite(t *testing.T) {
	tests := new(BaseTestSuite)
	suite.Run(t, tests)
}
