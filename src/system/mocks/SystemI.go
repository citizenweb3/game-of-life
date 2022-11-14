// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	system "gameoflife/system"

	mock "github.com/stretchr/testify/mock"

	testing "testing"

	utils "gameoflife/utils"
)

// SystemI is an autogenerated mock type for the SystemI type
type SystemI struct {
	mock.Mock
}

// CreateUserWithRamdomParam provides a mock function with given fields: user
func (_m *SystemI) CreateUserWithRamdomParam(user utils.UserID) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(utils.UserID) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCurrentTime provides a mock function with given fields:
func (_m *SystemI) GetCurrentTime() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetUserList provides a mock function with given fields:
func (_m *SystemI) GetUserList() []utils.UserID {
	ret := _m.Called()

	var r0 []utils.UserID
	if rf, ok := ret.Get(0).(func() []utils.UserID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]utils.UserID)
		}
	}

	return r0
}

// GetUserParam provides a mock function with given fields: userID
func (_m *SystemI) GetUserParam(userID utils.UserID) (system.UsersParam, error) {
	ret := _m.Called(userID)

	var r0 system.UsersParam
	if rf, ok := ret.Get(0).(func(utils.UserID) system.UsersParam); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(system.UsersParam)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(utils.UserID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MoveCurrentTimeForvard provides a mock function with given fields: _a0
func (_m *SystemI) MoveCurrentTimeForvard(_a0 int64) int64 {
	ret := _m.Called(_a0)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// NewSystemI creates a new instance of SystemI. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSystemI(t testing.TB) *SystemI {
	mock := &SystemI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
