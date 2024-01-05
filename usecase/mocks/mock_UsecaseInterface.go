// Code generated by mockery v2.26.1. DO NOT EDIT.

package usecase

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockUsecaseInterface is an autogenerated mock type for the UsecaseInterface type
type MockUsecaseInterface struct {
	mock.Mock
}

type MockUsecaseInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUsecaseInterface) EXPECT() *MockUsecaseInterface_Expecter {
	return &MockUsecaseInterface_Expecter{mock: &_m.Mock}
}

// Batch provides a mock function with given fields:
func (_m *MockUsecaseInterface) Batch() {
	_m.Called()
}

// MockUsecaseInterface_Batch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Batch'
type MockUsecaseInterface_Batch_Call struct {
	*mock.Call
}

// Batch is a helper method to define mock.On call
func (_e *MockUsecaseInterface_Expecter) Batch() *MockUsecaseInterface_Batch_Call {
	return &MockUsecaseInterface_Batch_Call{Call: _e.mock.On("Batch")}
}

func (_c *MockUsecaseInterface_Batch_Call) Run(run func()) *MockUsecaseInterface_Batch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUsecaseInterface_Batch_Call) Return() *MockUsecaseInterface_Batch_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUsecaseInterface_Batch_Call) RunAndReturn(run func()) *MockUsecaseInterface_Batch_Call {
	_c.Call.Return(run)
	return _c
}

// CreateQueue provides a mock function with given fields: ctx, idCard, mobileNo, input
func (_m *MockUsecaseInterface) CreateQueue(ctx context.Context, idCard string, mobileNo string, input map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(ctx, idCard, mobileNo, input)

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]interface{}) (map[string]interface{}, error)); ok {
		return rf(ctx, idCard, mobileNo, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(ctx, idCard, mobileNo, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, idCard, mobileNo, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUsecaseInterface_CreateQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateQueue'
type MockUsecaseInterface_CreateQueue_Call struct {
	*mock.Call
}

// CreateQueue is a helper method to define mock.On call
//   - ctx context.Context
//   - idCard string
//   - mobileNo string
//   - input map[string]interface{}
func (_e *MockUsecaseInterface_Expecter) CreateQueue(ctx interface{}, idCard interface{}, mobileNo interface{}, input interface{}) *MockUsecaseInterface_CreateQueue_Call {
	return &MockUsecaseInterface_CreateQueue_Call{Call: _e.mock.On("CreateQueue", ctx, idCard, mobileNo, input)}
}

func (_c *MockUsecaseInterface_CreateQueue_Call) Run(run func(ctx context.Context, idCard string, mobileNo string, input map[string]interface{})) *MockUsecaseInterface_CreateQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(map[string]interface{}))
	})
	return _c
}

func (_c *MockUsecaseInterface_CreateQueue_Call) Return(_a0 map[string]interface{}, _a1 error) *MockUsecaseInterface_CreateQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUsecaseInterface_CreateQueue_Call) RunAndReturn(run func(context.Context, string, string, map[string]interface{}) (map[string]interface{}, error)) *MockUsecaseInterface_CreateQueue_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteQueue provides a mock function with given fields: ctx, ids
func (_m *MockUsecaseInterface) DeleteQueue(ctx context.Context, ids string) (map[string]interface{}, error) {
	ret := _m.Called(ctx, ids)

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (map[string]interface{}, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]interface{}); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUsecaseInterface_DeleteQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteQueue'
type MockUsecaseInterface_DeleteQueue_Call struct {
	*mock.Call
}

// DeleteQueue is a helper method to define mock.On call
//   - ctx context.Context
//   - ids string
func (_e *MockUsecaseInterface_Expecter) DeleteQueue(ctx interface{}, ids interface{}) *MockUsecaseInterface_DeleteQueue_Call {
	return &MockUsecaseInterface_DeleteQueue_Call{Call: _e.mock.On("DeleteQueue", ctx, ids)}
}

func (_c *MockUsecaseInterface_DeleteQueue_Call) Run(run func(ctx context.Context, ids string)) *MockUsecaseInterface_DeleteQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockUsecaseInterface_DeleteQueue_Call) Return(_a0 map[string]interface{}, _a1 error) *MockUsecaseInterface_DeleteQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUsecaseInterface_DeleteQueue_Call) RunAndReturn(run func(context.Context, string) (map[string]interface{}, error)) *MockUsecaseInterface_DeleteQueue_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueue provides a mock function with given fields: ctx, input
func (_m *MockUsecaseInterface) GetQueue(ctx context.Context, input map[string]interface{}) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, input)

	var r0 []map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) ([]map[string]interface{}, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) []map[string]interface{}); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUsecaseInterface_GetQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueue'
