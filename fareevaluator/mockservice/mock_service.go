// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/andryanduta/taxi-fare/fareevaluator (interfaces: Service)

// Package mockservice is a generated GoMock package.
package mockservice

import (
	fareevaluator "github.com/andryanduta/taxi-fare/fareevaluator"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CalculateFare mocks base method
func (m *MockService) CalculateFare(arg0 []fareevaluator.DistanceMeter) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateFare", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateFare indicates an expected call of CalculateFare
func (mr *MockServiceMockRecorder) CalculateFare(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateFare", reflect.TypeOf((*MockService)(nil).CalculateFare), arg0)
}
