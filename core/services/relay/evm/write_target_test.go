package evm_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonevm "github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	gasmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/mocks"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/report/datafeeds"
	df_processor "github.com/smartcontractkit/chainlink-evm/pkg/report/datafeeds/processor"
	por_processor "github.com/smartcontractkit/chainlink-evm/pkg/report/por/processor"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink-framework/capabilities/writetarget/report/platform"

	ocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	evmmocks "github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	evmcapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func newMockedEncodeTransmissionInfo(state uint8) ([]byte, error) {
	info := evm.TransmissionInfo{
		GasLimit:        big.NewInt(0),
		InvalidReceiver: false,
		State:           state,
		Success:         false,
		TransmissionID:  [32]byte{},
		Transmitter:     common.HexToAddress("0x0"),
	}
	var buffer bytes.Buffer

	// 1. Encode TransmissionId (bytes32)
	buffer.Write(info.TransmissionID[:])

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
	var invReceiverByte byte
	if info.InvalidReceiver {
		invReceiverByte = 1
	}
	invalidReceiverSlot := make([]byte, 31)
	invalidReceiverSlot = append(invalidReceiverSlot, invReceiverByte)
	buffer.Write(invalidReceiverSlot)

	// 5. Encode Success as bool (32 bytes)
	var successByte byte
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
	poller := lpmocks.NewLogPoller(t)

	chain.On("Start", mock.Anything).Return(nil)
	chain.On("Close").Return(nil)
	chain.On("ID").Return(big.NewInt(11155111))
	chain.On("TxManager").Return(txManager)
	chain.On("LogPoller").Return(poller)
	chain.On("LatestHead", mock.Anything).Return(commontypes.Head{Height: "99"}, nil)
	chain.On("GetChainInfo", mock.Anything).Return(commontypes.ChainInfo{}, nil)

	ht := headstest.NewTracker[*evmtypes.Head](t)
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

	lggr := logger.TestLogger(t, zapcore.DebugLevel)
	cRegistry := evmcapabilities.NewRegistry(lggr)
	relayer, err := evm.NewRelayer(lggr, chain, evm.RelayerOpts{
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

	reportMetadata := ocr3types.Metadata{
		Version:          1,
		ExecutionID:      "0102030405060708090a0b0c0d0e0f1000000000000000000000000000000000",
		Timestamp:        1620000000,
		DONID:            1,
		DONConfigVersion: 1,
		WorkflowID:       "1234567890123456789012345678901234567890123456789012345678901234",
		WorkflowName:     "123456789",
		WorkflowOwner:    "1234567890123456789012345678901234567890",
		ReportID:         hex.EncodeToString(reportID[:]),
	}

	generateReportEncoded := func(reportType string) []byte {
		feedReports := datafeeds.Reports{
			{
				FeedID:    [32]byte{0x01},
				Price:     big.NewInt(1234567890123456789),
				Timestamp: 1620000000,
			},
		}

		ccipFeedReports := datafeeds.CCIPReports{
			{
				FeedID:    [32]byte{0x01},
				Timestamp: 1620000000,
				Price:     big.NewInt(1234567890123456789),
			},
		}

		porFeedReports := datafeeds.PORReports{
			{
				DataID:    [32]byte{0x01},
				Timestamp: 1620000000,
				Bundle:    []byte{0x01, 0x02, 0x03},
			},
		}

		var feedReportsEncoded []byte

		switch reportType {
		case "ccip":
			feedReportsEncoded, err = df_processor.GetCCIPDataFeedsSchema().Pack(ccipFeedReports)
			require.NoError(t, err)
		case "por":
			feedReportsEncoded, err = por_processor.GetPORSchema().Pack(porFeedReports)
			require.NoError(t, err)
		// normal non-ccip / POR report
		default:
			feedReportsEncoded, err = df_processor.GetDataFeedsSchema().Pack(feedReports)
			require.NoError(t, err)
		}

		report := platform.Report{
			Metadata: reportMetadata,
			Data:     feedReportsEncoded,
		}

		reportEncoded, encodeErr := report.Encode()
		require.NoError(t, encodeErr)

		return reportEncoded
	}

	signatures := [][]byte{}

	mockSuccessfulTransmission := func(reportType string) {
		// This is a very error-prone way to mock an on-chain response to a GetLatestValue("getTransmissionInfo") call
		// It's a bit of a hack, but it's the best way to do it without a lot of refactoring
		mockNotStarted, mockErr := newMockedEncodeTransmissionInfo(0)
		require.NoError(t, mockErr)

		evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(mockNotStarted, nil).Once()
		evmClient.On("CodeAt", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test"), nil)

		txManager.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(commontypes.Finalized, nil).Maybe()

		mockSucceeded, mockErr2 := newMockedEncodeTransmissionInfo(1)
		require.NoError(t, mockErr2)
		evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(mockSucceeded, nil).Maybe().Once()

		txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, nil).Run(func(args mock.Arguments) {
			req := args.Get(1).(txmgr.TxRequest)
			payload := make(map[string]any)
			method := forwardABI.Methods["report"]
			err = method.Inputs.UnpackIntoMap(payload, req.EncodedPayload[4:])
			require.NoError(t, err)
			require.Equal(t, generateReportEncoded(reportType), payload["rawReport"])
			require.Equal(t, signatures, payload["signatures"])
		}).Once()
		txManager.On("GetTransactionFee", mock.Anything, mock.Anything).Return(&commonevm.TransactionFee{TransactionFee: big.NewInt(10)}, nil)
	}

	generateValidInputs := func(reportType string) *values.Map {
		validInputs, inputErr := values.NewMap(map[string]any{
			"signed_report": map[string]any{
				"report":     generateReportEncoded(reportType),
				"signatures": signatures,
				"context":    []byte{4, 5},
				"id":         reportID[:],
			},
		})
		require.NoError(t, inputErr)
		return validInputs
	}

	// default inputs/report are not CCIP / POR
	validInputs := generateValidInputs("")

	validMetadata := capabilities.RequestMetadata{
		WorkflowID:          reportMetadata.WorkflowID,
		WorkflowOwner:       reportMetadata.WorkflowOwner,
		WorkflowName:        reportMetadata.WorkflowName,
		WorkflowExecutionID: reportMetadata.ExecutionID,
	}

	validConfig, err := values.NewMap(map[string]any{
		"address":   evmCfg.EVM().Workflow().ForwarderAddress().String(),
		"processor": "evm-data-feeds",
	})
	require.NoError(t, err)

	gasLimitDefault := uint64(400_000)

	t.Run("succeeds with valid report", func(t *testing.T) {
		mockSuccessfulTransmission("")
		ctx := testutils.Context(t)
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, lggr)
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   validConfig,
			Inputs:   validInputs,
		}

		_, err = capability.Execute(ctx, req)
		require.NoError(t, err)

		findLogMatch(t, observed, "[Beholder.emit]", "attributes", "FeedUpdated")
	})

	t.Run("succeeds with valid CCIP report", func(t *testing.T) {
		mockSuccessfulTransmission("ccip")
		ctx := testutils.Context(t)
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, lggr)
		require.NoError(t, err)

		config, err := values.NewMap(map[string]any{
			"address":   evmCfg.EVM().Workflow().ForwarderAddress().String(),
			"processor": "evm-data-feeds-ccip",
		})
		require.NoError(t, err)

		// special request with properly encoded CCIP report using ccip processor
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   generateValidInputs("ccip"),
		}

		_, err = capability.Execute(ctx, req)
		require.NoError(t, err)

		findLogMatch(t, observed, "[Beholder.emit]", "attributes", "FeedUpdated")
	})

	t.Run("succeeds with valid POR report", func(t *testing.T) {
		mockSuccessfulTransmission("por")
		ctx := testutils.Context(t)
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, lggr)
		require.NoError(t, err)

		config, err := values.NewMap(map[string]any{
			"address":   evmCfg.EVM().Workflow().ForwarderAddress().String(),
			"processor": "evm-por-feeds",
		})
		require.NoError(t, err)

		// special request with properly encoded CCIP report using ccip processor
		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   generateValidInputs("por"),
		}

		_, err = capability.Execute(ctx, req)
		require.NoError(t, err)

		findLogMatch(t, observed, "[Beholder.emit]", "attributes", "FeedUpdated")
	})

	t.Run("succeeds with valid report, but logs error for missing processor", func(t *testing.T) {
		mockSuccessfulTransmission("")

		ctx := testutils.Context(t)
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

		capability, err := evm.NewWriteTarget(ctx, relayer, chain, gasLimitDefault, lggr)
		require.NoError(t, err)

		config, err := values.NewMap(map[string]any{
			"address":   evmCfg.EVM().Workflow().ForwarderAddress().String(),
			"processor": "invalid-name",
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: validMetadata,
			Config:   config,
			Inputs:   validInputs,
		}

		_, err = capability.Execute(ctx, req)
		require.NoError(t, err)

		tests.RequireLogMessage(t, observed, "no matching processor for MetaCapabilityProcessor=invalid-name")
	})

	t.Run("succeeds when report already succeeded", func(t *testing.T) {
		mockCall, err := newMockedEncodeTransmissionInfo(1)
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
		mockCall, err := newMockedEncodeTransmissionInfo(0)
		require.NoError(t, err)
		evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(mockCall, nil).Once()
		txManager.On("CreateTransaction", mock.Anything, mock.Anything).Return(txmgr.Tx{}, errors.New("TXM error")).Once()

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

		relayer, err := evm.NewRelayer(lggr, testChain, evm.RelayerOpts{
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
}

func findLogMatch(t *testing.T, observed *observer.ObservedLogs, msg string, key string, value string) {
	require.Eventually(t, func() bool {
		filteredByMsg := observed.FilterMessage(msg)
		matches := filteredByMsg.
			Filter(func(le observer.LoggedEntry) bool {
				for _, field := range le.Context {
					if field.Key == key &&
						strings.Contains(fmt.Sprint(field.Interface),
							value) {
						return true
					}
				}
				return false
			}).
			All() // => []observer.LoggedEntry
		return len(matches) > 0
	}, 30*time.Second, 1*time.Second)
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
			networkName, err := chainselectors.ExtractNetworkEnvName(tc.networkName)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedName, networkName)
		})
	}
}
