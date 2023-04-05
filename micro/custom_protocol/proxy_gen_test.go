// Code generated by MockGen. DO NOT EDIT.
// Source: proxy.go

// Package proxy_v2 is a generated GoMock package.
package proxy_v2

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProxy is a mock of Proxy interface.
type MockProxy struct {
	ctrl     *gomock.Controller
	recorder *MockProxyMockRecorder
}

// MockProxyMockRecorder is the mock recorder for MockProxy.
type MockProxyMockRecorder struct {
	mock *MockProxy
}

// NewMockProxy creates a new mock instance.
func NewMockProxy(ctrl *gomock.Controller) *MockProxy {
	mock := &MockProxy{ctrl: ctrl}
	mock.recorder = &MockProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProxy) EXPECT() *MockProxyMockRecorder {
	return m.recorder
}

// Invoke mocks base method.
func (m *MockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invoke", ctx, req)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Invoke indicates an expected call of Invoke.
func (mr *MockProxyMockRecorder) Invoke(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockProxy)(nil).Invoke), ctx, req)
}
