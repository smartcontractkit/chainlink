package evm

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var sampleAddress = testutils.NewAddress()

type mockTransmitter struct {
	lastPayload []byte
}

func (m *mockTransmitter) CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte, _ *txmgr.TxMeta) error {
	m.lastPayload = payload
	return nil
}

func (*mockTransmitter) FromAddress() gethcommon.Address { return sampleAddress }

func TestContractTransmitter(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	c := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	ctx := testutils.Context(t)
	// scanLogs = false
	digestAndEpochDontScanLogs, _ := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000000" + // false
			"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(digestAndEpochDontScanLogs, nil).Once()
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	reportToEvmTxMeta := func(b []byte) (*txmgr.TxMeta, error) {
		return &txmgr.TxMeta{}, nil
	}
	ot, err := NewOCRContractTransmitter(ctx, gethcommon.Address{}, c, contractABI, &mockTransmitter{}, lp, lggr,
		WithReportToEthMetadata(reportToEvmTxMeta))
	require.NoError(t, err)
	digest, epoch, err := ot.LatestConfigDigestAndEpoch(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, "000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776", hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)

	// scanLogs = true
	digestAndEpochScanLogs, _ := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000001" + // true
			"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(digestAndEpochScanLogs, nil).Once()
	transmitted2, _ := hex.DecodeString(
		"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc777" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	lp.On("LatestLogByEventSigWithConfs",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&logpoller.Log{
		Data: transmitted2,
	}, nil)
	digest, epoch, err = ot.LatestConfigDigestAndEpoch(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, "000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc777", hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)
	from, err := ot.FromAccount()
	require.NoError(t, err)
	assert.Equal(t, sampleAddress.String(), string(from))
}

func Test_contractTransmitterNoSignatures_Transmit_SignaturesAreNotTransmitted(t *testing.T) {
	t.Parallel()

	transmitter := &mockTransmitter{}

	ctx := context.Background()
	reportCtx := types.ReportContext{}
	report := types.Report{}
	var signatures = oneSignature()

	oc := createContractTransmitter(ctx, t, transmitter, WithExcludeSignatures())

	err := oc.Transmit(ctx, reportCtx, report, signatures)
	require.NoError(t, err)

	var emptyRs [][32]byte
	var emptySs [][32]byte
	var emptyVs [32]byte
	emptySignaturesPayload, err := oc.contractABI.Pack("transmit", evmutil.RawReportContext(reportCtx), []byte(report), emptyRs, emptySs, emptyVs)
	require.NoError(t, err)
	require.Equal(t, transmitter.lastPayload, emptySignaturesPayload)
}

func Test_contractTransmitter_Transmit_SignaturesAreTransmitted(t *testing.T) {
	t.Parallel()

	transmitter := &mockTransmitter{}

	ctx := context.Background()
	reportCtx := types.ReportContext{}
	report := types.Report{}
	var signatures = oneSignature()

	oc := createContractTransmitter(ctx, t, transmitter)

	err := oc.Transmit(ctx, reportCtx, report, signatures)
	require.NoError(t, err)

	rs, ss, vs := signaturesAsPayload(t, signatures)
	withSignaturesPayload, err := oc.contractABI.Pack("transmit", evmutil.RawReportContext(reportCtx), []byte(report), rs, ss, vs)
	require.NoError(t, err)
	require.Equal(t, transmitter.lastPayload, withSignaturesPayload)
}

func signaturesAsPayload(t *testing.T, signatures []ocrtypes.AttributedOnchainSignature) ([][32]byte, [][32]byte, [32]byte) {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	r, s, v, err := evmutil.SplitSignature(signatures[0].Signature)
	require.NoError(t, err)
	rs = append(rs, r)
	ss = append(ss, s)
	vs[0] = v
	return rs, ss, vs
}

func oneSignature() []ocrtypes.AttributedOnchainSignature {
	signaturesData := make([]byte, 65)
	signaturesData[9] = 8
	signaturesData[7] = 6
	return []libocr.AttributedOnchainSignature{{Signature: signaturesData, Signer: commontypes.OracleID(54)}}
}

func createContractTransmitter(ctx context.Context, t *testing.T, transmitter Transmitter, ops ...OCRTransmitterOption) *contractTransmitter {
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	require.NoError(t, err)
	lp := lpmocks.NewLogPoller(t)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	contractTransmitter, err := NewOCRContractTransmitter(
		ctx,
		gethcommon.Address{},
		evmclimocks.NewClient(t),
		contractABI,
		transmitter,
		lp,
		logger.TestLogger(t),
		ops...,
	)
	require.NoError(t, err)
	return contractTransmitter
}
