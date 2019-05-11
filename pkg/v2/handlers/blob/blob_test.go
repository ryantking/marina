package blob

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

type BlobTestSuite struct {
	suite.Suite
	r http.Handler
}

func (suite *BlobTestSuite) SetupSuite() {
	e := echo.New()
	e.HTTPErrorHandler = web.ErrorHandler
	e.HEAD("/v2/:org/:repo/blobs/:digest", Exists)
	e.GET("/v2/:org/:repo/blobs/:digest", Get)
	e.DELETE("/v2/:org/:repo/blobs/:digest", Delete)
	suite.r = e
}

func (suite *BlobTestSuite) TestParsePath() {
	assert := suite.Assert()
	require := suite.Require()

	c := new(mocks.Context)
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	c.On("Param", "digest").Return("testDigest")
	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	blobExists = func(digest string) (bool, error) {
		return true, nil
	}

	digest, repoName, orgName, err := parsePath(c)
	require.NoError(err)
	assert.Equal("testDigest", digest)
	assert.Equal("testRepo", repoName)
	assert.Equal("testOrg", orgName)
	c.AssertExpectations(suite.T())
}

func (suite *BlobTestSuite) TestParsePathNonExistentRepo() {
	assert := suite.Assert()

	c := new(mocks.Context)
	c.On("Param", "digest").Return("testDigest")
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

func (suite *BlobTestSuite) TestParsePathNonExistentBlob() {
	assert := suite.Assert()

	c := new(mocks.Context)
	c.On("Param", "digest").Return("testDigest")
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	c.On("Set", "docker_err_code", docker.CodeBlobUnknown)
	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	blobExists = func(digest string) (bool, error) {
		return false, nil
	}

	_, _, _, err := parsePath(c)
	assert.Equal(echo.NewHTTPError(http.StatusNotFound, "not found"), err)
	c.AssertExpectations(suite.T())
}

func (suite *BlobTestSuite) TestParsePathExistsRepoError() {
	assert := suite.Assert()

	c := new(mocks.Context)
	c.On("Param", "digest").Return("testDigest")
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	repoExists = func(repoName, orgName string) (bool, error) {
		return false, fmt.Errorf("test error")
	}

	_, _, _, err := parsePath(c)
	assert.EqualError(errors.Cause(err), "test error")
	c.AssertExpectations(suite.T())
}

func (suite *BlobTestSuite) TestParsePathExistsBlobError() {
	assert := suite.Assert()

	c := new(mocks.Context)
	c.On("Param", "digest").Return("testDigest")
	c.On("Param", "repo").Return("testRepo")
	c.On("Param", "org").Return("testOrg")
	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	blobExists = func(digest string) (bool, error) {
		return false, fmt.Errorf("test error")
	}

	_, _, _, err := parsePath(c)
	assert.EqualError(errors.Cause(err), "test error")
	c.AssertExpectations(suite.T())
}

func (suite *BlobTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()

	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	blobExists = func(digest string) (bool, error) {
		return true, nil
	}

	req := httptest.NewRequest(http.MethodHead, "/v2/testOrg/testRepo/blobs/testDigest", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	assert.Equal(fmt.Sprint(len("testDigest")), rr.Header().Get(echo.HeaderContentLength))
	assert.Equal("testDigest", rr.Header().Get(docker.HeaderContentDigest))
}

func (suite *BlobTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	blob := []byte("testBlob")
	r := ioutil.NopCloser(bytes.NewReader(blob))
	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	blobExists = func(digest string) (bool, error) {
		return true, nil
	}
	getBlob = func(digest, repoName, orgName string) (io.ReadCloser, error) {
		return r, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v2/testOrg/testRepo/blobs/testDigest", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	assert.Equal(echo.MIMEOctetStream, rr.Header().Get(echo.HeaderContentType))
}

func (suite *BlobTestSuite) TestDelete() {
	assert := suite.Assert()

	repoExists = func(repoName, orgName string) (bool, error) {
		return true, nil
	}
	blobExists = func(digest string) (bool, error) {
		return true, nil
	}
	deleteBlob = func(digest, repoName, orgName string) error {
		return nil
	}
	deleteBlobEntry = func(digest string) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/v2/testOrg/testRepo/blobs/testDigest", nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	assert.Equal(http.StatusAccepted, rr.Code)
}

func TestBlobTestSuite(t *testing.T) {
	tests := new(BlobTestSuite)
	suite.Run(t, tests)
}
