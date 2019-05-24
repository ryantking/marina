package blob

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

type BlobTestSuite struct {
	suite.Suite
	r       http.Handler
	cleaner dbcleaner.DbCleaner
}

func (suite *BlobTestSuite) SetupSuite() {
	e := echo.New()
	e.HTTPErrorHandler = web.ErrorHandler
	e.HEAD("/v2/:org/:repo/blobs/:digest", Exists)
	e.GET("/v2/:org/:repo/blobs/:digest", Get)
	e.DELETE("/v2/:org/:repo/blobs/:digest", Delete)
	suite.r = e
	mysql := engine.NewMySQLEngine(config.Get().DB.DSN)
	suite.cleaner = dbcleaner.New()
	suite.cleaner.SetEngine(mysql)
}

func (suite *BlobTestSuite) SetupTest() {
	suite.cleaner.Acquire("Blob", "Repository", "Organization")
	testutil.Clear(context.Background())
	testutil.Seed(context.Background())
}

func (suite *BlobTestSuite) TearDownTest() {
	suite.cleaner.Clean("Blob", "Repository", "Organization")
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
		require.Equal(tt.code, rr.Code, "For: %s/%s", tt.org, tt.repo)
		if tt.code == http.StatusOK {
			assert.Equal(tt.digest, rr.Header().Get(docker.HeaderContentDigest), "For: %s/%s", tt.org, tt.repo)
			l := fmt.Sprint(len(tt.digest))
			assert.Equal(l, rr.Header().Get(echo.HeaderContentLength), "For: %s/%s", tt.org, tt.repo)
		} else {
			b, err := ioutil.ReadAll(rr.Body)
			require.NoError(err, "For: %s/%s", tt.org, tt.repo)
			assert.JSONEq(tt.body, string(b), "For: %s/%s", tt.org, tt.repo)
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
			assert.Equal(tt.digest, digest, "For: %s/%s", tt.org, tt.repo)
			assert.Equal(tt.repo, repo, "For: %s/%s", tt.org, tt.repo)
			assert.Equal(tt.org, org, "For: %s/%s", tt.org, tt.repo)
			return ioutil.NopCloser(strings.NewReader(tt.body)), nil
		}

		url := fmt.Sprintf("/v2/%s/%s/blobs/%s", tt.org, tt.repo, tt.digest)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code, "For: %s/%s", tt.org, tt.repo)
		b, err := ioutil.ReadAll(rr.Body)
		require.NoError(err, "For: %s/%s", tt.org, tt.repo)
		if tt.code == http.StatusOK {
			assert.Equal(tt.digest, rr.Header().Get(docker.HeaderContentDigest), "For: %s/%s", tt.org, tt.repo)
			assert.Equal(tt.body, string(b), "For: %s/%s", tt.org, tt.repo)
		} else {
			assert.JSONEq(tt.body, string(b), "For: %s/%s", tt.org, tt.repo)
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
			assert.Equal(tt.digest, digest, "For: %s/%s", tt.org, tt.repo)
			assert.Equal(tt.repo, repo, "For: %s/%s", tt.org, tt.repo)
			assert.Equal(tt.org, org, "For: %s/%s", tt.org, tt.repo)
			return nil
		}

		url := fmt.Sprintf("/v2/%s/%s/blobs/%s", tt.org, tt.repo, tt.digest)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code, "For: %s/%s", tt.org, tt.repo)
	}
}

func TestBlobTestSuite(t *testing.T) {
	tests := new(BlobTestSuite)
	suite.Run(t, tests)
}
