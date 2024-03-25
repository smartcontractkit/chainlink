package test

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// CommitStoreReader is a test implementation of the CommitStoreReader interface
var CommitStoreReader = staticCommitStoreReader{
	staticCommitStoreReaderConfig: staticCommitStoreReaderConfig{
		// change config test data
		changeConfigRequest: changeConfigRequest{
			onchainConfig:  []byte("onchainConfig"),
			offchainConfig: []byte("offchainConfig"),
		},
		changeConfigResponse: ccip.Address("new commit store address"),

		// en/decode commit test data
		decodeCommitReportRequest: []byte("encoded commit"),
		decodeCommitReportResponse: ccip.CommitStoreReport{
			TokenPrices: []ccip.TokenPrice{
				{
					Token: ccip.Address("token address"),
					Value: big.NewInt(1),
				},
				{
					Token: ccip.Address("token address 2"),
					Value: big.NewInt(2),
				},
			},
			GasPrices: []ccip.GasPrice{
				{
					DestChainSelector: 1,
					Value:             big.NewInt(1),
				},
				{
					DestChainSelector: 2,
					Value:             big.NewInt(2),
				},
			},
			Interval: ccip.CommitStoreInterval{
				Min: 1,
				Max: 99,
			},
			MerkleRoot: [32]byte{0: 1, 31: 7},
		},

		// gas price estimator test data
		gasPriceEstimatorResponse: GasPriceEstimatorCommit,

		// get accepted commit reports gte timestamp test data
		getAcceptedCommitReportsGteTimestampRequest: getAcceptedCommitReportsGteTimestampRequest{
			timestamp:     time.Unix(10000, 7).UTC(),
			confirmations: 1,
		},
		getAcceptedCommitReportsGteTimestampResponse: []ccip.CommitStoreReportWithTxMeta{
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 10000,
					BlockNumber:             1,
					TxHash:                  "first accepted hash",
					LogIndex:                1,
				},
				CommitStoreReport: ccip.CommitStoreReport{
					TokenPrices: []ccip.TokenPrice{
						{
							Token: ccip.Address("first accepted token address1"),
							Value: big.NewInt(7),
						},
						{
							Token: ccip.Address("first accepted token address2"),
							Value: big.NewInt(8),
						},
					},
					GasPrices: []ccip.GasPrice{
						{
							DestChainSelector: 7,
							Value:             big.NewInt(7),
						},
						{
							DestChainSelector: 8,
							Value:             big.NewInt(8),
						},
					},
					Interval: ccip.CommitStoreInterval{
						Min: 700,
						Max: 799,
					},
					MerkleRoot: [32]byte{1: 11, 2: 13},
				},
			},
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 20000,
					BlockNumber:             2,
					TxHash:                  "second accepted hash 2",
					LogIndex:                1,
				},
				CommitStoreReport: ccip.CommitStoreReport{
					TokenPrices: []ccip.TokenPrice{
						{
							Token: ccip.Address("seconde token address1"),
							Value: big.NewInt(7),
						},
						{
							Token: ccip.Address("second token address2"),
							Value: big.NewInt(8),
						},
					},
					GasPrices: []ccip.GasPrice{
						{
							DestChainSelector: 17,
							Value:             big.NewInt(17),
						},
						{
							DestChainSelector: 19,
							Value:             big.NewInt(19),
						},
					},
					Interval: ccip.CommitStoreInterval{
						Min: 900,
						Max: 999,
					},
					MerkleRoot: [32]byte{3: 23, 5: 27},
				},
			},
		},

		//get commit report matching seq num test data
		getCommitReportMatchingSeqNumRequest: getCommitReportMatchingSeqNumRequest{
			seqNum:        1,
			confirmations: 1,
		},
		// use the same response as get accepted commit reports gte timestamp

		// get commit store static config test data
		getCommitStoreStaticConfigResponse: ccip.CommitStoreStaticConfig{
			ChainSelector:       1,
			SourceChainSelector: 2,
			OnRamp:              ccip.Address("onramp address"),
			ArmProxy:            ccip.Address("arm proxy address"),
		},

		// get expected next sequence number test data
		getExpectedNextSequenceNumberResponse: 100,

		// get latest price epoch and round test data
		getLatestPriceEpochAndRoundResponse: 1000,

		// is blessed test data
		isBlessedRequest:  [32]byte{0: 1, 31: 7},
		isBlessedResponse: true,

		// is dest chain healthy test data
		isDestChainHealthyResponse: true,

		// is down test data
		isDownResponse: true,

		// offchain config test data
		offchainConfigResponse: ccip.CommitOffchainConfig{
			GasPriceDeviationPPB:   1000,
			GasPriceHeartBeat:      2 * time.Microsecond,
			TokenPriceDeviationPPB: 1000,
			TokenPriceHeartBeat:    3 * time.Millisecond,
			InflightCacheExpiry:    5 * time.Second,
		},

		// verify execution report test data
		verifyExecutionReportRequest: ccip.ExecReport{
			Messages: []ccip.EVM2EVMMessage{
				{
					SequenceNumber:      1,
					GasLimit:            big.NewInt(1),
					Nonce:               1,
					MessageID:           ccip.Hash{1},
					SourceChainSelector: 1,
					Sender:              ccip.Address("sender"),
					Receiver:            ccip.Address("receiver"),
					Strict:              true,
					FeeToken:            ccip.Address("feeToken"),
					FeeTokenAmount:      big.NewInt(1),
					Data:                []byte("data"),
					TokenAmounts: []ccip.TokenAmount{
						{
							Token:  ccip.Address("token"),
							Amount: big.NewInt(1),
						},
						{
							Token:  ccip.Address("token2"),
							Amount: big.NewInt(2),
						},
					},
					SourceTokenData: [][]byte{
						[]byte("sourceTokenData"),
						[]byte("sourceTokenData2"),
					},
				},
				{
					SequenceNumber:      2,
					GasLimit:            big.NewInt(2),
					Nonce:               2,
					MessageID:           ccip.Hash{2},
					SourceChainSelector: 2,
					Sender:              ccip.Address("sender2"),
					Receiver:            ccip.Address("receiver2"),
					Strict:              true,
					FeeToken:            ccip.Address("feeToken2"),
					FeeTokenAmount:      big.NewInt(2),
					Data:                []byte("data2"),
					TokenAmounts: []ccip.TokenAmount{
						{
							Token:  ccip.Address("second token"),
							Amount: big.NewInt(7),
						},
						{
							Token:  ccip.Address("second token2"),
							Amount: big.NewInt(11),
						},
					},
					SourceTokenData: [][]byte{
						[]byte("second sourceTokenData"),
						[]byte("second sourceTokenData2"),
					},
				},
			},
			Proofs: [][32]byte{
				{11},
				{79},
			},
			OffchainTokenData: [][][]byte{
				{
					[]byte("offchainTokenData"),
				},
			},
			ProofFlagBits: big.NewInt(1),
		},
	},
}

