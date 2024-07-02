package mocks

import (
	"context"

	mock "github.com/stretchr/testify/mock"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
)

var _ cctypes.HomeChainReader = (*HomeChainReader)(nil)

type HomeChainReader struct {
	mock.Mock
}

// GetOCRConfigs provides a mock function with given fields: ctx, donID, pluginType
func (_m *HomeChainReader) GetOCRConfigs(ctx context.Context, donID uint32, pluginType uint8) ([]ccipreaderpkg.OCR3ConfigWithMeta, error) {
	ret := _m.Called(ctx, donID, pluginType)

	if len(ret) == 0 {
		panic("no return value specified for GetOCRConfigs")
	}

	var r0 []ccipreaderpkg.OCR3ConfigWithMeta
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx context.Context, donID uint32, pluginType uint8) ([]ccipreaderpkg.OCR3ConfigWithMeta, error)); ok {
		return rf(ctx, donID, pluginType)
	}
	if rf, ok := ret.Get(0).(func(ctx context.Context, donID uint32, pluginType uint8) []ccipreaderpkg.OCR3ConfigWithMeta); ok {
		r0 = rf(ctx, donID, pluginType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ccipreaderpkg.OCR3ConfigWithMeta)
		}
	}

	if rf, ok := ret.Get(1).(func(ctx context.Context, donID uint32, pluginType uint8) error); ok {
		r1 = rf(ctx, donID, pluginType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *HomeChainReader) Ready() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewHomeChainReader creates a new instance of HomeChainReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHomeChainReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *HomeChainReader {
	mock := &HomeChainReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
