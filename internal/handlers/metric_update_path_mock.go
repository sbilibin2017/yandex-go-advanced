// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/yandex-go-advanced/internal/handlers/metric_update_path.go

// Package handlers is a generated GoMock package.
package handlers

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MockMetricUpdaterPath is a mock of MetricUpdaterPath interface.
type MockMetricUpdaterPath struct {
	ctrl     *gomock.Controller
	recorder *MockMetricUpdaterPathMockRecorder
}

// MockMetricUpdaterPathMockRecorder is the mock recorder for MockMetricUpdaterPath.
type MockMetricUpdaterPathMockRecorder struct {
	mock *MockMetricUpdaterPath
}

// NewMockMetricUpdaterPath creates a new mock instance.
func NewMockMetricUpdaterPath(ctrl *gomock.Controller) *MockMetricUpdaterPath {
	mock := &MockMetricUpdaterPath{ctrl: ctrl}
	mock.recorder = &MockMetricUpdaterPathMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricUpdaterPath) EXPECT() *MockMetricUpdaterPathMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockMetricUpdaterPath) Update(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, metrics)
	ret0, _ := ret[0].([]*types.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockMetricUpdaterPathMockRecorder) Update(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMetricUpdaterPath)(nil).Update), ctx, metrics)
}