type CommitStoreReaderEvaluator interface {
	ccip.CommitStoreReader
	testtypes.Evaluator[ccip.CommitStoreReader]
}

type staticCommitStoreReader struct {
	staticCommitStoreReaderConfig
}

// ChangeConfig implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (ccip.Address, error) {
	if !bytes.Equal(onchainConfig, s.changeConfigRequest.onchainConfig) {
		return ccip.Address(""), fmt.Errorf("change config expected onchain config %v, got %v", s.changeConfigRequest.onchainConfig, onchainConfig)
	}
	if !bytes.Equal(offchainConfig, s.changeConfigRequest.offchainConfig) {
		return ccip.Address(""), fmt.Errorf("change config expected offchain config %v, got %v", s.changeConfigRequest.offchainConfig, offchainConfig)
	}
	return s.changeConfigResponse, nil
}

// Close implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) Close() error {
	return nil
}

// DecodeCommitReport implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) DecodeCommitReport(ctx context.Context, report []byte) (ccip.CommitStoreReport, error) {
	if !bytes.Equal(report, s.decodeCommitReportRequest) {
		return ccip.CommitStoreReport{}, fmt.Errorf("decode commit report expected %v, got %v", s.decodeCommitReportRequest, report)
	}
	return s.decodeCommitReportResponse, nil
}

