// Code generated by mockery v2.35.4. DO NOT EDIT.

package types

import mock "github.com/stretchr/testify/mock"

// ISourceConfigMock is an autogenerated mock type for the ISourceConfig type
type ISourceConfigMock struct {
	mock.Mock
}

// GetIDFromURL provides a mock function with given fields: url
func (_m *ISourceConfigMock) GetIDFromURL(url string) (string, error) {
	ret := _m.Called(url)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SourceSettingFromConfig provides a mock function with given fields: data
func (_m *ISourceConfigMock) SourceSettingFromConfig(data map[string]interface{}) ISourceSettings {
	ret := _m.Called(data)

	var r0 ISourceSettings
	if rf, ok := ret.Get(0).(func(map[string]interface{}) ISourceSettings); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ISourceSettings)
		}
	}

	return r0
}

// SourceSettingFromFlags provides a mock function with given fields: cmd, original
func (_m *ISourceConfigMock) SourceSettingFromFlags(cmd ICommand, original ISourceSettings) (ISourceSettings, error) {
	ret := _m.Called(cmd, original)

	var r0 ISourceSettings
	var r1 error
	if rf, ok := ret.Get(0).(func(ICommand, ISourceSettings) (ISourceSettings, error)); ok {
		return rf(cmd, original)
	}
	if rf, ok := ret.Get(0).(func(ICommand, ISourceSettings) ISourceSettings); ok {
		r0 = rf(cmd, original)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ISourceSettings)
		}
	}

	if rf, ok := ret.Get(1).(func(ICommand, ISourceSettings) error); ok {
		r1 = rf(cmd, original)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewISourceConfigMock creates a new instance of ISourceConfigMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISourceConfigMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISourceConfigMock {
	mock := &ISourceConfigMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
