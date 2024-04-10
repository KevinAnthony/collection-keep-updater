// Code generated by mockery v2.35.4. DO NOT EDIT.

package di

import (
	types "github.com/kevinanthony/collection-keep-updater/types"
	mock "github.com/stretchr/testify/mock"
)

// IDepFactoryMock is an autogenerated mock type for the IDepFactory type
type IDepFactoryMock struct {
	mock.Mock
}

// Config provides a mock function with given fields: cmd, icfg
func (_m *IDepFactoryMock) Config(cmd types.ICommand, icfg types.IConfig) error {
	ret := _m.Called(cmd, icfg)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.ICommand, types.IConfig) error); ok {
		r0 = rf(cmd, icfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Sources provides a mock function with given fields: cmd
func (_m *IDepFactoryMock) Sources(cmd types.ICommand) error {
	ret := _m.Called(cmd)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.ICommand) error); ok {
		r0 = rf(cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIDepFactoryMock creates a new instance of IDepFactoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDepFactoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDepFactoryMock {
	mock := &IDepFactoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
