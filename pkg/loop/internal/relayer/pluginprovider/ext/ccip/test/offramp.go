package test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"slices"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var OffRampReader = staticOffRamp{
	staticOffRampConfig: staticOffRampConfig{
		// Address test data
		addressResponse: ccip.Address("addressResponse"),
		// ChangeConfig test data
		changeConfigRequest: changeConfigRequest{
			onchainConfig:  []byte("onchainConfig"),
			offchainConfig: []byte("offchainConfig"),
		},
		changeConfigResponse: changeConfigResponse{
			onchainConfigDigest:  ccip.Address("onchainConfigDigest"),
			offchainConfigDigest: ccip.Address("offchainConfigDigest"),
		},
		// CurrentRateLimiterState test data
		currentRateLimiterStateResponse: ccip.TokenBucketRateLimit{
			Tokens:      big.NewInt(1),
			IsEnabled:   true,
			LastUpdated: 7,
			Capacity:    big.NewInt(2),
			Rate:        big.NewInt(3),
		},
		// DecodeExecutionReport test data
		decodeExecutionReportResponse: ccip.ExecReport{
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
					},
					SourceTokenData: [][]byte{
						[]byte("sourceTokenData"),
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

		// EncodeExecutionReport test data
		encodeExecutionReportRequest: ccip.ExecReport{
			Messages: []ccip.EVM2EVMMessage{
				{
					SequenceNumber: 3,
				},
			},
			Proofs: [][32]byte{
				{3},
			},
		},
		encodeExecutionReportResponse: []byte("encodeExecutionReportResponse"),

		// GasPriceEstimator test data
		gasPriceEstimatorResponse: GasPriceEstimatorExec,

		// GetExecutionState test data
		getExecutionStateRequest:  4,
		getExecutionStateResponse: 5,

		// GetExecutionStateChangesBetweenSeqNums test data
		getExecutionStateChangesBetweenSeqNumsRequest: getExecutionStateChangesBetweenSeqNumsRequest{
			seqNumMin:     6,
			seqNumMax:     7,
			confirmations: 8,
		},
		getExecutionStateChangesBetweenSeqNumsResponse: getExecutionStateChangesBetweenSeqNumsResponse{
			executionStateChangedWithTxMeta: []ccip.ExecutionStateChangedWithTxMeta{
				{
					TxMeta: ccip.TxMeta{
						BlockTimestampUnixMilli: 1,
						BlockNumber:             2,
						TxHash:                  "txHash",
						LogIndex:                3,
					},
					ExecutionStateChanged: ccip.ExecutionStateChanged{
						SequenceNumber: 9,
						Finalized:      true,
					},
				},
			},
		},

		// GetSenderNonce test data
		getSenderNonceRequest:  ccip.Address("getSenderNonceRequest"),
		getSenderNonceResponse: 10,

		// ListSenderNonces test data
		listSenderNoncesRequest:  []ccip.Address{ccip.Address("listSenderNoncesRequest")},
		listSenderNoncesResponse: map[ccip.Address]uint64{ccip.Address("listSenderNoncesRequest"): 10},

		// GetSourceToDestTokensMapping test data
		getSourceToDestTokensMappingResponse: map[ccip.Address]ccip.Address{
			ccip.Address("source"): ccip.Address("dest"),
		},

		// GetStaticConfig test data
		getStaticConfigResponse: ccip.OffRampStaticConfig{
			CommitStore:         ccip.Address("commitStore"),
			ChainSelector:       1,
			SourceChainSelector: 2,
			OnRamp:              ccip.Address("onRamp"),
			PrevOffRamp:         ccip.Address("prevOffRamp"),
			ArmProxy:            ccip.Address("armProxy"),
		},

		// GetTokens test data
		getTokensResponse: ccip.OffRampTokens{
			DestinationTokens: []ccip.Address{
				ccip.Address("destinationToken1"),
				ccip.Address("destinationToken2"),
			},
			SourceTokens: []ccip.Address{
				ccip.Address("sourceToken1"),
				ccip.Address("sourceToken2"),
			},
			DestinationPool: map[ccip.Address]ccip.Address{
				ccip.Address("key1"): ccip.Address("value1"),
			},
		},

		// GetRouter test data
		getRouterResponse: ccip.Address("getRouterResponse"),
	},
}

type OffRampEvaluator interface {
	ccip.OffRampReader
	testtypes.Evaluator[ccip.OffRampReader]
}

var _ OffRampEvaluator = staticOffRamp{}

type staticOffRampConfig struct {
	addressResponse ccip.Address
	changeConfigRequest
	changeConfigResponse

	currentRateLimiterStateResponse ccip.TokenBucketRateLimit
	// DecodeExecutionReport test data
	decodeExecutionReportRequest  []byte
	decodeExecutionReportResponse ccip.ExecReport

	encodeExecutionReportRequest  ccip.ExecReport
	encodeExecutionReportResponse []byte

	gasPriceEstimatorResponse ccip.GasPriceEstimatorExec

	getExecutionStateRequest  uint64
	getExecutionStateResponse uint8

	getExecutionStateChangesBetweenSeqNumsRequest  getExecutionStateChangesBetweenSeqNumsRequest
	getExecutionStateChangesBetweenSeqNumsResponse getExecutionStateChangesBetweenSeqNumsResponse

	getSenderNonceRequest  ccip.Address
	getSenderNonceResponse uint64

	listSenderNoncesRequest  []ccip.Address
	listSenderNoncesResponse map[ccip.Address]uint64

	getSourceToDestTokensMappingResponse map[ccip.Address]ccip.Address

	getStaticConfigResponse ccip.OffRampStaticConfig

	getTokensResponse ccip.OffRampTokens

	getRouterResponse ccip.Address

	offchainConfigResponse ccip.ExecOffchainConfig

	onchainConfigResponse ccip.ExecOnchainConfig
}

type staticOffRamp struct {
	staticOffRampConfig
}

// Address implements OffRampEvaluator.
func (s staticOffRamp) Address(ctx context.Context) (ccip.Address, error) {
	return s.addressResponse, nil
}

// ChangeConfig implements OffRampEvaluator.
func (s staticOffRamp) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (ccip.Address, ccip.Address, error) {
	if !reflect.DeepEqual(onchainConfig, s.onchainConfig) {
		return ccip.Address(""), ccip.Address(""), fmt.Errorf("expected onchainConfig %v but got %v", s.onchainConfig, onchainConfig)
	}
	if !reflect.DeepEqual(offchainConfig, s.offchainConfig) {
		return ccip.Address(""), ccip.Address(""), fmt.Errorf("expected offchainConfig %v but got %v", s.offchainConfig, offchainConfig)
	}
	return s.onchainConfigDigest, s.offchainConfigDigest, nil
}

func (s staticOffRamp) Close() error {
	return nil
}

// CurrentRateLimiterState implements OffRampEvaluator.
func (s staticOffRamp) CurrentRateLimiterState(ctx context.Context) (ccip.TokenBucketRateLimit, error) {
	return s.currentRateLimiterStateResponse, nil
}

// DecodeExecutionReport implements OffRampEvaluator.
func (s staticOffRamp) DecodeExecutionReport(ctx context.Context, report []byte) (ccip.ExecReport, error) {
	if !reflect.DeepEqual(report, s.decodeExecutionReportRequest) {
		return ccip.ExecReport{}, fmt.Errorf("expected report %v but got %v", s.decodeExecutionReportRequest, report)
	}
	return s.decodeExecutionReportResponse, nil
}

// EncodeExecutionReport implements OffRampEvaluator.
func (s staticOffRamp) EncodeExecutionReport(ctx context.Context, report ccip.ExecReport) ([]byte, error) {
	// struggling to get full report equality via  reflect.DeepEqual or assert.ObjectsAreEqual
	// take a short cut and compare the fields we care about
	if len(report.Messages) != len(s.encodeExecutionReportRequest.Messages) {
		return nil, fmt.Errorf(" encodeExecutionReport message len %v but got %v", len(s.encodeExecutionReportRequest.Messages), len(report.Messages))
	}
	for i, message := range report.Messages {
		if message.SequenceNumber != s.encodeExecutionReportRequest.Messages[i].SequenceNumber {
			return nil, fmt.Errorf("expected sequenceNumber %d but got %d", s.encodeExecutionReportRequest.Messages[i].SequenceNumber, message.SequenceNumber)
		}
	}
	return s.encodeExecutionReportResponse, nil
}

// GasPriceEstimator implements OffRampEvaluator.
func (s staticOffRamp) GasPriceEstimator(ctx context.Context) (ccip.GasPriceEstimatorExec, error) {
	return s.gasPriceEstimatorResponse, nil
}

// GetExecutionState implements OffRampEvaluator.
func (s staticOffRamp) GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error) {
	if sequenceNumber != s.getExecutionStateRequest {
		return 0, fmt.Errorf("expected sequenceNumber %d but got %d", s.getExecutionStateRequest, sequenceNumber)
	}
	return s.getExecutionStateResponse, nil
}

