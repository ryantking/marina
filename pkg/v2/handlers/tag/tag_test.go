package tag

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ryantking/marina/pkg/prisma"
	"github.com/ryantking/marina/pkg/testutil"
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

func (suite *TagTestSuite) SetupTest() {
	testutil.Acquire("Tag")
}

func (suite *TagTestSuite) TearDownTest() {
	testutil.Clean("Tag")
}

func (suite *TagTestSuite) TestList() {
	assert := suite.Assert()
	require := suite.Require()

	repo := "alpine"
	org := "library"
	url := fmt.Sprintf("/v2/%s/%s/tags/list", org, repo)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	expected := fmt.Sprintf(`{"name": "%s/%s", "tags": ["3.9"]}`, org, repo)
	assert.JSONEq(expected, string(b))
}

func (suite *TagTestSuite) TestListPaginated() {
	assert := suite.Assert()
	require := suite.Require()

	repo := "alpine"
	org := "library"
	ctx := context.Background()
	images, err := client.Images(&prisma.ImagesParams{
		Where: &prisma.ImageWhereInput{
			Repo: &prisma.RepositoryWhereInput{
				Name: &repo,
				Org:  &prisma.OrganizationWhereInput{Name: &org},
			},
		},
	}).Exec(ctx)
	require.NoError(err)
	_, err = client.CreateTag(prisma.TagCreateInput{
		Ref: "latest",
		Image: prisma.ImageCreateOneWithoutTagsInput{
			Connect: &prisma.ImageWhereUniqueInput{ID: &images[0].ID},
		},
	}).Exec(ctx)
	require.NoError(err)

	url := fmt.Sprintf("/v2/%s/%s/tags/list?n=1", org, repo)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err := ioutil.ReadAll(rr.Body)
	require.NoError(err)
	expected := fmt.Sprintf(`{"name": "%s/%s", "tags": ["latest"]}`, org, repo)
	assert.JSONEq(expected, string(b))
	link := fmt.Sprintf("/v2/%s/%s/tags/list?n=1&last=latest", org, repo)
	assert.Equal(link, rr.Header().Get(headerLink))

	req = httptest.NewRequest(http.MethodGet, link, nil)
	rr = httptest.NewRecorder()
	suite.r.ServeHTTP(rr, req)
	require.Equal(http.StatusOK, rr.Code)
	b, err = ioutil.ReadAll(rr.Body)
	require.NoError(err)
	expected = fmt.Sprintf(`{"name": "%s/%s", "tags": ["3.9"]}`, org, repo)
	assert.JSONEq(expected, string(b))
	link = fmt.Sprintf("/v2/%s/%s/tags/list?n=1&last=3.9", org, repo)
	assert.Equal(link, rr.Header().Get(headerLink))
}

func TestTagTestSuite(t *testing.T) {
	tests := new(TagTestSuite)
	suite.Run(t, tests)
}
