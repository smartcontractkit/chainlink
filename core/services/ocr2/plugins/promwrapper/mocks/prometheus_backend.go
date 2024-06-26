// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PrometheusBackend is an autogenerated mock type for the PrometheusBackend type
type PrometheusBackend struct {
	mock.Mock
}

// SetAcceptFinalizedReportToTransmitAcceptedReportLatency provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetAcceptFinalizedReportToTransmitAcceptedReportLatency(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetCloseDuration provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetCloseDuration(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetObservationDuration provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetObservationDuration(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetObservationToReportLatency provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetObservationToReportLatency(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetQueryDuration provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetQueryDuration(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetQueryToObservationLatency provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetQueryToObservationLatency(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetReportDuration provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetReportDuration(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetReportToAcceptFinalizedReportLatency provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetReportToAcceptFinalizedReportLatency(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetShouldAcceptFinalizedReportDuration provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetShouldAcceptFinalizedReportDuration(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// SetShouldTransmitAcceptedReportDuration provides a mock function with given fields: _a0, _a1
func (_m *PrometheusBackend) SetShouldTransmitAcceptedReportDuration(_a0 []string, _a1 float64) {
	_m.Called(_a0, _a1)
}

// NewPrometheusBackend creates a new instance of PrometheusBackend. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPrometheusBackend(t interface {
	mock.TestingT
	Cleanup(func())
}) *PrometheusBackend {
	mock := &PrometheusBackend{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