// GetExecutionStateChangesBetweenSeqNums implements OffRampEvaluator.
func (s staticOffRamp) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64, confirmations int) ([]ccip.ExecutionStateChangedWithTxMeta, error) {
	if seqNumMin != s.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMin {
		return nil, fmt.Errorf("expected seqNumMin %d but got %d", s.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMin, seqNumMin)
	}
	if seqNumMax != s.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMax {
		return nil, fmt.Errorf("expected seqNumMax %d but got %d", s.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMax, seqNumMax)
	}
	if confirmations != s.getExecutionStateChangesBetweenSeqNumsRequest.confirmations {
		return nil, fmt.Errorf("expected confirmations %d but got %d", s.getExecutionStateChangesBetweenSeqNumsRequest.confirmations, confirmations)
	}
	return s.getExecutionStateChangesBetweenSeqNumsResponse.executionStateChangedWithTxMeta, nil
}

// ListSenderNonces implements OffRampEvaluator.
func (s staticOffRamp) ListSenderNonces(ctx context.Context, senders []ccip.Address) (map[ccip.Address]uint64, error) {
	if len(senders) == 0 || !slices.Equal(senders, s.listSenderNoncesRequest) {
		return nil, fmt.Errorf("expected sender %s but got %s", s.listSenderNoncesRequest, senders)
	}
	return s.listSenderNoncesResponse, nil
}

