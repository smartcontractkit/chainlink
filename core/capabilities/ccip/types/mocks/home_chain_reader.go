package mocks

import (
	"context"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/stretchr/testify/mock"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/libocr/ragep2p/types"
)

var _ ccipreaderpkg.HomeChain = (*HomeChainReader)(nil)

type HomeChainReader struct {
	mock.Mock
}

func (_m *HomeChainReader) GetChainConfig(chainSelector cciptypes.ChainSelector) (ccipreaderpkg.ChainConfig, error) {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) GetAllChainConfigs() (map[cciptypes.ChainSelector]ccipreaderpkg.ChainConfig, error) {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) GetSupportedChainsForPeer(id types.PeerID) (mapset.Set[cciptypes.ChainSelector], error) {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) GetKnownCCIPChains() (mapset.Set[cciptypes.ChainSelector], error) {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) GetFChain() (map[cciptypes.ChainSelector]int, error) {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) Start(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) Close() error {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) HealthReport() map[string]error {
	// TODO implement me
	panic("implement me")
}

func (_m *HomeChainReader) Name() string {
	// TODO implement me
	panic("implement me")
}

// GetOCRConfigs provides a mock function with given fields: ctx, donID, pluginType
func (_m *HomeChainReader) GetOCRConfigs(ctx context.Context, donID uint32, pluginType uint8) (ccipreaderpkg.ActiveAndCandidate, error) {
	ret := _m.Called(ctx, donID, pluginType)

	if len(ret) == 0 {
		panic("no return value specified for GetOCRConfigs")
	}

	var r0 ccipreaderpkg.ActiveAndCandidate
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx context.Context, donID uint32, pluginType uint8) (ccipreaderpkg.ActiveAndCandidate, error)); ok {
		return rf(ctx, donID, pluginType)
	}
	if rf, ok := ret.Get(0).(func(ctx context.Context, donID uint32, pluginType uint8) ccipreaderpkg.ActiveAndCandidate); ok {
		r0 = rf(ctx, donID, pluginType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ccipreaderpkg.ActiveAndCandidate)
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
