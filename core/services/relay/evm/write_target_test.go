package evm_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	evmcapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	relayevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

var forwardABI = types.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func TestEvmWrite(t *testing.T) {
	chain := evmmocks.NewChain(t)
	txManager := txmmocks.NewMockEvmTxManager(t)
	evmClient := evmclimocks.NewClient(t)

	// This probably isn't the best way to do this, but couldn't find a simpler way to mock the CallContract response
	var mockCall []byte
	for i := 0; i < 32; i++ {
		mockCall = append(mockCall, byte(0))
	}
	evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(mockCall, nil).Maybe()

	chain.On("ID").Return(big.NewInt(11155111))
	chain.On("TxManager").Return(txManager)
	chain.On("LogPoller").Return(nil)
	chain.On("Client").Return(evmClient)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		a := testutils.NewAddress()
		addr, err2 := types.NewEIP55Address(a.Hex())
		require.NoError(t, err2)
		c.EVM[0].Workflow.FromAddress = &addr

		forwarderA := testutils.NewAddress()
		forwarderAddr, err2 := types.NewEIP55Address(forwarderA.Hex())
		require.NoError(t, err2)
		c.EVM[0].Workflow.ForwarderAddress = &forwarderAddr
	})
	evmCfg := evmtest.NewChainScopedConfig(t, cfg)

	chain.On("Config").Return(evmCfg)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)

	lggr := logger.TestLogger(t)
	relayer, err := relayevm.NewRelayer(lggr, chain, relayevm.RelayerOpts{
		DS:                   db,
		CSAETHKeystore:       keyStore,
		CapabilitiesRegistry: evmcapabilities.NewRegistry(lggr),
	})
	require.NoError(t, err)

	txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, nil).Run(func(args mock.Arguments) {
		req := args.Get(1).(txmgr.TxRequest)
		payload := make(map[string]any)
		method := forwardABI.Methods["report"]
		err = method.Inputs.UnpackIntoMap(payload, req.EncodedPayload[4:])
		require.NoError(t, err)
		require.Equal(t, []byte{0x1, 0x2, 0x3}, payload["rawReport"])
		require.Equal(t, [][]byte{}, payload["signatures"])
	}).Once()

	t.Run("succeeds with valid report", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, lggr)
		require.NoError(t, err)

		config, err := values.NewMap(map[string]any{
			"Address": evmCfg.EVM().Workflow().ForwarderAddress().String(),
		})
		require.NoError(t, err)

		inputs, err := values.NewMap(map[string]any{
			"signed_report": map[string]any{
				"report":     []byte{1, 2, 3},
				"signatures": [][]byte{},
				"context":    []byte{4, 5},
				"id":         []byte{9, 9},
			},
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: config,
			Inputs: inputs,
		}

		ch, err := capability.Execute(ctx, req)
		require.NoError(t, err)

		response := <-ch
		require.Nil(t, response.Err)
	})

	t.Run("succeeds with empty report", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, logger.TestLogger(t))
		require.NoError(t, err)

		config, err := values.NewMap(map[string]any{
			"Address": evmCfg.EVM().Workflow().ForwarderAddress().String(),
		})
		require.NoError(t, err)

		inputs, err := values.NewMap(map[string]any{
			"signed_report": map[string]any{
				"report": nil,
			},
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: config,
			Inputs: inputs,
		}

		ch, err := capability.Execute(ctx, req)
		require.NoError(t, err)

		response := <-ch
		require.Nil(t, response.Err)
	})

	t.Run("fails with invalid config", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, logger.TestLogger(t))
		require.NoError(t, err)

		invalidConfig, err := values.NewMap(map[string]any{
			"Address": "invalid-address",
		})
		require.NoError(t, err)

		inputs, err := values.NewMap(map[string]any{
			"signed_report": map[string]any{
				"report": nil,
			},
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: invalidConfig,
			Inputs: inputs,
		}

		_, err = capability.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails when TXM CreateTransaction returns error", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, logger.TestLogger(t))
		require.NoError(t, err)

		config, err := values.NewMap(map[string]any{
			"Address": evmCfg.EVM().Workflow().ForwarderAddress().String(),
		})
		require.NoError(t, err)

		inputs, err := values.NewMap(map[string]any{
			"signed_report": map[string]any{
				"report":     []byte{1, 2, 3},
				"signatures": [][]byte{},
				"context":    []byte{4, 5},
				"id":         []byte{9, 9},
			},
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: "test-id",
			},
			Config: config,
			Inputs: inputs,
		}

		txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, errors.New("TXM error"))

		_, err = capability.Execute(ctx, req)
		require.Error(t, err)
	})
}
