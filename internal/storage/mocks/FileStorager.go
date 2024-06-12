// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	userdata "github.com/impr0ver/gophKeeper/internal/userdata"
)

// FileStorager is an autogenerated mock type for the FileStorager type
type FileStorager struct {
	mock.Mock
}

// CreateRecord provides a mock function with given fields: ctx, record
func (_m *FileStorager) CreateRecord(ctx context.Context, record userdata.Record) (string, error) {
	ret := _m.Called(ctx, record)

	if len(ret) == 0 {
		panic("no return value specified for CreateRecord")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, userdata.Record) (string, error)); ok {
		return rf(ctx, record)
	}
	if rf, ok := ret.Get(0).(func(context.Context, userdata.Record) string); ok {
		r0 = rf(ctx, record)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, userdata.Record) error); ok {
		r1 = rf(ctx, record)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRecord provides a mock function with given fields: ctx, recordID
func (_m *FileStorager) DeleteRecord(ctx context.Context, recordID string) error {
	ret := _m.Called(ctx, recordID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRecord")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, recordID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRecord provides a mock function with given fields: ctx, recordID
func (_m *FileStorager) GetRecord(ctx context.Context, recordID string) (userdata.Record, error) {
	ret := _m.Called(ctx, recordID)

	if len(ret) == 0 {
		panic("no return value specified for GetRecord")
	}

	var r0 userdata.Record
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (userdata.Record, error)); ok {
		return rf(ctx, recordID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) userdata.Record); ok {
		r0 = rf(ctx, recordID)
	} else {
		r0 = ret.Get(0).(userdata.Record)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, recordID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFileStorager creates a new instance of FileStorager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileStorager(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileStorager {
	mock := &FileStorager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
