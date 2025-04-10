package evm_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	gasmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/mocks"
	dftypes "github.com/smartcontractkit/chainlink-evm/pkg/report/datafeeds"
	"github.com/smartcontractkit/chainlink-evm/pkg/report/platform"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	evmcapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	pollermocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	relayevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func newMockedEncodeTransmissionInfo(state uint8) ([]byte, error) {
	info := evm.TransmissionInfo{
		GasLimit:        big.NewInt(0),
		InvalidReceiver: false,
		State:           state,
		Success:         false,
		TransmissionId:  [32]byte{},
		Transmitter:     common.HexToAddress("0x0"),
	}
	var buffer bytes.Buffer

	// 1. Encode TransmissionId (bytes32)
	buffer.Write(info.TransmissionId[:])

	// 2. Encode State (uint8, ABI pads to 32 bytes: 31 zeros + 1 byte)
	stateSlot := make([]byte, 31)
	stateSlot = append(stateSlot, info.State)
	buffer.Write(stateSlot)

	// 3. Encode Transmitter (address): address is 20 bytes; pad left with 12 zeros.
	txBytes := info.Transmitter.Bytes()
	// Ensure it is 20 bytes; if not, left-pad it (common.HexToAddress always returns 20 bytes).
	paddedTx := make([]byte, 32-len(txBytes))
	buffer.Write(paddedTx)
	buffer.Write(txBytes)

	// 4. Encode InvalidReceiver as bool (32 bytes: 31 zeros then byte(0) or byte(1))
	var invReceiverByte byte = 0
	if info.InvalidReceiver {
		invReceiverByte = 1
	}
	invalidReceiverSlot := make([]byte, 31)
	invalidReceiverSlot = append(invalidReceiverSlot, invReceiverByte)
	buffer.Write(invalidReceiverSlot)

	// 5. Encode Success as bool (32 bytes)
	var successByte byte = 0
	if info.Success {
		successByte = 1
	}
	successSlot := make([]byte, 31)
	successSlot = append(successSlot, successByte)
	buffer.Write(successSlot)

	// 6. Encode GasLimit as a uint (for uint80, still ABI-encoded in a 32-byte slot)
	gasLimitBytes := info.GasLimit.Bytes()
	// Left-pad the gas limit to 32 bytes.
	paddedGasLimit := make([]byte, 32-len(gasLimitBytes))
	buffer.Write(paddedGasLimit)
	buffer.Write(gasLimitBytes)

	return buffer.Bytes(), nil
}

