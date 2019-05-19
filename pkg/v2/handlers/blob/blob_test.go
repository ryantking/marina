package blob

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/ryantking/marina/pkg/web"
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

func (suite *BlobTestSuite) SetupTest() {
	testutil.Aquire("Blob", "Repository", "Organization")
}

func (suite *BlobTestSuite) TearDownTest() {
	testutil.Clean("Blob", "Repository", "Organization")
}

func (suite *BlobTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		org    string
		repo   string
		digest string
		code   int
		body   string
	}{
		{
			"library", "alpine",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusOK, "",
		},
		{
			"alpine", "mysql",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusNotFound,
			`{"errors":[{"code":"NAME_UNKNOWN","message":"repository name not known to registry"}]}`,
		},
		{
			"library", "redis",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusNotFound,
			`{"errors":[{"code":"BLOB_UNKNOWN","message":"blob unknown to registry"}]}`,
		},
	}

	for _, tt := range tests {
		url := fmt.Sprintf("/v2/%s/%s/blobs/%s", tt.org, tt.repo, tt.digest)
		req := httptest.NewRequest(http.MethodHead, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code)
		if tt.code == http.StatusOK {
			assert.Equal(tt.digest, rr.Header().Get(docker.HeaderContentDigest))
			assert.Equal(fmt.Sprint(len(tt.digest)), rr.Header().Get(echo.HeaderContentLength))
		} else {
			b, err := ioutil.ReadAll(rr.Body)
			require.NoError(err)
			assert.JSONEq(tt.body, string(b))
		}
	}
}

func (suite *BlobTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		org    string
		repo   string
		digest string
		code   int
		body   string
	}{
		{
			"library", "alpine",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusOK, "library/alpine test layer",
		},
		{
			"alpine", "mysql",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusNotFound,
			`{"errors":[{"code":"NAME_UNKNOWN","message":"repository name not known to registry"}]}`,
		},
		{
			"library", "redis",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusNotFound,
			`{"errors":[{"code":"BLOB_UNKNOWN","message":"blob unknown to registry"}]}`,
		},
	}

	for _, tt := range tests {
		getBlob = func(digest, repo, org string) (io.ReadCloser, error) {
			assert.Equal(tt.digest, digest)
			assert.Equal(tt.repo, repo)
			assert.Equal(tt.org, org)
			return ioutil.NopCloser(strings.NewReader(tt.body)), nil
		}

		url := fmt.Sprintf("/v2/%s/%s/blobs/%s", tt.org, tt.repo, tt.digest)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code)
		b, err := ioutil.ReadAll(rr.Body)
		require.NoError(err)
		if tt.code == http.StatusOK {
			assert.Equal(tt.digest, rr.Header().Get(docker.HeaderContentDigest))
			assert.Equal(tt.body, string(b))
		} else {
			assert.JSONEq(tt.body, string(b))
		}
	}
}

func (suite *BlobTestSuite) TestDelete() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		org    string
		repo   string
		digest string
		code   int
	}{
		{
			"library", "alpine",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusAccepted,
		},
		{
			"alpine", "mysql",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusNotFound,
		},
		{
			"library", "redis",
			"sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9",
			http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		deleteBlob = func(digest, repo, org string) error {
			assert.Equal(tt.digest, digest)
			assert.Equal(tt.repo, repo)
			assert.Equal(tt.org, org)
			return nil
		}

		url := fmt.Sprintf("/v2/%s/%s/blobs/%s", tt.org, tt.repo, tt.digest)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code)
	}
}

func TestBlobTestSuite(t *testing.T) {
	tests := new(BlobTestSuite)
	suite.Run(t, tests)
}
