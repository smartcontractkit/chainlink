package eal

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/capital-markets-projects/lib/web/jsonrpc"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlockchainClientEstimateGas(t *testing.T) {
	lggr := logger.TestLogger(t)
	fromAddress := "0x469aA2CD13e037DC5236320783dCfd0e641c0559"
	fromAddresses := []ethkey.EIP55Address{ethkey.EIP55Address(fromAddress)}
	ks := ksmocks.NewEth(t)
	ks.On("GetRoundRobinAddress", testutils.FixtureChainID, mock.Anything).Maybe().Return(common.HexToAddress(fromAddress), nil)
	txm := new(txmmocks.MockEvmTxManager)

	t.Run("happy case", func(t *testing.T) {
		mockClient := new(mocks.Client)
		mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(100_000), nil).Once()

		blockchainClient := &BlockchainClient{
			client:        mockClient,
			txm:           txm,
			lggr:          lggr,
			gethks:        ks,
			fromAddresses: fromAddresses,
			chainID:       testutils.FixtureChainID.Uint64(),
			maxGasLimit:   1_000_000,
		}

		ctx := testutils.Context(t)
		gasLimit, err := blockchainClient.EstimateGas(ctx, common.Address{}, []byte{})
		assert.NoError(t, err)
		assert.Equal(t, uint32(100_000), gasLimit)

		// Verify that the EstimateGas method on the mock client was called with the expected arguments
		mockClient.AssertCalled(t, "EstimateGas", ctx, mock.Anything)
	})

	t.Run("execution reverted", func(t *testing.T) {
		mockClient := new(mocks.Client)
		mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(0), errors.New(jsonrpc.ExecutionRevertedErrorMsg)).Once()

		blockchainClient := &BlockchainClient{
			client:        mockClient,
			txm:           txm,
			lggr:          lggr,
			gethks:        ks,
			fromAddresses: fromAddresses,
			chainID:       testutils.FixtureChainID.Uint64(),
			maxGasLimit:   1_000_000,
		}

		ctx := testutils.Context(t)
		_, err := blockchainClient.EstimateGas(ctx, common.Address{}, []byte{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), jsonrpc.ExecutionRevertedErrorMsg)

		// Verify that the EstimateGas method on the mock client was called with the expected arguments
		mockClient.AssertCalled(t, "EstimateGas", ctx, mock.Anything)
	})

	t.Run("gas limit exceeded", func(t *testing.T) {
		mockClient := new(mocks.Client)
		mockClient.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(100_000_000), nil).Once()

		blockchainClient := &BlockchainClient{
			client:        mockClient,
			txm:           txm,
			lggr:          lggr,
			gethks:        ks,
			fromAddresses: fromAddresses,
			chainID:       testutils.FixtureChainID.Uint64(),
			maxGasLimit:   1_000_000,
		}

		ctx := testutils.Context(t)
		_, err := blockchainClient.EstimateGas(ctx, common.Address{}, []byte{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), jsonrpc.EstimateGasExceededErrorMsg)

		// Verify that the EstimateGas method on the mock client was called with the expected arguments
		mockClient.AssertCalled(t, "EstimateGas", ctx, mock.Anything)
	})
}

func TestBlockchainClientSimulateTransaction(t *testing.T) {
	lggr := logger.TestLogger(t)
	fromAddress := "0x469aA2CD13e037DC5236320783dCfd0e641c0559"
	fromAddresses := []ethkey.EIP55Address{ethkey.EIP55Address(fromAddress)}
	ks := ksmocks.NewEth(t)
	ks.On("GetRoundRobinAddress", testutils.FixtureChainID, mock.Anything).Maybe().Return(common.HexToAddress(fromAddress), nil)
	txm := new(txmmocks.MockEvmTxManager)

	t.Run("happy case", func(t *testing.T) {
		mockClient := new(mocks.Client)
		result := []byte("result")
		mockClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(result, nil).Once()
		blockchainClient := &BlockchainClient{
			client:        mockClient,
			txm:           txm,
			lggr:          lggr,
			gethks:        ks,
			fromAddresses: fromAddresses,
			chainID:       testutils.FixtureChainID.Uint64(),
			maxGasLimit:   1_000_000,
		}
		ctx := testutils.Context(t)
		err := blockchainClient.SimulateTransaction(ctx, common.Address{}, []byte{}, uint32(100_000))
		assert.NoError(t, err)
	})

	t.Run("execution reverted", func(t *testing.T) {
		mockClient := new(mocks.Client)
		mockClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(jsonrpc.ExecutionRevertedErrorMsg)).Once()
		blockchainClient := &BlockchainClient{
			client:        mockClient,
			txm:           txm,
			lggr:          lggr,
			gethks:        ks,
			fromAddresses: fromAddresses,
			chainID:       testutils.FixtureChainID.Uint64(),
			maxGasLimit:   1_000_000,
		}
		ctx := testutils.Context(t)
		err := blockchainClient.SimulateTransaction(ctx, common.Address{}, []byte{}, uint32(100_000))
		assert.Error(t, err)
	})
}
