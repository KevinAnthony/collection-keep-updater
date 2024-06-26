// Code generated by mockery v2.35.4. DO NOT EDIT.

package wikipedia

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TableGetterMock is an autogenerated mock type for the TableGetter type
type TableGetterMock struct {
	mock.Mock
}

// GetTablesKeyValue provides a mock function with given fields: ctx, page, lang, cleanRef, keyRows, tables
func (_m *TableGetterMock) GetTablesKeyValue(ctx context.Context, page string, lang string, cleanRef bool, keyRows int, tables ...int) ([][]map[string]string, error) {
	_va := make([]interface{}, len(tables))
	for _i := range tables {
		_va[_i] = tables[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, page, lang, cleanRef, keyRows)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 [][]map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool, int, ...int) ([][]map[string]string, error)); ok {
		return rf(ctx, page, lang, cleanRef, keyRows, tables...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool, int, ...int) [][]map[string]string); ok {
		r0 = rf(ctx, page, lang, cleanRef, keyRows, tables...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, bool, int, ...int) error); ok {
		r1 = rf(ctx, page, lang, cleanRef, keyRows, tables...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTablesMatrix provides a mock function with given fields: ctx, page, lang, cleanRef, tables
func (_m *TableGetterMock) GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
	_va := make([]interface{}, len(tables))
	for _i := range tables {
		_va[_i] = tables[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, page, lang, cleanRef)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 [][][]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool, ...int) ([][][]string, error)); ok {
		return rf(ctx, page, lang, cleanRef, tables...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool, ...int) [][][]string); ok {
		r0 = rf(ctx, page, lang, cleanRef, tables...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][][]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, bool, ...int) error); ok {
		r1 = rf(ctx, page, lang, cleanRef, tables...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUserAgent provides a mock function with given fields: _a0
func (_m *TableGetterMock) SetUserAgent(_a0 string) {
	_m.Called(_a0)
}

// NewTableGetterMock creates a new instance of TableGetterMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTableGetterMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TableGetterMock {
	mock := &TableGetterMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
