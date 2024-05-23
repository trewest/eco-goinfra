// Code generated by MockGen. DO NOT EDIT.
// Source: execute.go
//
// Generated by this command:
//
//	mockgen -source=execute.go -package=ops -destination=mock_execute.go
//
// Package ops is a generated GoMock package.
package ops

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockExecute is a mock of Execute interface.
type MockExecute struct {
	ctrl     *gomock.Controller
	recorder *MockExecuteMockRecorder
}

// MockExecuteMockRecorder is the mock recorder for MockExecute.
type MockExecuteMockRecorder struct {
	mock *MockExecute
}

// NewMockExecute creates a new mock instance.
func NewMockExecute(ctrl *gomock.Controller) *MockExecute {
	mock := &MockExecute{ctrl: ctrl}
	mock.recorder = &MockExecuteMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecute) EXPECT() *MockExecuteMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockExecute) Execute(command string, args ...string) (string, error) {
	m.ctrl.T.Helper()
	varargs := []any{command}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockExecuteMockRecorder) Execute(command any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{command}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockExecute)(nil).Execute), varargs...)
}

// ExecuteWithLiveLogger mocks base method.
func (m *MockExecute) ExecuteWithLiveLogger(command string, args ...string) (string, error) {
	m.ctrl.T.Helper()
	varargs := []any{command}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecuteWithLiveLogger", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteWithLiveLogger indicates an expected call of ExecuteWithLiveLogger.
func (mr *MockExecuteMockRecorder) ExecuteWithLiveLogger(command any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{command}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteWithLiveLogger", reflect.TypeOf((*MockExecute)(nil).ExecuteWithLiveLogger), varargs...)
}
