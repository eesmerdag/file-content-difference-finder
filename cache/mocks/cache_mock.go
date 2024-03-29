// Code generated by MockGen. DO NOT EDIT.
// Source: ./cache/cache.go
//
// Generated by this command:
//
//	mockgen -source=./cache/cache.go -destination=./cache/mocks/cache_mock.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockICache is a mock of ICache interface.
type MockICache struct {
	ctrl     *gomock.Controller
	recorder *MockICacheMockRecorder
}

// MockICacheMockRecorder is the mock recorder for MockICache.
type MockICacheMockRecorder struct {
	mock *MockICache
}

// NewMockICache creates a new mock instance.
func NewMockICache(ctrl *gomock.Controller) *MockICache {
	mock := &MockICache{ctrl: ctrl}
	mock.recorder = &MockICacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICache) EXPECT() *MockICacheMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockICache) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockICacheMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockICache)(nil).Clear))
}

// Delete mocks base method.
func (m *MockICache) Delete(key string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", key)
}

// Delete indicates an expected call of Delete.
func (mr *MockICacheMockRecorder) Delete(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockICache)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockICache) Get(key string) (any, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockICacheMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockICache)(nil).Get), key)
}

// Set mocks base method.
func (m *MockICache) Set(key string, value any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value)
}

// Set indicates an expected call of Set.
func (mr *MockICacheMockRecorder) Set(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockICache)(nil).Set), key, value)
}
