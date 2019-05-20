package tag

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
)

type TagTestSuite struct {
	suite.Suite
	r http.Handler
}

func (suite *TagTestSuite) SetupSuite() {
	e := echo.New()
	e.GET("/v2/:org/:repo/tags/list", List)
	suite.r = e
}

func (suite *TagTestSuite) TestList() {
	assert := suite.Assert()
	require := suite.Require()

	req := httptest.NewRequest(http.MethodGet, "/v2/library/alpine/tags/list", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	expected := `{"name": "testOrg/testRepo", "tags": ["tag1", "tag2", "tag3", "tag4", "tag5"]}`
	assert.JSONEq(expected, string(b))
	assert.Equal("/v2/testOrg/testRepo/tags/list?n=5&last=tag5", rr.Header().Get(headerLink))
}

func TestBaseTestSuite(t *testing.T) {
	tests := new(TagTestSuite)
	suite.Run(t, tests)
}
