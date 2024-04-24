package testutils

import (
	"testing"

	evmclmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
)

func NewEthClientMock(t *testing.T) *evmclmocks.Client {
	return evmclmocks.NewClient(t)
}

func NewEthClientMockWithDefaultChain(t *testing.T) *evmclmocks.Client {
	c := NewEthClientMock(t)
	c.On("ConfiguredChainID").Return(FixtureChainID).Maybe()
	//c.On("IsL2").Return(false).Maybe()
	return c
}
