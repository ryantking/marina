package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ryantking/marina/pkg/store/mocks"
	"github.com/ryantking/marina/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	client *mocks.Client
}

func (suite *StoreTestSuite) SetupTest() {
	suite.client = new(mocks.Client)
	client = suite.client
	testutil.Aquire("Chunk", "Upload")
}

func (suite *StoreTestSuite) TearDownTest() {
	suite.client.AssertExpectations(suite.T())
	testutil.Clean("Chunk", "Upload")
}

func (suite *StoreTestSuite) TestGetBlob() {
	assert := suite.Assert()
	require := suite.Require()

	obj := ioutil.NopCloser(bytes.NewBuffer([]byte("test")))
	digest := "sha256:a464c54f93a9e88fc1d33df1e0e39cca427d60145a360962e8f19a1dbf900da9"
	repo := "alpine"
	org := "library"
	suite.client.On("Get", fmt.Sprintf("blobs/%s/%s/%s.tar.gz", org, repo, digest)).Return(obj, nil)

	r, err := GetBlob(digest, repo, org)
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
	suite.client.On("Put", fmt.Sprintf("uploads/%s/0.tar.gz", uuid), mock.Anything, int64(20)).Return(int64(20), nil)

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
	r1 := ioutil.NopCloser(bytes.NewBuffer([]byte("chunk1")))
	r2 := ioutil.NopCloser(bytes.NewBuffer([]byte("chunk2")))
	suite.client.On("Get", fmt.Sprintf("uploads/%s/0.tar.gz", uuid)).Return(r1, nil)
	suite.client.On("Get", fmt.Sprintf("uploads/%s/1024.tar.gz", uuid)).Return(r2, nil)
	suite.client.On("Put", fmt.Sprintf("blobs/%s/%s/%s.tar.gz", org, repo, digest),
		mock.Anything, int64(2048)).Return(int64(2048), nil)
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
