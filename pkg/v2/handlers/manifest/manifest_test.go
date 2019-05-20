package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/ryantking/marina/pkg/web"
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

func (suite *ManifestTestSuite) SetupTest() {
	testutil.Acquire("Tag", "Image")
}

func (suite *ManifestTestSuite) TearDownTest() {
	testutil.Clean("Tag", "Image")
}

func (suite *ManifestTestSuite) TestExists() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		ref    string
		repo   string
		org    string
		code   int
		digest string
	}{
		{
			"3.9", "alpine", "library",
			http.StatusOK, "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb",
		},
		{
			"sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb", "alpine", "library",
			http.StatusOK, "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb",
		},
		{
			"4.0", "alpine", "library",
			http.StatusNotFound, "",
		},
		{
			"4.0", "mysql", "library",
			http.StatusNotFound, "",
		},
	}

	for _, tt := range tests {
		url := fmt.Sprintf("/v2/%s/%s/manifests/%s", tt.org, tt.repo, tt.ref)
		req := httptest.NewRequest(http.MethodHead, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code)
		if tt.code == http.StatusOK {
			assert.Equal(tt.digest, rr.Header().Get(docker.HeaderContentDigest))
			assert.NotZero(rr.Header().Get(echo.HeaderContentLength))
		}
	}
}

func (suite *ManifestTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		ref      string
		repo     string
		org      string
		code     int
		digest   string
		manifest string
	}{
		{
			"3.9", "alpine", "library",
			http.StatusOK, "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb",
			`{"config": {"digest": "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb"}}`,
		},
		{
			"sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb", "alpine", "library",
			http.StatusOK, "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb",
			`{"config": {"digest": "sha256:cdf98d1859c1beb33ec70507249d34bacf888d59c24df3204057f9a6c758dddb"}}`,
		},
		{
			"4.0", "alpine", "library",
			http.StatusNotFound, "", "",
		},
		{
			"4.0", "mysql", "library",
			http.StatusNotFound, "", "",
		},
	}

	for _, tt := range tests {
		url := fmt.Sprintf("/v2/%s/%s/manifests/%s", tt.org, tt.repo, tt.ref)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		suite.r.ServeHTTP(rr, req)
		require.Equal(tt.code, rr.Code)
		if tt.code == http.StatusOK {
			assert.Equal(tt.digest, rr.Header().Get(docker.HeaderContentDigest))
			assert.Equal(docker.MIMEManifestV2, rr.Header().Get(echo.HeaderContentType))
			b, err := ioutil.ReadAll(rr.Body)
			require.NoError(err)
			assert.JSONEq(tt.manifest, string(b))
		}
	}
}

func (suite *ManifestTestSuite) TestUpdate() {
	assert := suite.Assert()
	require := suite.Require()

	digest := "sha256:055936d3920576da37aa9bc460d70c5f212028bda1c08c0879aedf03d7a66ea1"
	man := new(docker.ManifestV2)
	man.Config.Digest = digest
	b, err := json.Marshal(man)
	require.NoError(err)
	req := httptest.NewRequest(http.MethodPut, "/v2/library/alpine/manifests/latest", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, docker.MIMEManifestV2)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusCreated, rr.Code)
	assert.Equal(digest, rr.Header().Get(docker.HeaderContentDigest))
	assert.Equal("/v2/library/alpine/manifests/latest", rr.Header().Get(echo.HeaderLocation))
}

func (suite *ManifestTestSuite) TestDelete() {
	assert := suite.Assert()

	req := httptest.NewRequest(http.MethodDelete, "/v2/library/alpine/manifests/3.9", nil)
	req.Header.Set(echo.HeaderContentType, docker.MIMEManifestV2)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	assert.Equal(http.StatusAccepted, rr.Code)
}

func TestManifestTestSuite(t *testing.T) {
	tests := new(ManifestTestSuite)
	suite.Run(t, tests)
}
