// Code generated by MockGen. DO NOT EDIT.
// Source: collector/interfaces.go

// Package mockcollector is a generated GoMock package.
package mockcollector

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	collector "go.aporeto.io/trireme-lib/collector"
)

// MockEventCollector is a mock of EventCollector interface
// nolint
type MockEventCollector struct {
	ctrl     *gomock.Controller
	recorder *MockEventCollectorMockRecorder
}

// MockEventCollectorMockRecorder is the mock recorder for MockEventCollector
// nolint
type MockEventCollectorMockRecorder struct {
	mock *MockEventCollector
}

// NewMockEventCollector creates a new mock instance
// nolint
func NewMockEventCollector(ctrl *gomock.Controller) *MockEventCollector {
	mock := &MockEventCollector{ctrl: ctrl}
	mock.recorder = &MockEventCollectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockEventCollector) EXPECT() *MockEventCollectorMockRecorder {
	return m.recorder
}

// CollectFlowEvent mocks base method
// nolint
func (m *MockEventCollector) CollectFlowEvent(record *collector.FlowRecord) {
	m.ctrl.Call(m, "CollectFlowEvent", record)
}

// CollectFlowEvent indicates an expected call of CollectFlowEvent
// nolint
func (mr *MockEventCollectorMockRecorder) CollectFlowEvent(record interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectFlowEvent", reflect.TypeOf((*MockEventCollector)(nil).CollectFlowEvent), record)
}

// CollectContainerEvent mocks base method
// nolint
func (m *MockEventCollector) CollectContainerEvent(record *collector.ContainerRecord) {
	m.ctrl.Call(m, "CollectContainerEvent", record)
}

// CollectContainerEvent indicates an expected call of CollectContainerEvent
// nolint
func (mr *MockEventCollectorMockRecorder) CollectContainerEvent(record interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectContainerEvent", reflect.TypeOf((*MockEventCollector)(nil).CollectContainerEvent), record)
}

// CollectUserEvent mocks base method
// nolint
func (m *MockEventCollector) CollectUserEvent(record *collector.UserRecord) {
	m.ctrl.Call(m, "CollectUserEvent", record)
}

// CollectUserEvent indicates an expected call of CollectUserEvent
// nolint
func (mr *MockEventCollectorMockRecorder) CollectUserEvent(record interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectUserEvent", reflect.TypeOf((*MockEventCollector)(nil).CollectUserEvent), record)
}

// CollectTraceEvent mocks base method
// nolint
func (m *MockEventCollector) CollectTraceEvent(records []string) {
	m.ctrl.Call(m, "CollectTraceEvent", records)
}

// CollectTraceEvent indicates an expected call of CollectTraceEvent
// nolint
func (mr *MockEventCollectorMockRecorder) CollectTraceEvent(records interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectTraceEvent", reflect.TypeOf((*MockEventCollector)(nil).CollectTraceEvent), records)
}

// CollectPacketEvent mocks base method
// nolint
func (m *MockEventCollector) CollectPacketEvent(report *collector.PacketReport) {
	m.ctrl.Call(m, "CollectPacketEvent", report)
}

// CollectPacketEvent indicates an expected call of CollectPacketEvent
// nolint
func (mr *MockEventCollectorMockRecorder) CollectPacketEvent(report interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectPacketEvent", reflect.TypeOf((*MockEventCollector)(nil).CollectPacketEvent), report)
}

// CollectCounterEvent mocks base method
// nolint
func (m *MockEventCollector) CollectCounterEvent(counterReport *collector.CounterReport) {
	m.ctrl.Call(m, "CollectCounterEvent", counterReport)
}

// CollectCounterEvent indicates an expected call of CollectCounterEvent
// nolint
func (mr *MockEventCollectorMockRecorder) CollectCounterEvent(counterReport interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectCounterEvent", reflect.TypeOf((*MockEventCollector)(nil).CollectCounterEvent), counterReport)
}
