package eal

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlockchainClientEstimateGas(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(100000), nil)

	blockchainClient := &BlockchainClient{
		client: mockClient,
	}

	gasLimit, err := blockchainClient.EstimateGas(context.Background(), common.Address{}, []byte{})
	assert.NoError(t, err)
	assert.Equal(t, uint32(100000), gasLimit)

	// Verify that the EstimateGas method on the mock client was called with the expected arguments
	mockClient.AssertCalled(t, "EstimateGas", mock.Anything, mock.Anything)
}

func TestBlockchainClientSimulateTransaction(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	blockchainClient := &BlockchainClient{
		client: mockClient,
	}

	err := blockchainClient.SimulateTransaction(context.Background(), common.Address{}, []byte{}, uint32(100000))
	assert.NoError(t, err)

	mockClient.AssertCalled(t, "CallContract", mock.Anything, mock.Anything, mock.Anything)
}
