package targets_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var forwardABI = types.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func TestEvmWrite(t *testing.T) {
	chain := evmmocks.NewChain(t)

	txManager := txmmocks.NewMockEvmTxManager(t)
	chain.On("ID").Return(big.NewInt(11155111))
	chain.On("TxManager").Return(txManager)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		a := testutils.NewAddress()
		addr, err := types.NewEIP55Address(a.Hex())
		require.NoError(t, err)
		c.EVM[0].ChainWriter.FromAddress = &addr

		forwarderA := testutils.NewAddress()
		forwarderAddr, err := types.NewEIP55Address(forwarderA.Hex())
		require.NoError(t, err)
		c.EVM[0].ChainWriter.ForwarderAddress = &forwarderAddr
	})
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	chain.On("Config").Return(evmcfg)

	capability := targets.NewEvmWrite(chain, logger.TestLogger(t))
	ctx := testutils.Context(t)

	config, err := values.NewMap(map[string]any{})
	require.NoError(t, err)

	inputs, err := values.NewMap(map[string]any{
		"report":     []byte{1, 2, 3},
		"signatures": [][]byte{},
	})
	require.NoError(t, err)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "hello",
		},
		Config: config,
		Inputs: inputs,
	}

	txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, nil).Run(func(args mock.Arguments) {
		req := args.Get(1).(txmgr.TxRequest)
		payload := make(map[string]any)
		method := forwardABI.Methods["report"]
		err = method.Inputs.UnpackIntoMap(payload, req.EncodedPayload[4:])
		require.NoError(t, err)
		require.Equal(t, []byte{0x1, 0x2, 0x3}, payload["rawReport"])
		require.Equal(t, [][]byte{}, payload["signatures"])
	})

	ch, err := capability.Execute(ctx, req)
	require.NoError(t, err)

	response := <-ch
	require.Nil(t, response.Err)
}

func TestEvmWrite_EmptyReport(t *testing.T) {
	chain := evmmocks.NewChain(t)

	txManager := txmmocks.NewMockEvmTxManager(t)
	chain.On("ID").Return(big.NewInt(11155111))
	chain.On("TxManager").Return(txManager)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		a := testutils.NewAddress()
		addr, err := types.NewEIP55Address(a.Hex())
		require.NoError(t, err)
		c.EVM[0].ChainWriter.FromAddress = &addr

		forwarderA := testutils.NewAddress()
		forwarderAddr, err := types.NewEIP55Address(forwarderA.Hex())
		require.NoError(t, err)
		c.EVM[0].ChainWriter.ForwarderAddress = &forwarderAddr
	})
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	chain.On("Config").Return(evmcfg)

	capability := targets.NewEvmWrite(chain, logger.TestLogger(t))
	ctx := testutils.Context(t)

	config, err := values.NewMap(map[string]any{
		"abi":    "receive(report bytes)",
		"params": []any{"$(report)"},
	})
	require.NoError(t, err)

	inputs, err := values.NewMap(map[string]any{
		"report": nil,
	})
	require.NoError(t, err)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "hello",
		},
		Config: config,
		Inputs: inputs,
	}

	ch, err := capability.Execute(ctx, req)
	require.NoError(t, err)

	response := <-ch
	require.Nil(t, response.Err)
}