// GetSourceToDestTokensMapping implements OffRampEvaluator.
func (s staticOffRamp) GetSourceToDestTokensMapping(ctx context.Context) (map[ccip.Address]ccip.Address, error) {
	return s.getSourceToDestTokensMappingResponse, nil
}

// GetStaticConfig implements OffRampEvaluator.
func (s staticOffRamp) GetStaticConfig(ctx context.Context) (ccip.OffRampStaticConfig, error) {
	return s.getStaticConfigResponse, nil
}

// GetTokens implements OffRampEvaluator.
func (s staticOffRamp) GetTokens(ctx context.Context) (ccip.OffRampTokens, error) {
	return s.getTokensResponse, nil
}

// GetRouter implements OffRampEvaluator.
func (s staticOffRamp) GetRouter(ctx context.Context) (ccip.Address, error) {
	return s.getRouterResponse, nil
}

// OffchainConfig implements OffRampEvaluator.
func (s staticOffRamp) OffchainConfig(ctx context.Context) (ccip.ExecOffchainConfig, error) {
	return s.offchainConfigResponse, nil
}

// OnchainConfig implements OffRampEvaluator.
func (s staticOffRamp) OnchainConfig(ctx context.Context) (ccip.ExecOnchainConfig, error) {
	return s.onchainConfigResponse, nil
}

