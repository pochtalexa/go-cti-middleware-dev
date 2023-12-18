// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	sync "sync"

	storage "github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	mock "github.com/stretchr/testify/mock"
)

// IntAgent is an autogenerated mock type for the IntAgent type
type IntAgent struct {
	mock.Mock
}

// DropAgentEvents provides a mock function with given fields: login
func (_m *IntAgent) DropAgentEvents(login string) {
	_m.Called(login)
}

// GetEvents provides a mock function with given fields: login
func (_m *IntAgent) GetEvents(login string) (storage.StAgentEvents, bool) {
	ret := _m.Called(login)

	var r0 storage.StAgentEvents
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (storage.StAgentEvents, bool)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) storage.StAgentEvents); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Get(0).(storage.StAgentEvents)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetMutex provides a mock function with given fields: login
func (_m *IntAgent) GetMutex(login string) (*sync.RWMutex, bool) {
	ret := _m.Called(login)

	var r0 *sync.RWMutex
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (*sync.RWMutex, bool)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) *sync.RWMutex); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sync.RWMutex)
		}
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// IsUpdated provides a mock function with given fields: login
func (_m *IntAgent) IsUpdated(login string) (bool, bool) {
	ret := _m.Called(login)

	var r0 bool
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (bool, bool)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// SetEvent provides a mock function with given fields: event, message
func (_m *IntAgent) SetEvent(event *storage.StWsEvent, message []byte) error {
	ret := _m.Called(event, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(*storage.StWsEvent, []byte) error); ok {
		r0 = rf(event, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetUpdated provides a mock function with given fields: login, val
func (_m *IntAgent) SetUpdated(login string, val bool) {
	_m.Called(login, val)
}

type mockConstructorTestingTNewIntAgent interface {
	mock.TestingT
	Cleanup(func())
}

// NewIntAgent creates a new instance of IntAgent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIntAgent(t mockConstructorTestingTNewIntAgent) *IntAgent {
	mock := &IntAgent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}