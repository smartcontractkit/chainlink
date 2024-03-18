package test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var GasPriceEstimatorExec = staticGasPriceEstimatorExec{
	commonStaticGasPriceEstimator: commonGasPriceEstimator,
	staticGasPriceEstimatorExecConfig: staticGasPriceEstimatorExecConfig{
		estimateMsgCostUSDRequest: estimateMsgCostUSDRequest{
			p:                  big.NewInt(1),
			wrappedNativePrice: big.NewInt(2),
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					SequenceNumber: 1,
					GasLimit:       big.NewInt(3),
					Data:           []byte{4},
					TokenAmounts: []cciptypes.TokenAmount{
						{
							Token:  cciptypes.Address("token1"),
							Amount: big.NewInt(5),
						},
					},
					SourceTokenData: [][]byte{
						{6},
					},
				},
			},
		},
		estimateMsgCostUSDResponse: big.NewInt(7),
	},
}

type GasPriceEstimatorExecEvaluator interface {
	cciptypes.GasPriceEstimatorExec
	testtypes.Evaluator[cciptypes.GasPriceEstimatorExec]
}

type staticGasPriceEstimatorExecConfig struct {
	estimateMsgCostUSDRequest  estimateMsgCostUSDRequest
	estimateMsgCostUSDResponse *big.Int
}

type staticGasPriceEstimatorExec struct {
	commonStaticGasPriceEstimator
	staticGasPriceEstimatorExecConfig
}

// EstimateMsgCostUSD implements GasPriceEstimatorExecEvaluator.
func (s staticGasPriceEstimatorExec) EstimateMsgCostUSD(p *big.Int, wrappedNativePrice *big.Int, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error) {
	if s.estimateMsgCostUSDRequest.p.Cmp(p) != 0 {
		return nil, fmt.Errorf("expected p %v, got %v", s.estimateMsgCostUSDRequest.p, p)
	}
	if s.estimateMsgCostUSDRequest.wrappedNativePrice.Cmp(wrappedNativePrice) != 0 {
		return nil, fmt.Errorf("expected wrappedNativePrice %v, got %v", s.estimateMsgCostUSDRequest.wrappedNativePrice, wrappedNativePrice)
	}
	if !reflect.DeepEqual(s.estimateMsgCostUSDRequest.msg, msg) {
		return nil, fmt.Errorf("expected msg %v, got %v", s.estimateMsgCostUSDRequest.msg, msg)
	}
	return s.estimateMsgCostUSDResponse, nil
}

var _ GasPriceEstimatorExecEvaluator = staticGasPriceEstimatorExec{}

// DenoteInUSD implements GasPriceEstimatorExecEvaluator.
func (s staticGasPriceEstimatorExec) DenoteInUSD(p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	if s.denoteInUSDRequest.p.Cmp(p) != 0 {
		return nil, fmt.Errorf("expected p %v, got %v", s.denoteInUSDRequest.p, p)
	}
	if s.denoteInUSDRequest.wrappedNativePrice.Cmp(wrappedNativePrice) != 0 {
		return nil, fmt.Errorf("expected wrappedNativePrice %v, got %v", s.denoteInUSDRequest.wrappedNativePrice, wrappedNativePrice)
	}
	return s.denoteInUSDResponse.result, nil
}

// Evaluate implements GasPriceEstimatorExecEvaluator.
func (s staticGasPriceEstimatorExec) Evaluate(ctx context.Context, other cciptypes.GasPriceEstimatorExec) error {
	// GetGasPrice test case
	gotGas, err := other.GetGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to other.GetGasPrice: %w", err)
	}
	if s.getGasPriceResponse.Cmp(gotGas) != 0 {
		return fmt.Errorf("expected other.GetGasPrice %v, got %v", s.getGasPriceResponse, gotGas)
	}

	// Median test case
	gotMedian, err := other.Median(s.medianRequest.gasPrices)
	if err != nil {
		return fmt.Errorf("failed to other.Median: %w", err)
	}
	if s.medianResponse.Cmp(gotMedian) != 0 {
		return fmt.Errorf("expected other.Median %v, got %v", s.medianResponse, gotMedian)
	}

	// EstimateMsgCostUSD test case
	gotEstimate, err := other.EstimateMsgCostUSD(s.estimateMsgCostUSDRequest.p, s.estimateMsgCostUSDRequest.wrappedNativePrice, s.estimateMsgCostUSDRequest.msg)
	if err != nil {
		return fmt.Errorf("failed to other.EstimateMsgCostUSD: %w", err)
	}
	if s.estimateMsgCostUSDResponse.Cmp(gotEstimate) != 0 {
		return fmt.Errorf("expected other.EstimateMsgCostUSD %v, got %v", s.estimateMsgCostUSDResponse, gotEstimate)
	}

	gotDenoteInUSD, err := other.DenoteInUSD(s.denoteInUSDRequest.p, s.denoteInUSDRequest.wrappedNativePrice)
	if err != nil {
		return fmt.Errorf("failed to other.DenoteInUSD: %w", err)
	}
	if s.denoteInUSDResponse.result.Cmp(gotDenoteInUSD) != 0 {
		return fmt.Errorf("expected other.DenoteInUSD %v, got %v", s.denoteInUSDResponse.result, gotDenoteInUSD)
	}

	return nil
}

// GetGasPrice implements GasPriceEstimatorExecEvaluator.
func (s staticGasPriceEstimatorExec) GetGasPrice(ctx context.Context) (*big.Int, error) {
	return s.getGasPriceResponse, nil
}

// Median implements GasPriceEstimatorExecEvaluator.
func (s staticGasPriceEstimatorExec) Median(gasPrices []*big.Int) (*big.Int, error) {
	if len(gasPrices) != len(s.medianRequest.gasPrices) {
		return nil, fmt.Errorf("expected gas prices len %d, got %d", len(s.medianRequest.gasPrices), len(gasPrices))
	}
	for i, p := range gasPrices {
		if s.medianRequest.gasPrices[i].Cmp(p) != 0 {
			return nil, fmt.Errorf("expected gas price %d %v, got %v", i, s.medianRequest.gasPrices[i], p)
		}
	}
	return s.medianResponse, nil
}

type estimateMsgCostUSDRequest struct {
	p                  *big.Int
	wrappedNativePrice *big.Int
	msg                cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
}
