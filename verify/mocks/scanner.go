// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/legosx/gopro-media-library-verifier/verify (interfaces: Scanner)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/scanner.go -package=mocks github.com/legosx/gopro-media-library-verifier/verify Scanner
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dirscan "github.com/legosx/gopro-media-library-verifier/dirscan"
	gomock "go.uber.org/mock/gomock"
)

// MockScanner is a mock of Scanner interface.
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner.
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance.
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// GetFileList mocks base method.
func (m *MockScanner) GetFileList(arg0 string) ([]dirscan.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileList", arg0)
	ret0, _ := ret[0].([]dirscan.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileList indicates an expected call of GetFileList.
func (mr *MockScannerMockRecorder) GetFileList(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileList", reflect.TypeOf((*MockScanner)(nil).GetFileList), arg0)
}
