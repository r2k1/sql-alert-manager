// Code generated by mockery v1.0.0. DO NOT EDIT.

package alert

import mock "github.com/stretchr/testify/mock"

// MockDestination is an autogenerated mock type for the Destination type
type MockDestination struct {
	mock.Mock
}

// Name provides a mock function with given fields:
func (_m *MockDestination) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ResolveAlert provides a mock function with given fields: a, msg
func (_m *MockDestination) ResolveAlert(a *Alert, msg string) error {
	ret := _m.Called(a, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Alert, string) error); ok {
		r0 = rf(a, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendAlert provides a mock function with given fields: a, msg
func (_m *MockDestination) SendAlert(a *Alert, msg string) error {
	ret := _m.Called(a, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Alert, string) error); ok {
		r0 = rf(a, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}