// EncodeCommitReport implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) EncodeCommitReport(ctx context.Context, report ccip.CommitStoreReport) ([]byte, error) {
	if !reflect.DeepEqual(s.decodeCommitReportResponse, report) {
		return nil, fmt.Errorf("encode commit report expected %v, got %v", s.decodeCommitReportResponse, report)
	}
	return s.decodeCommitReportRequest, nil
}

// Evaluate implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) Evaluate(ctx context.Context, other ccip.CommitStoreReader) error {
	// change config
	gotConfig, err := other.ChangeConfig(ctx, s.changeConfigRequest.onchainConfig, s.changeConfigRequest.offchainConfig)
	if err != nil {
		return fmt.Errorf("failed to call other.ChangeConfig: %w", err)
	}
	if gotConfig != s.changeConfigResponse {
		return fmt.Errorf("change config expected %v, got %v", s.changeConfigResponse, gotConfig)
	}

	// decode commit
	gotCommit, err := other.DecodeCommitReport(ctx, s.decodeCommitReportRequest)
	if err != nil {
		return fmt.Errorf("failed to call other.DecodeCommitReport: %w", err)
	}
	if !reflect.DeepEqual(s.decodeCommitReportResponse, gotCommit) {
		return fmt.Errorf("decode commit expected %v, got %v", s.decodeCommitReportResponse, gotCommit)
	}

	// encode commit
	gotEncodedCommit, err := other.EncodeCommitReport(ctx, s.decodeCommitReportResponse)
	if err != nil {
		return fmt.Errorf("failed to call other.EncodeCommitReport: %w", err)
	}
	if !bytes.Equal(s.decodeCommitReportRequest, gotEncodedCommit) {
		return fmt.Errorf("encode commit expected %v, got %v", s.decodeCommitReportRequest, gotEncodedCommit)
	}

	// gas price estimator
	gotGasPriceEstimator, err := other.GasPriceEstimator(ctx)
	if err != nil {
		return fmt.Errorf("failed to call other.GasPriceEstimator: %w", err)
	}
	err = s.gasPriceEstimatorResponse.Evaluate(ctx, gotGasPriceEstimator)
	if err != nil {
		return fmt.Errorf("failed to evaluate gas price estimator: %w", err)
	}

	// get accepted commit reports gte timestamp
	gotAcceptedCommitReportsGteTimestamp, err := other.GetAcceptedCommitReportsGteTimestamp(ctx, s.getAcceptedCommitReportsGteTimestampRequest.timestamp, s.getAcceptedCommitReportsGteTimestampRequest.confirmations)
	if err != nil {
		return fmt.Errorf("failed to call other.GetAcceptedCommitReportsGteTimestamp: %w", err)
	}
	if !reflect.DeepEqual(s.getAcceptedCommitReportsGteTimestampResponse, gotAcceptedCommitReportsGteTimestamp) {
		return fmt.Errorf("get accepted commit reports gte timestamp expected %v, got %v", s.getAcceptedCommitReportsGteTimestampResponse, gotAcceptedCommitReportsGteTimestamp)
	}

	// get commit report matching seq num
	gotCommitReportMatchingSeqNum, err := other.GetCommitReportMatchingSeqNum(ctx, s.getCommitReportMatchingSeqNumRequest.seqNum, s.getCommitReportMatchingSeqNumRequest.confirmations)
	if err != nil {
		return fmt.Errorf("failed to call other.GetCommitReportMatchingSeqNum: %w", err)
	}
	// for simplicity, just use the same response as get accepted commit reports gte timestamp
	if !reflect.DeepEqual(s.getAcceptedCommitReportsGteTimestampResponse, gotCommitReportMatchingSeqNum) {
		return fmt.Errorf("get commit report matching seq num expected %v, got %v", s.getAcceptedCommitReportsGteTimestampResponse, gotCommitReportMatchingSeqNum)
	}

	// get commit store static config
	gotCommitStoreStaticConfig, err := other.GetCommitStoreStaticConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to call other.GetCommitStoreStaticConfig: %w", err)
	}
	if !reflect.DeepEqual(s.getCommitStoreStaticConfigResponse, gotCommitStoreStaticConfig) {
		return fmt.Errorf("get commit store static config expected %v, got %v", s.getCommitStoreStaticConfigResponse, gotCommitStoreStaticConfig)
	}

	// get expected next sequence number
	gotExpectedNextSequenceNumber, err := other.GetExpectedNextSequenceNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to call other.GetExpectedNextSequenceNumber: %w", err)
	}
	if gotExpectedNextSequenceNumber != s.getExpectedNextSequenceNumberResponse {
		return fmt.Errorf("get expected next sequence number expected %v, got %v", s.getExpectedNextSequenceNumberResponse, gotExpectedNextSequenceNumber)
	}

	// get latest price epoch and round
	gotLatestPriceEpochAndRound, err := other.GetLatestPriceEpochAndRound(ctx)
	if err != nil {
		return fmt.Errorf("failed to call other.GetLatestPriceEpochAndRound: %w", err)
	}
	if gotLatestPriceEpochAndRound != s.getLatestPriceEpochAndRoundResponse {
		return fmt.Errorf("get latest price epoch and round expected %v, got %v", s.getLatestPriceEpochAndRoundResponse, gotLatestPriceEpochAndRound)
	}

	// is blessed
	gotIsBlessed, err := other.IsBlessed(ctx, s.isBlessedRequest)
	if err != nil {
		return fmt.Errorf("failed to call other.IsBlessed: %w", err)
	}
	if gotIsBlessed != s.isBlessedResponse {
		return fmt.Errorf("is blessed expected %v, got %v", s.isBlessedResponse, gotIsBlessed)
	}

	// is down
	gotIsDown, err := other.IsDown(ctx)
	if err != nil {
		return fmt.Errorf("failed to call other.IsDown: %w", err)
	}
	if gotIsDown != s.isDownResponse {
		return fmt.Errorf("is down expected %v, got %v", s.isDownResponse, gotIsDown)
	}

	// offchain config
	gotOffchainConfig, err := other.OffchainConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to call other.OffchainConfig: %w", err)
	}
	if !reflect.DeepEqual(s.offchainConfigResponse, gotOffchainConfig) {
		return fmt.Errorf("offchain config expected %v, got %v", s.offchainConfigResponse, gotOffchainConfig)
	}

	// verify execution report
	gotVerifyExecutionReport, err := other.VerifyExecutionReport(ctx, s.verifyExecutionReportRequest)
	if err != nil {
		return fmt.Errorf("failed to call other.VerifyExecutionReport: %w", err)
	}
	if gotVerifyExecutionReport != s.verifyExecutionReportResponse {
		return fmt.Errorf("verify execution report expected %v, got %v", s.verifyExecutionReportResponse, gotVerifyExecutionReport)
	}

	return nil
}

