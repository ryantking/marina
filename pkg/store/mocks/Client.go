// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Get provides a mock function with given fields: path, start, end
func (_m *Client) Get(path string, start int64, end int64) (io.ReadCloser, error) {
	ret := _m.Called(path, start, end)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string, int64, int64) io.ReadCloser); ok {
		r0 = rf(path, start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64, int64) error); ok {
		r1 = rf(path, start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Merge provides a mock function with given fields: toPath, fromPaths
func (_m *Client) Merge(toPath string, fromPaths ...string) error {
	_va := make([]interface{}, len(fromPaths))
	for _i := range fromPaths {
		_va[_i] = fromPaths[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, toPath)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...string) error); ok {
		r0 = rf(toPath, fromPaths...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Put provides a mock function with given fields: path, r, sz
func (_m *Client) Put(path string, r io.Reader, sz int32) (int32, error) {
	ret := _m.Called(path, r, sz)

	var r0 int32
	if rf, ok := ret.Get(0).(func(string, io.Reader, int32) int32); ok {
		r0 = rf(path, r, sz)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, io.Reader, int32) error); ok {
		r1 = rf(path, r, sz)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: path
func (_m *Client) Remove(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
