// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/store (interfaces: AppKVStore)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAppKVStore is a mock of AppKVStore interface.
type MockAppKVStore struct {
	ctrl     *gomock.Controller
	recorder *MockAppKVStoreMockRecorder
}

// MockAppKVStoreMockRecorder is the mock recorder for MockAppKVStore.
type MockAppKVStoreMockRecorder struct {
	mock *MockAppKVStore
}

// NewMockAppKVStore creates a new mock instance.
func NewMockAppKVStore(ctrl *gomock.Controller) *MockAppKVStore {
	mock := &MockAppKVStore{ctrl: ctrl}
	mock.recorder = &MockAppKVStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppKVStore) EXPECT() *MockAppKVStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAppKVStore) Delete(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAppKVStoreMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAppKVStore)(nil).Delete), arg0, arg1, arg2)
}

// DeleteAll mocks base method.
func (m *MockAppKVStore) DeleteAll(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockAppKVStoreMockRecorder) DeleteAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockAppKVStore)(nil).DeleteAll), arg0)
}

// Get mocks base method.
func (m *MockAppKVStore) Get(arg0, arg1, arg2 string, arg3 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockAppKVStoreMockRecorder) Get(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAppKVStore)(nil).Get), arg0, arg1, arg2, arg3)
}

// Set mocks base method.
func (m *MockAppKVStore) Set(arg0, arg1, arg2 string, arg3 interface{}) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockAppKVStoreMockRecorder) Set(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockAppKVStore)(nil).Set), arg0, arg1, arg2, arg3)
}
