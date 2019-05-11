package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/web"
	"github.com/ryantking/marina/pkg/web/mocks"
	"github.com/stretchr/testify/suite"
)

type ManifestTestSuite struct {
	suite.Suite
	r http.Handler
}

func (suite *ManifestTestSuite) SetupSuite() {
	e := echo.New()
	e.HTTPErrorHandler = web.ErrorHandler
	e.Binder = new(docker.Binder)
	e.HEAD("/v2/:org/:repo/manifests/:ref", Exists)
	e.GET("/v2/:org/:repo/manifests/:ref", Get)
	e.PUT("/v2/:org/:repo/manifests/:ref", Update)
	e.DELETE("/v2/:org/:repo/manifests/:ref", Delete)
	suite.r = e
}

func (suite *ManifestTestSuite) TestParsePath() {
	assert := suite.Assert()
	require := suite.Require()

	c := new(mocks.Context)
	c.On("Param", "ref").Return("testTag")
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}

	ref, repoName, orgName, err := parsePath(c)
	require.NoError(err)
	assert.Equal("testTag", ref)
	assert.Equal("testRepo", repoName)
	assert.Equal("testOrg", orgName)
	c.AssertExpectations(suite.T())
}

func (suite *ManifestTestSuite) TestParsePathNonExistentRepo() {
	assert := suite.Assert()

	c := new(mocks.Context)
	c.On("Param", "ref").Return("testTag")
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	c.On("Set", "docker_err_code", docker.CodeNameUnknown)
	repoExists = func(repoName, orgName string) (bool, error) {
		return false, nil
	}

	_, _, _, err := parsePath(c)
	assert.Equal(echo.NewHTTPError(http.StatusNotFound, "no such repository"), err)
	c.AssertExpectations(suite.T())
}

func (suite *ManifestTestSuite) TestParsePathExistsError() {
	assert := suite.Assert()

	c := new(mocks.Context)
	c.On("Param", "ref").Return("testTag")
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	repoExists = func(repoName, orgName string) (bool, error) {
		return false, fmt.Errorf("test error")
	}

	_, _, _, err := parsePath(c)
	assert.EqualError(errors.Cause(err), "test error")
	c.AssertExpectations(suite.T())
}

func (suite *ManifestTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()

	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	getManifest = func(ref, repoName, orgName string) (docker.Manifest, string, error) {
		man := new(docker.ManifestV2)
		man.Config.Digest = "testdigest"
		return man, docker.MIMEManifestV2, nil
	}

	req := httptest.NewRequest(http.MethodHead, "/v2/testOrg/testRepo/manifests/testTag", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	assert.Equal("testdigest", rr.Header().Get(docker.HeaderContentDigest))
	assert.NotZero(rr.Header().Get(echo.HeaderContentLength))
}

func (suite *ManifestTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	getManifest = func(ref, repoName, orgName string) (docker.Manifest, string, error) {
		man := new(docker.ManifestV2)
		man.Config.Digest = "testdigest"
		return man, docker.MIMEManifestV2, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v2/testOrg/testRepo/manifests/testTag", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	assert.Equal(docker.MIMEManifestV2, rr.Header().Get(echo.HeaderContentType))
	assert.Equal("testdigest", rr.Header().Get(docker.HeaderContentDigest))
}

func (suite *ManifestTestSuite) TestUpdate() {
	assert := suite.Assert()
	require := suite.Require()

	man := new(docker.ManifestV2)
	man.Config.Digest = "testdigest"
	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	updateManifest = func(ref, repoName, orgName string, manifest docker.Manifest, manifestType string) error {
		assert.Equal("testTag", ref)
		assert.Equal("testRepo", repoName)
		assert.Equal("testOrg", orgName)
		assert.Equal(man, manifest)
		assert.Equal(docker.MIMEManifestV2, manifestType)
		return nil
	}

	b, err := json.Marshal(man)
	require.NoError(err)
	buf := bytes.NewBuffer(b)
	req := httptest.NewRequest(http.MethodPut, "/v2/testOrg/testRepo/manifests/testTag", buf)
	req.Header.Set(echo.HeaderContentType, docker.MIMEManifestV2)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusCreated, rr.Code)
	assert.Equal("testdigest", rr.Header().Get(docker.HeaderContentDigest))
	assert.Equal("/v2/testOrg/testRepo/manifests/testTag", rr.Header().Get(echo.HeaderLocation))
}

func (suite *ManifestTestSuite) TestDelete() {
	assert := suite.Assert()

	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	deleteManifest = func(ref string) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/v2/testOrg/testRepo/manifests/testTag", nil)
	req.Header.Set(echo.HeaderContentType, docker.MIMEManifestV2)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	assert.Equal(http.StatusAccepted, rr.Code)
}

func TestManifestTestSuite(t *testing.T) {
	tests := new(ManifestTestSuite)
	suite.Run(t, tests)
}
