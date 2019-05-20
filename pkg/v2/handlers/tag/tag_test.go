package tag

import (
	"net/http"
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

// func (suite *TagTestSuite) TestParsePath() {
// 	assert := suite.Assert()
// 	require := suite.Require()
//
// 	c := new(mocks.Context)
// 	c.On("Param", "repo").Return("testRepo")
// 	c.On("Param", "org").Return("testOrg")
// 	repoExists = func(repoName, orgName string) (bool, error) {
// 		return true, nil
// 	}
//
// 	repoName, orgName, err := parsePath(c)
// 	require.NoError(err)
// 	assert.Equal("testRepo", repoName)
// 	assert.Equal("testOrg", orgName)
// 	c.AssertExpectations(suite.T())
// }
//
// func (suite *TagTestSuite) TestParsePathNonExistentRepo() {
// 	assert := suite.Assert()
//
// 	c := new(mocks.Context)
// 	c.On("Param", "repo").Return("testRepo")
// 	c.On("Param", "org").Return("testOrg")
// 	c.On("Set", "docker_err_code", docker.CodeNameUnknown)
// 	repoExists = func(repoName, orgName string) (bool, error) {
// 		return false, nil
// 	}
//
// 	_, _, err := parsePath(c)
// 	assert.Equal(echo.NewHTTPError(http.StatusNotFound, "no such repository"), err)
// 	c.AssertExpectations(suite.T())
// }
//
// func (suite *TagTestSuite) TestParsePathExistsError() {
// 	assert := suite.Assert()
//
// 	c := new(mocks.Context)
// 	c.On("Param", "repo").Return("testRepo")
// 	c.On("Param", "org").Return("testOrg")
// 	repoExists = func(repoName, orgName string) (bool, error) {
// 		return false, fmt.Errorf("test error")
// 	}
//
// 	_, _, err := parsePath(c)
// 	assert.EqualError(errors.Cause(err), "test error")
// 	c.AssertExpectations(suite.T())
// }
//
// func (suite *TagTestSuite) TestList() {
// 	assert := suite.Assert()
// 	require := suite.Require()
//
// 	repoExists = func(repoName, orgName string) (bool, error) {
// 		assert.Equal("testRepo", repoName)
// 		assert.Equal("testOrg", orgName)
// 		return true, nil
// 	}
// 	tagList = func(repoName, orgName string, pageSize uint, last string) ([]string, string, error) {
// 		assert.Equal("testRepo", repoName)
// 		assert.Equal("testOrg", orgName)
// 		assert.EqualValues(5, pageSize)
// 		assert.Zero(last)
// 		return []string{"tag1", "tag2", "tag3", "tag4", "tag5"}, "tag5", nil
// 	}
//
// 	req := httptest.NewRequest(http.MethodGet, "/v2/testOrg/testRepo/tags/list?n=5", nil)
// 	rr := httptest.NewRecorder()
// 	suite.r.ServeHTTP(rr, req)
// 	require.Equal(http.StatusOK, rr.Code)
// 	b, err := ioutil.ReadAll(rr.Body)
// 	require.NoError(err)
// 	expected := `{"name": "testOrg/testRepo", "tags": ["tag1", "tag2", "tag3", "tag4", "tag5"]}`
// 	assert.JSONEq(expected, string(b))
// 	assert.Equal("/v2/testOrg/testRepo/tags/list?n=5&last=tag5", rr.Header().Get(headerLink))
// }

func TestBaseTestSuite(t *testing.T) {
	tests := new(TagTestSuite)
	suite.Run(t, tests)
}
