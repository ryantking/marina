// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import db "upper.io/db.v3"
import mock "github.com/stretchr/testify/mock"

// Result is an autogenerated mock type for the Result type
type Result struct {
	mock.Mock
}

// All provides a mock function with given fields: sliceOfStructs
func (_m *Result) All(sliceOfStructs interface{}) error {
	ret := _m.Called(sliceOfStructs)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(sliceOfStructs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// And provides a mock function with given fields: _a0
func (_m *Result) And(_a0 ...interface{}) db.Result {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	ret := _m.Called(_ca...)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(...interface{}) db.Result); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Result) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Count provides a mock function with given fields:
func (_m *Result) Count() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Cursor provides a mock function with given fields: cursorColumn
func (_m *Result) Cursor(cursorColumn string) db.Result {
	ret := _m.Called(cursorColumn)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(string) db.Result); ok {
		r0 = rf(cursorColumn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Delete provides a mock function with given fields:
func (_m *Result) Delete() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Err provides a mock function with given fields:
func (_m *Result) Err() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields:
func (_m *Result) Exists() (bool, error) {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Group provides a mock function with given fields: _a0
func (_m *Result) Group(_a0 ...interface{}) db.Result {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	ret := _m.Called(_ca...)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(...interface{}) db.Result); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Limit provides a mock function with given fields: _a0
func (_m *Result) Limit(_a0 int) db.Result {
	ret := _m.Called(_a0)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(int) db.Result); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Next provides a mock function with given fields: ptrToStruct
func (_m *Result) Next(ptrToStruct interface{}) bool {
	ret := _m.Called(ptrToStruct)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(ptrToStruct)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NextPage provides a mock function with given fields: cursorValue
func (_m *Result) NextPage(cursorValue interface{}) db.Result {
	ret := _m.Called(cursorValue)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(interface{}) db.Result); ok {
		r0 = rf(cursorValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Offset provides a mock function with given fields: _a0
func (_m *Result) Offset(_a0 int) db.Result {
	ret := _m.Called(_a0)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(int) db.Result); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// One provides a mock function with given fields: ptrToStruct
func (_m *Result) One(ptrToStruct interface{}) error {
	ret := _m.Called(ptrToStruct)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(ptrToStruct)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderBy provides a mock function with given fields: _a0
func (_m *Result) OrderBy(_a0 ...interface{}) db.Result {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	ret := _m.Called(_ca...)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(...interface{}) db.Result); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Page provides a mock function with given fields: pageNumber
func (_m *Result) Page(pageNumber uint) db.Result {
	ret := _m.Called(pageNumber)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(uint) db.Result); ok {
		r0 = rf(pageNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Paginate provides a mock function with given fields: pageSize
func (_m *Result) Paginate(pageSize uint) db.Result {
	ret := _m.Called(pageSize)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(uint) db.Result); ok {
		r0 = rf(pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// PrevPage provides a mock function with given fields: cursorValue
func (_m *Result) PrevPage(cursorValue interface{}) db.Result {
	ret := _m.Called(cursorValue)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(interface{}) db.Result); ok {
		r0 = rf(cursorValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// Select provides a mock function with given fields: _a0
func (_m *Result) Select(_a0 ...interface{}) db.Result {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	ret := _m.Called(_ca...)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(...interface{}) db.Result); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *Result) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TotalEntries provides a mock function with given fields:
func (_m *Result) TotalEntries() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TotalPages provides a mock function with given fields:
func (_m *Result) TotalPages() (uint, error) {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Result) Update(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Where provides a mock function with given fields: _a0
func (_m *Result) Where(_a0 ...interface{}) db.Result {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	ret := _m.Called(_ca...)

	var r0 db.Result
	if rf, ok := ret.Get(0).(func(...interface{}) db.Result); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Result)
		}
	}

	return r0
}