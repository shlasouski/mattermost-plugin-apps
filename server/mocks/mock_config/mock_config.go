// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/config (interfaces: Service)

// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/mattermost/mattermost-plugin-apps/server/config"
	model "github.com/mattermost/mattermost-server/v5/model"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetConfig mocks base method.
func (m *MockService) GetConfig() config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(config.Config)
	return ret0
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockServiceMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockService)(nil).GetConfig))
}

// GetMattermostConfig mocks base method.
func (m *MockService) GetMattermostConfig() *model.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMattermostConfig")
	ret0, _ := ret[0].(*model.Config)
	return ret0
}

// GetMattermostConfig indicates an expected call of GetMattermostConfig.
func (mr *MockServiceMockRecorder) GetMattermostConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMattermostConfig", reflect.TypeOf((*MockService)(nil).GetMattermostConfig))
}

// Reconfigure mocks base method.
func (m *MockService) Reconfigure(arg0 config.StoredConfig, arg1 ...config.Configurable) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Reconfigure", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reconfigure indicates an expected call of Reconfigure.
func (mr *MockServiceMockRecorder) Reconfigure(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reconfigure", reflect.TypeOf((*MockService)(nil).Reconfigure), varargs...)
}

// StoreConfig mocks base method.
func (m *MockService) StoreConfig(arg0 config.StoredConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreConfig indicates an expected call of StoreConfig.
func (mr *MockServiceMockRecorder) StoreConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreConfig", reflect.TypeOf((*MockService)(nil).StoreConfig), arg0)
}