// Evaluate implements OffRampEvaluator.
func (s staticOffRamp) Evaluate(ctx context.Context, other ccip.OffRampReader) error {
	// Address test case
	address, err := other.Address(ctx)
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}
	if address != s.addressResponse {
		return fmt.Errorf("expected address %s but got %s", s.addressResponse, address)
	}

	// ChangeConfig test case
	gotState, err := other.CurrentRateLimiterState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get currentRateLimiterState: %w", err)
	}
	if !assert.ObjectsAreEqual(gotState, s.currentRateLimiterStateResponse) {
		return fmt.Errorf("expected currentRateLimiterState %v but got %v", s.currentRateLimiterStateResponse, gotState)
	}

	// DecodeExecutionReport test case
	gotReport, err := other.DecodeExecutionReport(ctx, s.decodeExecutionReportRequest)
	if err != nil {
		return fmt.Errorf("failed to decodeExecutionReport: %w", err)
	}
	if !reflect.DeepEqual(gotReport, s.decodeExecutionReportResponse) {
		return fmt.Errorf("expected decodeExecutionReport %v but got %v", s.decodeExecutionReportResponse, gotReport)
	}

	// EncodeExecutionReport test case
	encodeExecutionReport, err := other.EncodeExecutionReport(ctx, s.encodeExecutionReportRequest)
	if err != nil {
		return fmt.Errorf("failed to encodeExecutionReport: %w", err)
	}
	if !reflect.DeepEqual(encodeExecutionReport, s.encodeExecutionReportResponse) {
		return fmt.Errorf("expected encodeExecutionReport %v but got %v", s.encodeExecutionReportResponse, encodeExecutionReport)
	}

	gasPriceEstimator, err := other.GasPriceEstimator(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gasPriceEstimator: %w", err)
	}
	// exercise all the gas price estimator methods
	// GetGasPrice test case
	price, err := gasPriceEstimator.GetGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get other gasPrice: %w", err)
	}
	expectedGas, err := GasPriceEstimatorExec.GetGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expected gasPrice: %w", err)
	}
	if price.Cmp(expectedGas) != 0 {
		return fmt.Errorf("expected gasPrice %v but got %v", GasPriceEstimatorExec.getGasPriceResponse, price)
	}
	// DenoteInUSD test case
	gotusd, err := gasPriceEstimator.DenoteInUSD(GasPriceEstimatorExec.denoteInUSDRequest.p, GasPriceEstimatorExec.denoteInUSDRequest.wrappedNativePrice)
	if err != nil {
		return fmt.Errorf("failed to get other usd: %w", err)
	}
	if gotusd.Cmp(GasPriceEstimatorExec.denoteInUSDResponse.result) != 0 {
		return fmt.Errorf("expected usd %v but got %v", GasPriceEstimatorExec.denoteInUSDResponse.result, gotusd)
	}
	// EstimateMsgCostUSD test case
	cost, err := gasPriceEstimator.EstimateMsgCostUSD(
		GasPriceEstimatorExec.estimateMsgCostUSDRequest.p,
		GasPriceEstimatorExec.estimateMsgCostUSDRequest.wrappedNativePrice,
		GasPriceEstimatorExec.estimateMsgCostUSDRequest.msg,
	)
	if err != nil {
		return fmt.Errorf("failed to get other cost: %w", err)
	}
	if cost.Cmp(GasPriceEstimatorExec.estimateMsgCostUSDResponse) != 0 {
		return fmt.Errorf("expected cost %v but got %v", GasPriceEstimatorExec.estimateMsgCostUSDResponse, cost)
	}
	// Median test case
	median, err := gasPriceEstimator.Median(GasPriceEstimatorExec.medianRequest.gasPrices)
	if err != nil {
		return fmt.Errorf("failed to get other median: %w", err)
	}
	if median.Cmp(GasPriceEstimatorExec.medianResponse) != 0 {
		return fmt.Errorf("expected median %v but got %v", GasPriceEstimatorExec.medianResponse, median)
	}

	getExecutionState, err := other.GetExecutionState(ctx, s.getExecutionStateRequest)
	if err != nil {
		return fmt.Errorf("failed to get getExecutionState: %w", err)
	}
	if getExecutionState != s.getExecutionStateResponse {
		return fmt.Errorf("expected getExecutionState %d but got %d", s.getExecutionStateResponse, getExecutionState)
	}

	getExecutionStateChangesBetweenSeqNums, err := other.GetExecutionStateChangesBetweenSeqNums(ctx,
		s.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMin,
		s.getExecutionStateChangesBetweenSeqNumsRequest.seqNumMax,
		s.getExecutionStateChangesBetweenSeqNumsRequest.confirmations)
	if err != nil {
		return fmt.Errorf("failed to get getExecutionStateChangesBetweenSeqNums: %w", err)
	}
	if !reflect.DeepEqual(getExecutionStateChangesBetweenSeqNums, s.getExecutionStateChangesBetweenSeqNumsResponse.executionStateChangedWithTxMeta) {
		return fmt.Errorf("expected getExecutionStateChangesBetweenSeqNums %v but got %v", s.getExecutionStateChangesBetweenSeqNumsResponse.executionStateChangedWithTxMeta, getExecutionStateChangesBetweenSeqNums)
	}

	getSourceToDestTokensMapping, err := other.GetSourceToDestTokensMapping(ctx)
	if err != nil {
		return fmt.Errorf("failed to get getSourceToDestTokensMapping: %w", err)
	}
	if !reflect.DeepEqual(getSourceToDestTokensMapping, s.getSourceToDestTokensMappingResponse) {
		return fmt.Errorf("expected getSourceToDestTokensMapping %v but got %v", s.getSourceToDestTokensMappingResponse, getSourceToDestTokensMapping)
	}

	getStaticConfig, err := other.GetStaticConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get getStaticConfig: %w", err)
	}
	if getStaticConfig != s.getStaticConfigResponse {
		return fmt.Errorf("expected getStaticConfig %v but got %v", s.getStaticConfigResponse, getStaticConfig)
	}

	getTokens, err := other.GetTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed to get getTokens: %w", err)
	}
	if !assert.ObjectsAreEqual(getTokens, s.getTokensResponse) {
		return fmt.Errorf("expected getTokens %v but got %v", s.getTokensResponse, getTokens)
	}

	getRouter, err := other.GetRouter(ctx)
	if err != nil {
		return fmt.Errorf("failed to get getRouter: %w", err)
	}
	if getRouter != s.getRouterResponse {
		return fmt.Errorf("expected getRouter %s but got %s", s.getRouterResponse, getRouter)
	}

	offchainConfig, err := other.OffchainConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get offchainConfig: %w", err)
	}
	if offchainConfig != s.offchainConfigResponse {
		return fmt.Errorf("expected offchainConfig %v but got %v", s.offchainConfigResponse, offchainConfig)
	}

	onchainConfig, err := other.OnchainConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get onchainConfig: %w", err)
	}
	if onchainConfig != s.onchainConfigResponse {
		return fmt.Errorf("expected onchainConfig %v but got %v", s.onchainConfigResponse, onchainConfig)
	}

	return nil
}

type changeConfigRequest struct {
	onchainConfig  []byte
	offchainConfig []byte
}

type changeConfigResponse struct {
	onchainConfigDigest  ccip.Address
	offchainConfigDigest ccip.Address
}

type getExecutionStateChangesBetweenSeqNumsRequest struct {
	seqNumMin     uint64
	seqNumMax     uint64
	confirmations int
}

type getExecutionStateChangesBetweenSeqNumsResponse struct {
	executionStateChangedWithTxMeta []ccip.ExecutionStateChangedWithTxMeta
}