// GasPriceEstimator implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) GasPriceEstimator(ctx context.Context) (ccip.GasPriceEstimatorCommit, error) {
	return s.gasPriceEstimatorResponse, nil
}

// GetAcceptedCommitReportsGteTimestamp implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]ccip.CommitStoreReportWithTxMeta, error) {
	if ts != s.getAcceptedCommitReportsGteTimestampRequest.timestamp {
		return nil, fmt.Errorf("get accepted commit reports gte timestamp expected %v, got %v", s.getAcceptedCommitReportsGteTimestampRequest.timestamp, ts)
	}
	if confirmations != s.getAcceptedCommitReportsGteTimestampRequest.confirmations {
		return nil, fmt.Errorf("get accepted commit reports gte timestamp expected confirmations %v, got %v", s.getAcceptedCommitReportsGteTimestampRequest.confirmations, confirmations)
	}
	return s.getAcceptedCommitReportsGteTimestampResponse, nil
}

// GetCommitReportMatchingSeqNum implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]ccip.CommitStoreReportWithTxMeta, error) {
	if seqNum != s.getCommitReportMatchingSeqNumRequest.seqNum {
		return nil, fmt.Errorf("get commit report matching seq num expected %v, got %v", s.getCommitReportMatchingSeqNumRequest.seqNum, seqNum)
	}
	if confirmations != s.getCommitReportMatchingSeqNumRequest.confirmations {
		return nil, fmt.Errorf("get commit report matching seq num expected confirmations %v, got %v", s.getCommitReportMatchingSeqNumRequest.confirmations, confirmations)
	}

	return s.getAcceptedCommitReportsGteTimestampResponse, nil
}

