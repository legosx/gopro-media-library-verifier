// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/legosx/gopro-media-library-verifier/dirscan (interfaces: OS)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/os.go -package=mocks github.com/legosx/gopro-media-library-verifier/dirscan OS
//

// Package mocks is a generated GoMock package.
package mocks

import (
	fs "io/fs"
	reflect "reflect"

	dirscan "github.com/legosx/gopro-media-library-verifier/dirscan"
	gomock "go.uber.org/mock/gomock"
)

// MockOS is a mock of OS interface.
type MockOS struct {
	ctrl     *gomock.Controller
	recorder *MockOSMockRecorder
}

// MockOSMockRecorder is the mock recorder for MockOS.
type MockOSMockRecorder struct {
	mock *MockOS
}

// NewMockOS creates a new mock instance.
func NewMockOS(ctrl *gomock.Controller) *MockOS {
	mock := &MockOS{ctrl: ctrl}
	mock.recorder = &MockOSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOS) EXPECT() *MockOSMockRecorder {
	return m.recorder
}

// IsNotExist mocks base method.
func (m *MockOS) IsNotExist(arg0 error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNotExist", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNotExist indicates an expected call of IsNotExist.
func (mr *MockOSMockRecorder) IsNotExist(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNotExist", reflect.TypeOf((*MockOS)(nil).IsNotExist), arg0)
}

// Open mocks base method.
func (m *MockOS) Open(arg0 string) (dirscan.OSFile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", arg0)
	ret0, _ := ret[0].(dirscan.OSFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockOSMockRecorder) Open(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockOS)(nil).Open), arg0)
}

// Stat mocks base method.
func (m *MockOS) Stat(arg0 string) (fs.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", arg0)
	ret0, _ := ret[0].(fs.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *MockOSMockRecorder) Stat(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockOS)(nil).Stat), arg0)
}
