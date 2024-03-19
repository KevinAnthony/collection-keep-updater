// Code generated by mockery v2.35.4. DO NOT EDIT.

package types

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	pflag "github.com/spf13/pflag"
)

// ICommandMock is an autogenerated mock type for the ICommand type
type ICommandMock struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *ICommandMock) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// Flag provides a mock function with given fields: name
func (_m *ICommandMock) Flag(name string) *pflag.Flag {
	ret := _m.Called(name)

	var r0 *pflag.Flag
	if rf, ok := ret.Get(0).(func(string) *pflag.Flag); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pflag.Flag)
		}
	}

	return r0
}

// Flags provides a mock function with given fields:
func (_m *ICommandMock) Flags() *pflag.FlagSet {
	ret := _m.Called()

	var r0 *pflag.FlagSet
	if rf, ok := ret.Get(0).(func() *pflag.FlagSet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pflag.FlagSet)
		}
	}

	return r0
}

// OutOrStdout provides a mock function with given fields:
func (_m *ICommandMock) OutOrStdout() io.Writer {
	ret := _m.Called()

	var r0 io.Writer
	if rf, ok := ret.Get(0).(func() io.Writer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Writer)
		}
	}

	return r0
}

// PersistentFlags provides a mock function with given fields:
func (_m *ICommandMock) PersistentFlags() *pflag.FlagSet {
	ret := _m.Called()

	var r0 *pflag.FlagSet
	if rf, ok := ret.Get(0).(func() *pflag.FlagSet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pflag.FlagSet)
		}
	}

	return r0
}

// SetContext provides a mock function with given fields: ctx
func (_m *ICommandMock) SetContext(ctx context.Context) {
	_m.Called(ctx)
}

// NewICommandMock creates a new instance of ICommandMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICommandMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICommandMock {
	mock := &ICommandMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}