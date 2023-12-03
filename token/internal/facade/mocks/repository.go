// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// DeleteToken provides a mock function with given fields: ctx, ssh
func (_m *Repository) DeleteToken(ctx context.Context, ssh string) error {
	ret := _m.Called(ctx, ssh)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, ssh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_DeleteToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteToken'
type Repository_DeleteToken_Call struct {
	*mock.Call
}

// DeleteToken is a helper method to define mock.On call
//   - ctx context.Context
//   - ssh string
func (_e *Repository_Expecter) DeleteToken(ctx interface{}, ssh interface{}) *Repository_DeleteToken_Call {
	return &Repository_DeleteToken_Call{Call: _e.mock.On("DeleteToken", ctx, ssh)}
}

func (_c *Repository_DeleteToken_Call) Run(run func(ctx context.Context, ssh string)) *Repository_DeleteToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_DeleteToken_Call) Return(_a0 error) *Repository_DeleteToken_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_DeleteToken_Call) RunAndReturn(run func(context.Context, string) error) *Repository_DeleteToken_Call {
	_c.Call.Return(run)
	return _c
}

// GetSsh provides a mock function with given fields: ctx, token
func (_m *Repository) GetSsh(ctx context.Context, token string) (string, error) {
	ret := _m.Called(ctx, token)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetSsh_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSsh'
type Repository_GetSsh_Call struct {
	*mock.Call
}

// GetSsh is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *Repository_Expecter) GetSsh(ctx interface{}, token interface{}) *Repository_GetSsh_Call {
	return &Repository_GetSsh_Call{Call: _e.mock.On("GetSsh", ctx, token)}
}

func (_c *Repository_GetSsh_Call) Run(run func(ctx context.Context, token string)) *Repository_GetSsh_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_GetSsh_Call) Return(_a0 string, _a1 error) *Repository_GetSsh_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetSsh_Call) RunAndReturn(run func(context.Context, string) (string, error)) *Repository_GetSsh_Call {
	_c.Call.Return(run)
	return _c
}

// GetToken provides a mock function with given fields: ctx, ssh
func (_m *Repository) GetToken(ctx context.Context, ssh string) (string, error) {
	ret := _m.Called(ctx, ssh)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ssh)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ssh)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ssh)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetToken'
type Repository_GetToken_Call struct {
	*mock.Call
}

// GetToken is a helper method to define mock.On call
//   - ctx context.Context
//   - ssh string
func (_e *Repository_Expecter) GetToken(ctx interface{}, ssh interface{}) *Repository_GetToken_Call {
	return &Repository_GetToken_Call{Call: _e.mock.On("GetToken", ctx, ssh)}
}

func (_c *Repository_GetToken_Call) Run(run func(ctx context.Context, ssh string)) *Repository_GetToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_GetToken_Call) Return(_a0 string, _a1 error) *Repository_GetToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetToken_Call) RunAndReturn(run func(context.Context, string) (string, error)) *Repository_GetToken_Call {
	_c.Call.Return(run)
	return _c
}

// SetToken provides a mock function with given fields: ctx, ssh, token
func (_m *Repository) SetToken(ctx context.Context, ssh string, token string) error {
	ret := _m.Called(ctx, ssh, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, ssh, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_SetToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetToken'
type Repository_SetToken_Call struct {
	*mock.Call
}

// SetToken is a helper method to define mock.On call
//   - ctx context.Context
//   - ssh string
//   - token string
func (_e *Repository_Expecter) SetToken(ctx interface{}, ssh interface{}, token interface{}) *Repository_SetToken_Call {
	return &Repository_SetToken_Call{Call: _e.mock.On("SetToken", ctx, ssh, token)}
}

func (_c *Repository_SetToken_Call) Run(run func(ctx context.Context, ssh string, token string)) *Repository_SetToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Repository_SetToken_Call) Return(_a0 error) *Repository_SetToken_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_SetToken_Call) RunAndReturn(run func(context.Context, string, string) error) *Repository_SetToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
