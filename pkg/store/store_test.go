package store

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ryantking/marina/pkg/config"
	"github.com/ryantking/marina/pkg/store/mocks"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

var negOne int64 = -1

type StoreTestSuite struct {
	suite.Suite
	client  *mocks.Client
	cleaner dbcleaner.DbCleaner
}

func (suite *StoreTestSuite) SetupSuite() {
	mysql := engine.NewMySQLEngine(config.Get().DB.DSN)
	suite.cleaner = dbcleaner.New()
	suite.cleaner.SetEngine(mysql)
}

func (suite *StoreTestSuite) SetupTest() {
	suite.client = new(mocks.Client)
	client = suite.client
	suite.cleaner.Acquire("Chunk", "Upload", "Repository", "Organization")
	testutil.Clear(context.Background())
	testutil.Seed(context.Background())
}

func (suite *StoreTestSuite) TearDownTest() {
	suite.client.AssertExpectations(suite.T())
	suite.cleaner.Clean("Chunk", "Upload", "Repository", "Organization")
}

func (suite *StoreTestSuite) TestGetBlob() {
	assert := suite.Assert()
	require := suite.Require()

	obj := ioutil.NopCloser(bytes.NewBuffer([]byte("test")))
	digest := "sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9"
	repo := "alpine"
	org := "library"
	suite.client.On("Get", fmt.Sprintf("blobs/%s/%s/%s.tar.gz", org, repo, digest), negOne, negOne).Return(obj, nil)

	r, err := GetBlob(digest, repo, org, -1, -1)
	require.NoError(err)
	assert.EqualValues(obj, r)
}

func (suite *StoreTestSuite) TestDeleteBlob() {
	assert := suite.Assert()

	digest := "sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9"
	repo := "alpine"
	org := "library"
	suite.client.On("Remove", fmt.Sprintf("blobs/%s/%s/%s.tar.gz", org, repo, digest)).Return(nil)

	err := DeleteBlob(digest, repo, org)
	assert.NoError(err)
}

func (suite *StoreTestSuite) TestUploadChunk() {
	assert := suite.Assert()
	require := suite.Require()

	r := bytes.NewBuffer([]byte("testChunk"))
	uuid := "6b3c9a93-af5d-473f-a4ce-9710022185cd"
	suite.client.On("Put", fmt.Sprintf("uploads/%s/0.tar.gz", uuid), mock.Anything, int32(20)).Return(int32(20), nil)

	n, err := UploadChunk(uuid, r, 20, 0)
	require.NoError(err)
	assert.EqualValues(20, n)
}

func (suite *StoreTestSuite) TestFinishUpload() {
	require := suite.Require()

	uuid := "6b3c9a93-af5d-473f-a4ce-9710022185cd"
	digest := "sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9"
	repo := "alpine"
	org := "library"
	to := fmt.Sprintf("blobs/%s/%s/%s.tar.gz", org, repo, digest)
	froms := []string{fmt.Sprintf("uploads/%s/0.tar.gz", uuid), fmt.Sprintf("uploads/%s/1024.tar.gz", uuid)}
	suite.client.On("Merge", to, froms[0], froms[1]).Return(nil)
	suite.client.On("Remove", fmt.Sprintf("uploads/%s/0.tar.gz", uuid)).Return(nil)
	suite.client.On("Remove", fmt.Sprintf("uploads/%s/1024.tar.gz", uuid)).Return(nil)

	err := FinishUpload(digest, uuid, repo, org)
	require.NoError(err)
}

func (suite *StoreTestSuite) TestDeleteUpload() {
	require := suite.Require()

	uuid := "6b3c9a93-af5d-473f-a4ce-9710022185cd"
	suite.client.On("Remove", fmt.Sprintf("uploads/%s/0.tar.gz", uuid)).Return(nil)
	suite.client.On("Remove", fmt.Sprintf("uploads/%s/1024.tar.gz", uuid)).Return(nil)

	err := DeleteUpload(uuid)
	require.NoError(err)
}

func TestStoreTestSuite(t *testing.T) {
	tests := new(StoreTestSuite)
	suite.Run(t, tests)
}