// GetCommitStoreStaticConfig implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) GetCommitStoreStaticConfig(ctx context.Context) (ccip.CommitStoreStaticConfig, error) {
	return s.getCommitStoreStaticConfigResponse, nil
}

// GetExpectedNextSequenceNumber implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return s.getExpectedNextSequenceNumberResponse, nil
}

// GetLatestPriceEpochAndRound implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return s.getLatestPriceEpochAndRoundResponse, nil
}

// IsBlessed implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	if root != s.isBlessedRequest {
		return false, fmt.Errorf("is blessed expected %v, got %v", s.isBlessedRequest, root)
	}
	return s.isBlessedResponse, nil
}

// IsDestChainHealthy implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) IsDestChainHealthy(ctx context.Context) (bool, error) {
	return s.isDestChainHealthyResponse, nil
}

// IsDown implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	return s.isDownResponse, nil
}

// OffchainConfig implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) OffchainConfig(ctx context.Context) (ccip.CommitOffchainConfig, error) {
	return s.offchainConfigResponse, nil
}

// VerifyExecutionReport implements CommitStoreReaderEvaluator.
func (s staticCommitStoreReader) VerifyExecutionReport(ctx context.Context, report ccip.ExecReport) (bool, error) {
	if !reflect.DeepEqual(s.verifyExecutionReportRequest, report) {
		return false, fmt.Errorf("verify execution report expected %v, got %v", s.verifyExecutionReportRequest, report)
	}
	return s.verifyExecutionReportResponse, nil
}

var _ CommitStoreReaderEvaluator = staticCommitStoreReader{}

type staticCommitStoreReaderConfig struct {
	changeConfigRequest  changeConfigRequest
	changeConfigResponse ccip.Address

	decodeCommitReportRequest  []byte
	decodeCommitReportResponse ccip.CommitStoreReport

	// rather than explicit encode request/response, we just use decode in reverse

	gasPriceEstimatorResponse GasPriceEstimatorCommitEvaluator

	getAcceptedCommitReportsGteTimestampRequest  getAcceptedCommitReportsGteTimestampRequest
	getAcceptedCommitReportsGteTimestampResponse []ccip.CommitStoreReportWithTxMeta

	getCommitStoreStaticConfigResponse ccip.CommitStoreStaticConfig

	getCommitReportMatchingSeqNumRequest getCommitReportMatchingSeqNumRequest

	getExpectedNextSequenceNumberResponse uint64

	getLatestPriceEpochAndRoundResponse uint64

	isBlessedRequest  [32]byte
	isBlessedResponse bool

	isDestChainHealthyResponse bool

	isDownResponse bool

	offchainConfigResponse ccip.CommitOffchainConfig

	verifyExecutionReportRequest  ccip.ExecReport
	verifyExecutionReportResponse bool
}

type getAcceptedCommitReportsGteTimestampRequest struct {
	timestamp     time.Time
	confirmations int
}

type getCommitReportMatchingSeqNumRequest struct {
	seqNum        uint64
	confirmations int
}