func TestEvmWrite(t *testing.T) {
	chain := evmmocks.NewChain(t)
	txManager := txmmocks.NewMockEvmTxManager(t)
	evmClient := clienttest.NewClient(t)
	poller := pollermocks.NewLogPoller(t)

	// This is a very error-prone way to mock an on-chain response to a GetLatestValue("getTransmissionInfo") call
	// It's a bit of a hack, but it's the best way to do it without a lot of refactoring
	mockCall, err := newMockedEncodeTransmissionInfo(0)
	require.NoError(t, err)

	evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(mockCall, nil).Times(3)
	evmClient.On("CodeAt", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test"), nil)

	txManager.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(commontypes.Finalized, nil).Maybe()

	chain.On("Start", mock.Anything).Return(nil)
	chain.On("Close").Return(nil)
	chain.On("ID").Return(big.NewInt(11155111))
	chain.On("TxManager").Return(txManager)
	chain.On("LogPoller").Return(poller)
	chain.On("LatestHead", mock.Anything).Return(commontypes.Head{Height: "99"}, nil)

	ht := headstest.NewTracker[*evmtypes.Head, common.Hash](t)
	ht.On("LatestAndFinalizedBlock", mock.Anything).Return(&evmtypes.Head{Number: 99}, &evmtypes.Head{}, nil)
	chain.On("HeadTracker").Return(ht)

	chain.On("Client").Return(evmClient)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		a := testutils.NewAddress()
		addr, err2 := evmtypes.NewEIP55Address(a.Hex())
		require.NoError(t, err2)
		c.EVM[0].Workflow.FromAddress = &addr

		forwarderA := testutils.NewAddress()
		forwarderAddr, err2 := evmtypes.NewEIP55Address(forwarderA.Hex())
		require.NoError(t, err2)
		c.EVM[0].Workflow.ForwarderAddress = &forwarderAddr
	})
	evmCfg := evmtest.NewChainScopedConfig(t, cfg)
	ge := gasmocks.NewEvmFeeEstimator(t)

	chain.On("Config").Return(evmCfg)
	chain.On("GasEstimator").Return(ge)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)

	lggr := logger.TestLogger(t)
	cRegistry := evmcapabilities.NewRegistry(lggr)
	relayer, err := relayevm.NewRelayer(lggr, chain, relayevm.RelayerOpts{
		DS:                   db,
		EVMKeystore:          keys.NewChainStore(keystore.NewEthSigner(keyStore.Eth(), chain.ID()), chain.ID()),
		CSAKeystore:          &keystore.CSASigner{CSA: keyStore.CSA()},
		CapabilitiesRegistry: cRegistry,
	})
	require.NoError(t, err)
	servicetest.Run(t, relayer)
	registeredCapabilities, err := cRegistry.List(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, registeredCapabilities, 1) // WriteTarget should be added to the registry

	reportID := [2]byte{0x00, 0x01}
	reportMetadata := dftypes.Metadata{
		Version:             1,
		WorkflowExecutionID: [32]byte{},
		Timestamp:           0,
		DonID:               0,
		DonConfigVersion:    0,
		WorkflowCID:         [32]byte{},
		WorkflowName:        [10]byte{},
		WorkflowOwner:       [20]byte{},
		ReportID:            reportID,
	}

	feedReports := dftypes.Reports{
		{
			FeedID:    [32]byte{0x01},
			Price:     big.NewInt(1234567890123456789),
			Timestamp: 1620000000,
		},
	}

	feedReportsEncoded, err := dftypes.GetSchema().Pack(feedReports)
	require.NoError(t, err)

	report := platform.Report{
		Metadata: reportMetadata,
		Data:     feedReportsEncoded,
	}

	reportEncoded, err := report.Encode()
	require.NoError(t, err)

	signatures := [][]byte{}

	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     reportEncoded,
			"signatures": signatures,
			"context":    []byte{4, 5},
			"id":         reportID[:],
		},
	})
	require.NoError(t, err)

	validMetadata := capabilities.RequestMetadata{
		WorkflowID:          hex.EncodeToString(reportMetadata.WorkflowCID[:]),
		WorkflowOwner:       hex.EncodeToString(reportMetadata.WorkflowOwner[:]),
		WorkflowName:        hex.EncodeToString(reportMetadata.WorkflowName[:]),
		WorkflowExecutionID: hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]),
	}

	validConfig, err := values.NewMap(map[string]any{
		"Address": evmCfg.EVM().Workflow().ForwarderAddress().String(),
	})
	require.NoError(t, err)

	txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, nil).Run(func(args mock.Arguments) {
		req := args.Get(1).(txmgr.TxRequest)
		payload := make(map[string]any)
		method := forwardABI.Methods["report"]
		err = method.Inputs.UnpackIntoMap(payload, req.EncodedPayload[4:])
		require.NoError(t, err)
		require.Equal(t, reportEncoded, payload["rawReport"])
		require.Equal(t, signatures, payload["signatures"])
	}).Once()

	gasLimitDefault := uint64(400_000)

	t.Run("succeeds with valid report", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, lggr)
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   validConfig,
			Inputs:   validInputs,
		}

		_, err = capability.Execute(ctx, req)
		require.NoError(t, err)
	})

	t.Run("fails with invalid config", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, logger.TestLogger(t))
		require.NoError(t, err)

		invalidConfig, err := values.NewMap(map[string]any{
			"Address": "invalid-address",
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   invalidConfig,
			Inputs:   validInputs,
		}

		_, err = capability.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("fails when TXM CreateTransaction returns error", func(t *testing.T) {
		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, logger.TestLogger(t))
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   validConfig,
			Inputs:   validInputs,
		}

		txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, errors.New("TXM error"))

		_, err = capability.Execute(ctx, req)
		require.Error(t, err)
	})

	t.Run("Relayer fails to start WriteTarget capability on missing config", func(t *testing.T) {
		ctx := testutils.Context(t)
		testChain := evmmocks.NewChain(t)
		testCfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].Workflow.FromAddress = nil
			forwarderA := testutils.NewAddress()
			forwarderAddr, err2 := evmtypes.NewEIP55Address(forwarderA.Hex())
			require.NoError(t, err2)
			c.EVM[0].Workflow.ForwarderAddress = &forwarderAddr
		})
		testChain.On("Start", mock.Anything).Return(nil)
		testChain.On("Close").Return(nil)
		testChain.On("ID").Return(big.NewInt(11155111))
		testChain.On("Config").Return(evmtest.NewChainScopedConfig(t, testCfg))
		capabilityRegistry := evmcapabilities.NewRegistry(lggr)

		relayer, err := relayevm.NewRelayer(lggr, testChain, relayevm.RelayerOpts{
			DS:                   db,
			EVMKeystore:          keys.NewChainStore(keystore.NewEthSigner(keyStore.Eth(), chain.ID()), chain.ID()),
			CSAKeystore:          &keystore.CSASigner{CSA: keyStore.CSA()},
			CapabilitiesRegistry: capabilityRegistry,
		})
		require.NoError(t, err)
		servicetest.Run(t, relayer)

		l, err := capabilityRegistry.List(ctx)
		require.NoError(t, err)

		assert.Empty(t, l)
	})

	t.Run("succeeds when report already succeeded", func(t *testing.T) {
		mockCall, err = newMockedEncodeTransmissionInfo(1)
		require.NoError(t, err)

		evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(mockCall, nil).Once()

		ctx := testutils.Context(t)
		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, lggr)
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   validConfig,
			Inputs:   validInputs,
		}

		_, err = capability.Execute(ctx, req)
		require.NoError(t, err)
	})
}

func TestExtractNetwork(t *testing.T) {
	testCases := []struct {
		networkName  string
		expectedName string
		expectedErr  bool
	}{
		{
			networkName:  "ethereum-testnet-goerli",
			expectedName: "testnet",
			expectedErr:  false,
		},
		{
			networkName:  "ethereum-mainnet",
			expectedName: "mainnet",
			expectedErr:  false,
		},
		{
			networkName:  "polygon-devnet",
			expectedName: "devnet",
			expectedErr:  false,
		},
		{
			networkName:  "ethereum_test",
			expectedName: "",
			expectedErr:  true,
		},
		{
			networkName:  "ethereum",
			expectedName: "",
			expectedErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.networkName, func(t *testing.T) {
			networkName, err := evm.ExtractNetwork(tc.networkName)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedName, networkName)
		})
	}
}
