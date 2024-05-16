package functions_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
)

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func TestContractTransmitter_LatestConfigDigestAndEpoch(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

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
	txm := txmmocks.NewMockEvmTxManager(t)
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	strategy := newMockTxStrategy(t)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)

	functionsTransmitter, err := functions.NewFunctionsContractTransmitter(
		c,
		contractABI,
		lp,
		lggr,
		1,
		txm,
		[]gethcommon.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)
	require.NoError(t, functionsTransmitter.UpdateRoutes(ctx, gethcommon.Address{}, gethcommon.Address{}))

	digest, epoch, err := functionsTransmitter.LatestConfigDigestAndEpoch(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, digestStr, hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)
}

func TestContractTransmitter_Transmit_V1(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	contractVersion := uint32(1)
	configuredDestAddress, coordinatorAddress := testutils.NewAddress(), testutils.NewAddress()
	lggr := logger.TestLogger(t)
	c := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	txm := txmmocks.NewMockEvmTxManager(t)
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	strategy := newMockTxStrategy(t)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)

	ot, err := functions.NewFunctionsContractTransmitter(
		c,
		contractABI,
		lp,
		lggr,
		contractVersion,
		txm,
		[]gethcommon.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)
	require.NoError(t, ot.UpdateRoutes(ctx, configuredDestAddress, configuredDestAddress))

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
	rawReportCtx := evmutil.RawReportContext(ocrtypes.ReportContext{})
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	payload, err := contractABI.Pack("transmit", rawReportCtx, reportBytes, rs, ss, vs)
	require.NoError(t, err)

	// success
	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress,
		ToAddress:        coordinatorAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: gethcommon.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	require.NoError(t, ot.Transmit(testutils.Context(t), ocrtypes.ReportContext{}, reportBytes, []ocrtypes.AttributedOnchainSignature{}))

	// failure on too many signatures
	signatures := []ocrtypes.AttributedOnchainSignature{}
	for i := 0; i < 33; i++ {
		signatures = append(signatures, ocrtypes.AttributedOnchainSignature{})
	}
	require.Error(t, ot.Transmit(testutils.Context(t), ocrtypes.ReportContext{}, reportBytes, signatures))
}

func TestContractTransmitter_Transmit_V1_CoordinatorMismatch(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	contractVersion := uint32(1)
	configuredDestAddress, coordinatorAddress1, coordinatorAddress2 := testutils.NewAddress(), testutils.NewAddress(), testutils.NewAddress()
	lggr := logger.TestLogger(t)
	c := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	txm := txmmocks.NewMockEvmTxManager(t)
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	strategy := newMockTxStrategy(t)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)

	ot, err := functions.NewFunctionsContractTransmitter(
		c,
		contractABI,
		lp,
		lggr,
		contractVersion,
		txm,
		[]gethcommon.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)
	require.NoError(t, ot.UpdateRoutes(ctx, configuredDestAddress, configuredDestAddress))

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
	rawReportCtx := evmutil.RawReportContext(ocrtypes.ReportContext{})
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	payload, err := contractABI.Pack("transmit", rawReportCtx, reportBytes, rs, ss, vs)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress,
		ToAddress:        coordinatorAddress1,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: gethcommon.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	require.NoError(t, ot.Transmit(testutils.Context(t), ocrtypes.ReportContext{}, reportBytes, []ocrtypes.AttributedOnchainSignature{}))
}