type MockUsecaseInterface_GetQueue_Call struct {
	*mock.Call
}

// GetQueue is a helper method to define mock.On call
//   - ctx context.Context
//   - input map[string]interface{}
func (_e *MockUsecaseInterface_Expecter) GetQueue(ctx interface{}, input interface{}) *MockUsecaseInterface_GetQueue_Call {
	return &MockUsecaseInterface_GetQueue_Call{Call: _e.mock.On("GetQueue", ctx, input)}
}

func (_c *MockUsecaseInterface_GetQueue_Call) Run(run func(ctx context.Context, input map[string]interface{})) *MockUsecaseInterface_GetQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *MockUsecaseInterface_GetQueue_Call) Return(_a0 []map[string]interface{}, _a1 error) *MockUsecaseInterface_GetQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUsecaseInterface_GetQueue_Call) RunAndReturn(run func(context.Context, map[string]interface{}) ([]map[string]interface{}, error)) *MockUsecaseInterface_GetQueue_Call {
	_c.Call.Return(run)
	return _c
}

// SendSmsAllToday provides a mock function with given fields:
func (_m *MockUsecaseInterface) SendSmsAllToday() {
	_m.Called()
}

// MockUsecaseInterface_SendSmsAllToday_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendSmsAllToday'
type MockUsecaseInterface_SendSmsAllToday_Call struct {
	*mock.Call
}

// SendSmsAllToday is a helper method to define mock.On call
func (_e *MockUsecaseInterface_Expecter) SendSmsAllToday() *MockUsecaseInterface_SendSmsAllToday_Call {
	return &MockUsecaseInterface_SendSmsAllToday_Call{Call: _e.mock.On("SendSmsAllToday")}
}

func (_c *MockUsecaseInterface_SendSmsAllToday_Call) Run(run func()) *MockUsecaseInterface_SendSmsAllToday_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUsecaseInterface_SendSmsAllToday_Call) Return() *MockUsecaseInterface_SendSmsAllToday_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUsecaseInterface_SendSmsAllToday_Call) RunAndReturn(run func()) *MockUsecaseInterface_SendSmsAllToday_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateQueue provides a mock function with given fields: ctx, ids, newDate, newSlot, input
func (_m *MockUsecaseInterface) UpdateQueue(ctx context.Context, ids string, newDate string, newSlot int, input map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(ctx, ids, newDate, newSlot, input)

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, map[string]interface{}) (map[string]interface{}, error)); ok {
		return rf(ctx, ids, newDate, newSlot, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(ctx, ids, newDate, newSlot, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, map[string]interface{}) error); ok {
		r1 = rf(ctx, ids, newDate, newSlot, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUsecaseInterface_UpdateQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateQueue'
type MockUsecaseInterface_UpdateQueue_Call struct {
	*mock.Call
}

// UpdateQueue is a helper method to define mock.On call
//   - ctx context.Context
//   - ids string
//   - newDate string
//   - newSlot int
//   - input map[string]interface{}
func (_e *MockUsecaseInterface_Expecter) UpdateQueue(ctx interface{}, ids interface{}, newDate interface{}, newSlot interface{}, input interface{}) *MockUsecaseInterface_UpdateQueue_Call {
	return &MockUsecaseInterface_UpdateQueue_Call{Call: _e.mock.On("UpdateQueue", ctx, ids, newDate, newSlot, input)}
}

func (_c *MockUsecaseInterface_UpdateQueue_Call) Run(run func(ctx context.Context, ids string, newDate string, newSlot int, input map[string]interface{})) *MockUsecaseInterface_UpdateQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(map[string]interface{}))
	})
	return _c
}

func (_c *MockUsecaseInterface_UpdateQueue_Call) Return(_a0 map[string]interface{}, _a1 error) *MockUsecaseInterface_UpdateQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUsecaseInterface_UpdateQueue_Call) RunAndReturn(run func(context.Context, string, string, int, map[string]interface{}) (map[string]interface{}, error)) *MockUsecaseInterface_UpdateQueue_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockUsecaseInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockUsecaseInterface creates a new instance of MockUsecaseInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUsecaseInterface(t mockConstructorTestingTNewMockUsecaseInterface) *MockUsecaseInterface {
	mock := &MockUsecaseInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}