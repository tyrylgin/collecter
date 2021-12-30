// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package storagemock is a generated GoMock package.
package storagemock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/tyrylgin/collecter/model"
)

// MockMetricStorer is a mock of MetricStorer interface.
type MockMetricStorer struct {
	ctrl     *gomock.Controller
	recorder *MockMetricStorerMockRecorder
}

// MockMetricStorerMockRecorder is the mock recorder for MockMetricStorer.
type MockMetricStorerMockRecorder struct {
	mock *MockMetricStorer
}

// NewMockMetricStorer creates a new mock instance.
func NewMockMetricStorer(ctrl *gomock.Controller) *MockMetricStorer {
	mock := &MockMetricStorer{ctrl: ctrl}
	mock.recorder = &MockMetricStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricStorer) EXPECT() *MockMetricStorerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMetricStorer) Get(name string) model.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", name)
	ret0, _ := ret[0].(model.Metric)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockMetricStorerMockRecorder) Get(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMetricStorer)(nil).Get), name)
}

// GetAll mocks base method.
func (m *MockMetricStorer) GetAll() map[string]model.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(map[string]model.Metric)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricStorerMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricStorer)(nil).GetAll))
}

// Save mocks base method.
func (m *MockMetricStorer) Save(name string, metric model.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", name, metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockMetricStorerMockRecorder) Save(name, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockMetricStorer)(nil).Save), name, metric)
}
