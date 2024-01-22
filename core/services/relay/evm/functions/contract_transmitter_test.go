package functions_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
)

type mockTransmitter struct {
	toAddress gethcommon.Address
}

func (m *mockTransmitter) CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte, _ *txmgr.TxMeta) error {
	m.toAddress = toAddress
	return nil
}
func (mockTransmitter) FromAddress() gethcommon.Address { return testutils.NewAddress() }

func TestContractTransmitter_LatestConfigDigestAndEpoch(t *testing.T) {
	t.Parallel()

	digestStr := "000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776"
	lggr := logger.TestLogger(t)
	c := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	digestAndEpochDontScanLogs, err := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000000" + // scan logs = false
			digestStr +
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	require.NoError(t, err)
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(digestAndEpochDontScanLogs, nil).Once()
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	require.NoError(t, err)
	lp.On("RegisterFilter", mock.Anything).Return(nil)

	functionsTransmitter, err := functions.NewFunctionsContractTransmitter(c, contractABI, &mockTransmitter{}, lp, lggr, func(b []byte) (*txmgr.TxMeta, error) {
		return &txmgr.TxMeta{}, nil
	}, 1)
	require.NoError(t, err)
	require.NoError(t, functionsTransmitter.UpdateRoutes(gethcommon.Address{}, gethcommon.Address{}))

	digest, epoch, err := functionsTransmitter.LatestConfigDigestAndEpoch(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, digestStr, hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)
}

func TestContractTransmitter_Transmit_V1(t *testing.T) {
	t.Parallel()

	contractVersion := uint32(1)
	configuredDestAddress, coordinatorAddress := testutils.NewAddress(), testutils.NewAddress()
	lggr := logger.TestLogger(t)
	c := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	lp.On("RegisterFilter", mock.Anything).Return(nil)

	ocrTransmitter := mockTransmitter{}
	ot, err := functions.NewFunctionsContractTransmitter(c, contractABI, &ocrTransmitter, lp, lggr, func(b []byte) (*txmgr.TxMeta, error) {
		return &txmgr.TxMeta{}, nil
	}, contractVersion)
	require.NoError(t, err)
	require.NoError(t, ot.UpdateRoutes(configuredDestAddress, configuredDestAddress))

	reqId, err := hex.DecodeString("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	require.NoError(t, err)
	processedRequests := []*encoding.ProcessedRequest{
		{
			RequestID:           reqId,
			CoordinatorContract: coordinatorAddress.Bytes(),
		},
	}
	codec, err := encoding.NewReportCodec(contractVersion)
	require.NoError(t, err)
	reportBytes, err := codec.EncodeReport(processedRequests)
	require.NoError(t, err)

	// success
	require.NoError(t, ot.Transmit(testutils.Context(t), ocrtypes.ReportContext{}, reportBytes, []ocrtypes.AttributedOnchainSignature{}))
	require.Equal(t, coordinatorAddress, ocrTransmitter.toAddress)

	// failure on too many signatures
	signatures := []ocrtypes.AttributedOnchainSignature{}
	for i := 0; i < 33; i++ {
		signatures = append(signatures, ocrtypes.AttributedOnchainSignature{})
	}
	require.Error(t, ot.Transmit(testutils.Context(t), ocrtypes.ReportContext{}, reportBytes, signatures))
}

func TestContractTransmitter_Transmit_V1_CoordinatorMismatch(t *testing.T) {
	t.Parallel()

	contractVersion := uint32(1)
	configuredDestAddress, coordinatorAddress1, coordinatorAddress2 := testutils.NewAddress(), testutils.NewAddress(), testutils.NewAddress()
	lggr := logger.TestLogger(t)
	c := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	lp.On("RegisterFilter", mock.Anything).Return(nil)

	ocrTransmitter := mockTransmitter{}
	ot, err := functions.NewFunctionsContractTransmitter(c, contractABI, &ocrTransmitter, lp, lggr, func(b []byte) (*txmgr.TxMeta, error) {
		return &txmgr.TxMeta{}, nil
	}, contractVersion)
	require.NoError(t, err)
	require.NoError(t, ot.UpdateRoutes(configuredDestAddress, configuredDestAddress))

	reqId1, err := hex.DecodeString("110102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	require.NoError(t, err)
	reqId2, err := hex.DecodeString("220102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	require.NoError(t, err)
	processedRequests := []*encoding.ProcessedRequest{
		{
			RequestID:           reqId1,
			CoordinatorContract: coordinatorAddress1.Bytes(),
		},
		{
			RequestID:           reqId2,
			CoordinatorContract: coordinatorAddress2.Bytes(),
		},
	}
	codec, err := encoding.NewReportCodec(contractVersion)
	require.NoError(t, err)
	reportBytes, err := codec.EncodeReport(processedRequests)
	require.NoError(t, err)

	require.NoError(t, ot.Transmit(testutils.Context(t), ocrtypes.ReportContext{}, reportBytes, []ocrtypes.AttributedOnchainSignature{}))
	require.Equal(t, coordinatorAddress1, ocrTransmitter.toAddress)
}
