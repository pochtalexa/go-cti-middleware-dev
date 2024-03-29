// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	storage "github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	mock "github.com/stretchr/testify/mock"
)

// IntStorage is an autogenerated mock type for the IntStorage type
type IntStorage struct {
	mock.Mock
}

// GetAgent provides a mock function with given fields: login
func (_m *IntStorage) GetAgent(login string) (*storage.StAgent, error) {
	ret := _m.Called(login)

	var r0 *storage.StAgent
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*storage.StAgent, error)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) *storage.StAgent); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.StAgent)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveAgent provides a mock function with given fields: login, passHash
func (_m *IntStorage) SaveAgent(login string, passHash []byte) (int64, error) {
	ret := _m.Called(login, passHash)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []byte) (int64, error)); ok {
		return rf(login, passHash)
	}
	if rf, ok := ret.Get(0).(func(string, []byte) int64); ok {
		r0 = rf(login, passHash)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(login, passHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIntStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewIntStorage creates a new instance of IntStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIntStorage(t mockConstructorTestingTNewIntStorage) *IntStorage {
	mock := &IntStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
