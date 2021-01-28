// Code generated by MockGen. DO NOT EDIT.
// Source: aggregate.go

// Package mock_aggregate is a generated GoMock package.
package mock_aggregate

import (
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	aggregate "github.com/modernice/goes/aggregate"
	event "github.com/modernice/goes/event"
	reflect "reflect"
)

// MockAggregate is a mock of Aggregate interface
type MockAggregate struct {
	ctrl     *gomock.Controller
	recorder *MockAggregateMockRecorder
}

// MockAggregateMockRecorder is the mock recorder for MockAggregate
type MockAggregateMockRecorder struct {
	mock *MockAggregate
}

// NewMockAggregate creates a new mock instance
func NewMockAggregate(ctrl *gomock.Controller) *MockAggregate {
	mock := &MockAggregate{ctrl: ctrl}
	mock.recorder = &MockAggregateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAggregate) EXPECT() *MockAggregateMockRecorder {
	return m.recorder
}

// AggregateID mocks base method
func (m *MockAggregate) AggregateID() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateID")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// AggregateID indicates an expected call of AggregateID
func (mr *MockAggregateMockRecorder) AggregateID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateID", reflect.TypeOf((*MockAggregate)(nil).AggregateID))
}

// AggregateName mocks base method
func (m *MockAggregate) AggregateName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateName")
	ret0, _ := ret[0].(string)
	return ret0
}

// AggregateName indicates an expected call of AggregateName
func (mr *MockAggregateMockRecorder) AggregateName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateName", reflect.TypeOf((*MockAggregate)(nil).AggregateName))
}

// AggregateVersion mocks base method
func (m *MockAggregate) AggregateVersion() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// AggregateVersion indicates an expected call of AggregateVersion
func (mr *MockAggregateMockRecorder) AggregateVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateVersion", reflect.TypeOf((*MockAggregate)(nil).AggregateVersion))
}

// AggregateChanges mocks base method
func (m *MockAggregate) AggregateChanges() []event.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateChanges")
	ret0, _ := ret[0].([]event.Event)
	return ret0
}

// AggregateChanges indicates an expected call of AggregateChanges
func (mr *MockAggregateMockRecorder) AggregateChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateChanges", reflect.TypeOf((*MockAggregate)(nil).AggregateChanges))
}

// TrackChange mocks base method
func (m *MockAggregate) TrackChange(arg0 ...event.Event) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "TrackChange", varargs...)
}

// TrackChange indicates an expected call of TrackChange
func (mr *MockAggregateMockRecorder) TrackChange(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrackChange", reflect.TypeOf((*MockAggregate)(nil).TrackChange), arg0...)
}

// FlushChanges mocks base method
func (m *MockAggregate) FlushChanges() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FlushChanges")
}

// FlushChanges indicates an expected call of FlushChanges
func (mr *MockAggregateMockRecorder) FlushChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushChanges", reflect.TypeOf((*MockAggregate)(nil).FlushChanges))
}

// ApplyEvent mocks base method
func (m *MockAggregate) ApplyEvent(arg0 event.Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ApplyEvent", arg0)
}

// ApplyEvent indicates an expected call of ApplyEvent
func (mr *MockAggregateMockRecorder) ApplyEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyEvent", reflect.TypeOf((*MockAggregate)(nil).ApplyEvent), arg0)
}

// MockFactory is a mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// Make mocks base method
func (m *MockFactory) Make(name string, id uuid.UUID) (aggregate.Aggregate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Make", name, id)
	ret0, _ := ret[0].(aggregate.Aggregate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Make indicates an expected call of Make
func (mr *MockFactoryMockRecorder) Make(name, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Make", reflect.TypeOf((*MockFactory)(nil).Make), name, id)
}
