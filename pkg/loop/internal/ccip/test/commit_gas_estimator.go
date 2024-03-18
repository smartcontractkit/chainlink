package test

import (
	"context"
	"fmt"
	"math/big"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var GasPriceEstimatorCommit = staticGasPriceEstimatorCommit{
	commonStaticGasPriceEstimator: commonGasPriceEstimator,
	staticGasPriceEstimatorCommitConfig: staticGasPriceEstimatorCommitConfig{
		deviatesRequest: deviatesRequest{
			p1: big.NewInt(1),
			p2: big.NewInt(2),
		},
		deviatesResponse: true,
	},
}

type GasPriceEstimatorCommitEvaluator interface {
	cciptypes.GasPriceEstimatorCommit
	testtypes.Evaluator[cciptypes.GasPriceEstimatorCommit]
}

type staticGasPriceEstimatorCommitConfig struct {
	deviatesRequest
	deviatesResponse bool
}

type staticGasPriceEstimatorCommit struct {
	commonStaticGasPriceEstimator
	staticGasPriceEstimatorCommitConfig
}

var _ GasPriceEstimatorCommitEvaluator = staticGasPriceEstimatorCommit{}

// DenoteInUSD implements GasPriceEstimatorCommitEvaluator.
func (s staticGasPriceEstimatorCommit) DenoteInUSD(p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error) {
	if s.denoteInUSDRequest.p.Cmp(p) != 0 {
		return nil, fmt.Errorf("expected p %v, got %v", s.denoteInUSDRequest.p, p)
	}
	if s.denoteInUSDRequest.wrappedNativePrice.Cmp(wrappedNativePrice) != 0 {
		return nil, fmt.Errorf("expected wrappedNativePrice %v, got %v", s.denoteInUSDRequest.wrappedNativePrice, wrappedNativePrice)
	}
	return s.denoteInUSDResponse.result, nil
}

// Deviates implements GasPriceEstimatorCommitEvaluator.
func (s staticGasPriceEstimatorCommit) Deviates(p1 *big.Int, p2 *big.Int) (bool, error) {
	if s.deviatesRequest.p1.Cmp(p1) != 0 {
		return false, fmt.Errorf("expected p1 %v, got %v", s.deviatesRequest.p1, p1)
	}
	if s.deviatesRequest.p2.Cmp(p2) != 0 {
		return false, fmt.Errorf("expected p2 %v, got %v", s.deviatesRequest.p2, p2)
	}
	return s.deviatesResponse, nil
}

// Evaluate implements GasPriceEstimatorCommitEvaluator.
func (s staticGasPriceEstimatorCommit) Evaluate(ctx context.Context, other cciptypes.GasPriceEstimatorCommit) error {
	gotGas, err := other.GetGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to other.GetGasPrice: %w", err)
	}
	if s.getGasPriceResponse.Cmp(gotGas) != 0 {
		return fmt.Errorf("expected other.GetGasPrice %v, got %v", s.getGasPriceResponse, gotGas)
	}

	gotMedian, err := other.Median(s.medianRequest.gasPrices)
	if err != nil {
		return fmt.Errorf("failed to other.Median: %w", err)
	}
	if s.medianResponse.Cmp(gotMedian) != 0 {
		return fmt.Errorf("expected other.Median %v, got %v", s.medianResponse, gotMedian)
	}

	gotDeviates, err := other.Deviates(s.deviatesRequest.p1, s.deviatesRequest.p2)
	if err != nil {
		return fmt.Errorf("failed to other.Deviates: %w", err)
	}
	if s.deviatesResponse != gotDeviates {
		return fmt.Errorf("expected other.Deviates %v, got %v", s.deviatesResponse, gotDeviates)
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

// GetGasPrice implements GasPriceEstimatorCommitEvaluator.
func (s staticGasPriceEstimatorCommit) GetGasPrice(ctx context.Context) (*big.Int, error) {
	return s.getGasPriceResponse, nil
}

// Median implements GasPriceEstimatorCommitEvaluator.
func (s staticGasPriceEstimatorCommit) Median(gasPrices []*big.Int) (*big.Int, error) {
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
