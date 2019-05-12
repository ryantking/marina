package store

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/ryantking/marina/pkg/db/models/upload/chunk"
	"github.com/ryantking/marina/pkg/store/mocks"
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
}

func (suite *StoreTestSuite) TearDownTest() {
	suite.client.AssertExpectations(suite.T())
}

func (suite *StoreTestSuite) TestGetBlob() {
	assert := suite.Assert()
	require := suite.Require()

	obj := ioutil.NopCloser(bytes.NewBuffer([]byte("test")))
	suite.client.On("Get", "blobs/testOrg/testRepo/testDigest.tar.gz").Return(obj, nil)

	r, err := GetBlob("testDigest", "testRepo", "testOrg")
	require.NoError(err)
	assert.EqualValues(obj, r)
}

func (suite *StoreTestSuite) TestDeleteBlob() {
	assert := suite.Assert()

	suite.client.On("Remove", "blobs/testOrg/testRepo/testDigest.tar.gz").Return(nil)

	err := DeleteBlob("testDigest", "testRepo", "testOrg")
	assert.NoError(err)
}

func (suite *StoreTestSuite) TestUploadChunk() {
	assert := suite.Assert()
	require := suite.Require()

	r := bytes.NewBuffer([]byte("testChunk"))
	suite.client.On("Put", "uploads/testUUID/0.tar.gz", mock.Anything, int64(20)).Return(int64(20), nil)

	n, err := UploadChunk("testUUID", r, 20, 0)
	require.NoError(err)
	assert.EqualValues(20, n)
}

func (suite *StoreTestSuite) TestFinishUpload() {
	assert := suite.Assert()
	require := suite.Require()

	r1 := ioutil.NopCloser(bytes.NewBuffer([]byte("chunk1")))
	r2 := ioutil.NopCloser(bytes.NewBuffer([]byte("chunk2")))
	getChunks = func(uuid string) ([]*chunk.Model, error) {
		assert.Equal("testUUID", uuid)
		chunks := []*chunk.Model{
			&chunk.Model{UUID: "testUUID", RangeStart: 0, RangeEnd: 9},
			&chunk.Model{UUID: "testUUID", RangeStart: 10, RangeEnd: 19},
		}
		return chunks, nil
	}
	suite.client.On("Get", "uploads/testUUID/0.tar.gz").Return(r1, nil)
	suite.client.On("Get", "uploads/testUUID/10.tar.gz").Return(r2, nil)
	suite.client.On("Put", "blobs/testOrg/testRepo/testDigest.tar.gz", mock.Anything, int64(20)).Return(int64(20), nil)
	suite.client.On("Remove", "uploads/testUUID/0.tar.gz").Return(nil)
	suite.client.On("Remove", "uploads/testUUID/10.tar.gz").Return(nil)

	err := FinishUpload("testDigest", "testUUID", "testRepo", "testOrg")
	require.NoError(err)
}

func (suite *StoreTestSuite) TestDeleteUpload() {
	assert := suite.Assert()
	require := suite.Require()

	getChunks = func(uuid string) ([]*chunk.Model, error) {
		assert.Equal("testUUID", uuid)
		chunks := []*chunk.Model{
			&chunk.Model{UUID: "testUUID", RangeStart: 0, RangeEnd: 9},
			&chunk.Model{UUID: "testUUID", RangeStart: 10, RangeEnd: 19},
		}
		return chunks, nil
	}
	suite.client.On("Remove", "uploads/testUUID/0.tar.gz").Return(nil)
	suite.client.On("Remove", "uploads/testUUID/10.tar.gz").Return(nil)

	err := DeleteUpload("testUUID")
	require.NoError(err)
}

func TestStoreTestSuite(t *testing.T) {
	tests := new(StoreTestSuite)
	suite.Run(t, tests)
}
