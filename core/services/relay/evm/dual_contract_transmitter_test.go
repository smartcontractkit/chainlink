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

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var sampleAddressPrimary = testutils.NewAddress()

type mockDualTransmitter struct {
	lastPrimaryPayload   []byte
	lastSecondaryPayload []byte
}

func (m *mockDualTransmitter) SecondaryFromAddress(ctx context.Context) (gethcommon.Address, error) {
	return gethcommon.Address{}, nil
}

func (*mockDualTransmitter) FromAddress(ctx context.Context) gethcommon.Address {
	return sampleAddressPrimary
}

func (m *mockDualTransmitter) CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte, _ *txmgr.TxMeta) error {
	m.lastPrimaryPayload = payload
	return nil
}

func (m *mockDualTransmitter) CreateSecondaryEthTransaction(ctx context.Context, payload []byte, _ *txmgr.TxMeta) error {
	m.lastSecondaryPayload = payload
	return nil
}

func TestDualContractTransmitter(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	c := clienttest.NewClient(t)
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
	ot, err := NewOCRDualContractTransmitter(ctx, gethcommon.Address{}, c, contractABI, &mockDualTransmitter{}, lp, lggr, &keystest.FakeChainStore{}, WithReportToEthMetadata(reportToEvmTxMeta))
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
	from, err := ot.FromAccount(t.Context())
	require.NoError(t, err)
	assert.Equal(t, sampleAddressPrimary.String(), string(from))
}

func Test_dualContractTransmitterNoSignatures_Transmit_SignaturesAreNotTransmitted(t *testing.T) {
	t.Parallel()

	transmitter := &mockDualTransmitter{}

	ctx := context.Background()
	reportCtx := types.ReportContext{}
	report := types.Report{}
	var signatures = oneSignature()

	oc := createDualContractTransmitter(ctx, t, transmitter, WithExcludeSignatures())

	err := oc.Transmit(ctx, reportCtx, report, signatures)
	require.NoError(t, err)

	var emptyRs [][32]byte
	var emptySs [][32]byte
	var emptyVs [32]byte
	emptySignaturesPayloadPrimary, err := oc.contractABI.Pack("transmit", evmutil.RawReportContext(reportCtx), []byte(report), emptyRs, emptySs, emptyVs)
	require.NoError(t, err)
	emptySignaturesPayloadSecondary, err := oc.dualTransmissionABI.Pack("transmitSecondary", evmutil.RawReportContext(reportCtx), []byte(report), emptyRs, emptySs, emptyVs)
	require.NoError(t, err)
	require.Equal(t, transmitter.lastPrimaryPayload, emptySignaturesPayloadPrimary, "primary payload not equal")
	require.Equal(t, transmitter.lastSecondaryPayload, emptySignaturesPayloadSecondary, "secondary payload not equal")
}

func Test_dualContractTransmitter_Transmit_SignaturesAreTransmitted(t *testing.T) {
	t.Parallel()

	transmitter := &mockDualTransmitter{}

	ctx := context.Background()
	reportCtx := types.ReportContext{}
	report := types.Report{}
	var signatures = oneSignature()

	oc := createDualContractTransmitter(ctx, t, transmitter)

	err := oc.Transmit(ctx, reportCtx, report, signatures)
	require.NoError(t, err)

	rs, ss, vs := signaturesAsPayload(t, signatures)
	withSignaturesPayloadPrimary, err := oc.contractABI.Pack("transmit", evmutil.RawReportContext(reportCtx), []byte(report), rs, ss, vs)
	require.NoError(t, err)
	withSignaturesPayloadSecondary, err := oc.dualTransmissionABI.Pack("transmitSecondary", evmutil.RawReportContext(reportCtx), []byte(report), rs, ss, vs)
	require.NoError(t, err)
	require.Equal(t, transmitter.lastPrimaryPayload, withSignaturesPayloadPrimary, "primary payload not equal")
	require.Equal(t, transmitter.lastSecondaryPayload, withSignaturesPayloadSecondary, "secondary payload not equal")
}

func createDualContractTransmitter(ctx context.Context, t *testing.T, transmitter Transmitter, ops ...OCRTransmitterOption) *dualContractTransmitter {
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	require.NoError(t, err)
	lp := lpmocks.NewLogPoller(t)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	contractTransmitter, err := NewOCRDualContractTransmitter(
		ctx,
		gethcommon.Address{},
		clienttest.NewClient(t),
		contractABI,
		transmitter,
		lp,
		logger.Test(t),
		&keystest.FakeChainStore{},
		ops...,
	)
	require.NoError(t, err)
	return contractTransmitter
}
