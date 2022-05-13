// Code generated by MockGen. DO NOT EDIT.
// Source: alice.go

// Package mock_alice is a generated GoMock package.
package mock_alice

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	alice "github.com/sku4/alice-checklist/models/alice"
)

// MockChanAnswer is a mock of ChanAnswer interface.
type MockChanAnswer struct {
	ctrl     *gomock.Controller
	recorder *MockChanAnswerMockRecorder
}

// MockChanAnswerMockRecorder is the mock recorder for MockChanAnswer.
type MockChanAnswerMockRecorder struct {
	mock *MockChanAnswer
}

// NewMockChanAnswer creates a new mock instance.
func NewMockChanAnswer(ctrl *gomock.Controller) *MockChanAnswer {
	mock := &MockChanAnswer{ctrl: ctrl}
	mock.recorder = &MockChanAnswerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChanAnswer) EXPECT() *MockChanAnswerMockRecorder {
	return m.recorder
}

// ColdAnswer mocks base method.
func (m *MockChanAnswer) ColdAnswer(arg0 int, arg1 alice.Response, arg2 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ColdAnswer", arg0, arg1, arg2)
}

// ColdAnswer indicates an expected call of ColdAnswer.
func (mr *MockChanAnswerMockRecorder) ColdAnswer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ColdAnswer", reflect.TypeOf((*MockChanAnswer)(nil).ColdAnswer), arg0, arg1, arg2)
}

// DropAnswer mocks base method.
func (m *MockChanAnswer) DropAnswer(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DropAnswer", arg0)
}

// DropAnswer indicates an expected call of DropAnswer.
func (mr *MockChanAnswerMockRecorder) DropAnswer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropAnswer", reflect.TypeOf((*MockChanAnswer)(nil).DropAnswer), arg0)
}

// HotAnswer mocks base method.
func (m *MockChanAnswer) HotAnswer(arg0 int, arg1 alice.Response) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HotAnswer", arg0, arg1)
}

// HotAnswer indicates an expected call of HotAnswer.
func (mr *MockChanAnswerMockRecorder) HotAnswer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HotAnswer", reflect.TypeOf((*MockChanAnswer)(nil).HotAnswer), arg0, arg1)
}
