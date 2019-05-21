package upload

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/docker"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type UploadTestSuite struct {
	suite.Suite
	r http.Handler
}

func (suite *UploadTestSuite) SetupSuite() {
	e := echo.New()
	e.GET("/v2/:org/:repo/blobs/uploads/:uuid", Get)
	e.POST("/v2/:org/:repo/blobs/uploads", Start)
	suite.r = e
}

func (suite *UploadTestSuite) SetupTest() {
	testutil.Acquire("Chunk", "Upload")
}

func (suite *UploadTestSuite) TearDownTest() {
	testutil.Clean("Chunk", "Upload")
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

func TestUploadTestSuite(t *testing.T) {
	tests := new(UploadTestSuite)
	suite.Run(t, tests)
}
