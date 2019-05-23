package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

type UploadTestSuite struct {
	suite.Suite
	r       http.Handler
	cleaner dbcleaner.DbCleaner
}

func (suite *UploadTestSuite) SetupSuite() {
	e := echo.New()
	e.GET("/v2/:org/:repo/blobs/uploads/:uuid", Get)
	e.POST("/v2/:org/:repo/blobs/uploads", Start)
	e.PATCH("/v2/:org/:repo/blobs/uploads/:uuid", Chunk)
	e.PUT("/v2/:org/:repo/blobs/uploads/:uuid", Finish)
	e.DELETE("/v2/:org/:repo/blobs/uploads/:uuid", Cancel)
	suite.r = e
	mysql := engine.NewMySQLEngine(config.Get().DB.DSN)
	suite.cleaner = dbcleaner.New()
	suite.cleaner.SetEngine(mysql)
}

func (suite *UploadTestSuite) SetupTest() {
	suite.cleaner.Acquire("Chunk", "Upload", "Layer", "Repository", "Organization")
	testutil.Clear(context.Background())
	testutil.Seed(context.Background())
}

func (suite *UploadTestSuite) TearDownTest() {
	suite.cleaner.Clean("Chunk", "Upload", "Layer", "Repository", "Organization")
}

func (suite *UploadTestSuite) TestGet() {
	assert := suite.Assert()
	require := suite.Require()

	repo := "alpine"
	org := "library"
	uuid := "6b3c9a93-af5d-473f-a4ce-9710022185cd"
	url := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", org, repo, uuid)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusNoContent, rr.Code)
	assert.Equal(uuid, rr.Header().Get(docker.HeaderUploadUUID))
	assert.Equal("1024-2047", rr.Header().Get(headerRange))
}

func (suite *UploadTestSuite) TestStart() {
	assert := suite.Assert()
	require := suite.Require()

	repo := "alpine"
	org := "library"
	url := fmt.Sprintf("/v2/%s/%s/blobs/uploads", org, repo)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusAccepted, rr.Code)
	assert.Contains(rr.Header().Get(echo.HeaderLocation), fmt.Sprintf("/v2/%s/%s/blobs/uploads/", org, repo))
	assert.Equal("0-0", rr.Header().Get(headerRange))
}

func (suite *UploadTestSuite) TestChunk() {
	assert := suite.Assert()
	require := suite.Require()

	uuid := "3f497dc6-9458-4c2d-8368-2e71d35c77e5"
	repo := "alpine"
	org := "library"
	r := bytes.NewReader([]byte("chunk1"))
	storeChunk = func(uuid string, r io.Reader, sz, start int32) (int32, error) {
		b, err := ioutil.ReadAll(r)
		require.NoError(err)
		assert.EqualValues("chunk1", b)
		return 6, nil
	}

	url := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", org, repo, uuid)
	req := httptest.NewRequest(http.MethodPatch, url, r)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusNoContent, rr.Code)
	assert.Equal(url, rr.Header().Get(echo.HeaderLocation))
	assert.Equal(uuid, rr.Header().Get(docker.HeaderUploadUUID))
	assert.Equal("0-5", rr.Header().Get(headerRange))
}

func (suite *UploadTestSuite) TestFinish() {
	assert := suite.Assert()
	require := suite.Require()

	digest := "sha256:8a1a56c55249a7e7085ba7482de00d83083d4ebe0c1e782a8ce9d56dd7d3f0a0"
	uuid := "3f497dc6-9458-4c2d-8368-2e71d35c77e5"
	repo := "alpine"
	org := "library"
	r := bytes.NewReader([]byte("chunk1"))
	storeChunk = func(uuid string, r io.Reader, sz, start int32) (int32, error) {
		b, err := ioutil.ReadAll(r)
		require.NoError(err)
		assert.EqualValues("chunk1", b)
		return 6, nil
	}
	finishUpload = func(digest, uuid, repo, org string) error {
		return nil
	}

	url := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s?digest=%s", org, repo, uuid, digest)
	req := httptest.NewRequest(http.MethodPut, url, r)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusCreated, rr.Code)
	loc := fmt.Sprintf("/v2/%s/%s/blobs/%s", org, repo, digest)
	assert.Equal(loc, rr.Header().Get(echo.HeaderLocation))
	assert.Equal(digest, rr.Header().Get(docker.HeaderContentDigest))
}

func (suite *UploadTestSuite) TestCancel() {
	assert := suite.Assert()
	require := suite.Require()

	uuid := "3f497dc6-9458-4c2d-8368-2e71d35c77e5"
	repo := "alpine"
	org := "library"
	deleteUpload = func(uuid string) error {
		assert.Equal("3f497dc6-9458-4c2d-8368-2e71d35c77e5", uuid)
		return nil
	}

	url := fmt.Sprintf("/v2/%s/%s/blobs/uploads/%s", org, repo, uuid)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusNoContent, rr.Code)
}

func TestUploadTestSuite(t *testing.T) {
	tests := new(UploadTestSuite)
	suite.Run(t, tests)
}
