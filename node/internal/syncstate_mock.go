// Code generated by MockGen. DO NOT EDIT.
// Source: state/state.go

// Package internal is a generated GoMock package.
package internal

import (
	gomock "github.com/golang/mock/gomock"
	state "github.com/vacp2p/mvds/state"
	reflect "reflect"
)

// MockSyncState is a mock of SyncState interface
type MockSyncState struct {
	ctrl     *gomock.Controller
	recorder *MockSyncStateMockRecorder
}

// MockSyncStateMockRecorder is the mock recorder for MockSyncState
type MockSyncStateMockRecorder struct {
	mock *MockSyncState
}

// NewMockSyncState creates a new mock instance
func NewMockSyncState(ctrl *gomock.Controller) *MockSyncState {
	mock := &MockSyncState{ctrl: ctrl}
	mock.recorder = &MockSyncStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSyncState) EXPECT() *MockSyncStateMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockSyncState) Add(newState state.State) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", newState)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockSyncStateMockRecorder) Add(newState interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSyncState)(nil).Add), newState)
}

// Remove mocks base method
func (m *MockSyncState) Remove(id state.MessageID, peer state.PeerID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", id, peer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockSyncStateMockRecorder) Remove(id, peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockSyncState)(nil).Remove), id, peer)
}

// All mocks base method
func (m *MockSyncState) All(epoch int64) ([]state.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", epoch)
	ret0, _ := ret[0].([]state.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All
func (mr *MockSyncStateMockRecorder) All(epoch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockSyncState)(nil).All), epoch)
}

// Map mocks base method
func (m *MockSyncState) Map(epoch int64, process func(state.State) state.State) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Map", epoch, process)
	ret0, _ := ret[0].(error)
	return ret0
}

// Map indicates an expected call of Map
func (mr *MockSyncStateMockRecorder) Map(epoch, process interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Map", reflect.TypeOf((*MockSyncState)(nil).Map), epoch, process)
}