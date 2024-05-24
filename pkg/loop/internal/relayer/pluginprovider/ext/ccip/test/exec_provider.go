package test

import (
	"context"
	"fmt"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type ExecProviderEvaluator interface {
	types.CCIPExecProvider
	testtypes.Evaluator[types.CCIPExecProvider]
}

type ExecProviderTester interface {
	types.CCIPExecProvider
	testtypes.Evaluator[types.CCIPExecProvider]
	testtypes.AssertEqualer[types.CCIPExecProvider]
}

// ExecutionProvider is a static implementation of the ExecProviderTester interface.
// It is to be used in tests the verify grpc implementations of the ExecProvider interface.
var ExecutionProvider = staticExecProvider{
	staticExecProviderConfig: staticExecProviderConfig{
		addr:                      ccip.Address("some address"),
		offchainDigester:          ocr2test.OffchainConfigDigester,
		contractTracker:           ocr2test.ContractConfigTracker,
		contractTransmitter:       ocr2test.ContractTransmitter,
		commitStoreReader:         CommitStoreReader,
		offRampReader:             OffRampReader,
		onRampReader:              OnRampReader,
		priceRegistryReader:       PriceRegistryReader,
		sourceNativeTokenResponse: ccip.Address("source native token response"),
		tokenDataReader:           TokenDataReader,
		tokenPoolBatchedReader:    TokenPoolBatchedReader,
	},
}

var _ ExecProviderTester = staticExecProvider{}

type staticExecProviderConfig struct {
	addr                ccip.Address
	offchainDigester    testtypes.OffchainConfigDigesterEvaluator
	contractTracker     testtypes.ContractConfigTrackerEvaluator
	contractTransmitter testtypes.ContractTransmitterEvaluator

	commitStoreReader         CommitStoreReaderEvaluator
	offRampReader             OffRampEvaluator
	onRampReader              OnRampEvaluator
	priceRegistryReader       PriceRegistryReaderEvaluator
	sourceNativeTokenResponse ccip.Address
	tokenDataReader           TokenDataReaderEvaluator
	tokenPoolBatchedReader    TokenPoolBatchedReaderEvaluator
}

type staticExecProvider struct {
	staticExecProviderConfig
}

// ChainReader implements ExecProviderEvaluator.
func (s staticExecProvider) ChainReader() types.ContractReader {
	return nil
}

// Close implements ExecProviderEvaluator.
func (s staticExecProvider) Close() error {
	return nil
}

// Codec implements ExecProviderEvaluator.
func (s staticExecProvider) Codec() types.Codec {
	return nil
}

// ContractConfigTracker implements ExecProviderEvaluator.
func (s staticExecProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return s.contractTracker
}

// ContractTransmitter implements ExecProviderEvaluator.
func (s staticExecProvider) ContractTransmitter() libocr.ContractTransmitter {
	return s.contractTransmitter
}

// Evaluate implements ExecProviderEvaluator.
func (s staticExecProvider) Evaluate(ctx context.Context, other types.CCIPExecProvider) error {
	// CommitStoreReader test case
	otherCommitStore, err := other.NewCommitStoreReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other commit store reader: %w", err)
	}
	err = s.commitStoreReader.Evaluate(ctx, otherCommitStore)
	if err != nil {
		return evaluationError{err: err, component: "CommitStoreReader"}
	}

	// OffRampReader test case
	otherOffRamp, err := other.NewOffRampReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other off ramp reader: %w", err)
	}
	err = s.offRampReader.Evaluate(ctx, otherOffRamp)
	if err != nil {
		return evaluationError{err: err, component: offRampComponent}
	}

	// OnRampReader test case
	otherOnRamp, err := other.NewOnRampReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other on ramp reader: %w", err)
	}
	err = s.onRampReader.Evaluate(ctx, otherOnRamp)
	if err != nil {
		return evaluationError{err: err, component: onRampComponent}
	}

	// PriceRegistryReader test case
	otherPriceRegistry, err := other.NewPriceRegistryReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other price registry reader: %w", err)
	}
	err = s.priceRegistryReader.Evaluate(ctx, otherPriceRegistry)
	if err != nil {
		return evaluationError{err: err, component: priceRegistryComponent}
	}

	// TokenDataReader test case
	otherTokenData, err := other.NewTokenDataReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other token data reader: %w", err)
	}
	err = s.tokenDataReader.Evaluate(ctx, otherTokenData)
	if err != nil {
		return evaluationError{err: err, component: "TokenDataReader"}
	}

	// TokenPoolBatchedReader test case
	otherPool, err := other.NewTokenPoolBatchedReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create other token pool batched reader: %w", err)
	}
	err = s.tokenPoolBatchedReader.Evaluate(ctx, otherPool)
	if err != nil {
		return evaluationError{err: err, component: "TokenPoolBatchedReader"}
	}

	// SourceNativeToken test case
	otherSourceNativeToken, err := other.SourceNativeToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get other source native token: %w", err)
	}
	if otherSourceNativeToken != s.sourceNativeTokenResponse {
		return fmt.Errorf("expected source native token %s but got %s", s.sourceNativeTokenResponse, otherSourceNativeToken)
	}
	return nil
}

// HealthReport implements ExecProviderEvaluator.
func (s staticExecProvider) HealthReport() map[string]error {
	panic("unimplemented")
}

// Name implements ExecProviderEvaluator.
func (s staticExecProvider) Name() string {
	panic("unimplemented")
}

// NewCommitStoreReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewCommitStoreReader(ctx context.Context, addr ccip.Address) (ccip.CommitStoreReader, error) {
	return s.commitStoreReader, nil
}

// NewOffRampReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewOffRampReader(ctx context.Context, addr ccip.Address) (ccip.OffRampReader, error) {
	return s.offRampReader, nil
}

// NewOnRampReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewOnRampReader(ctx context.Context, addr ccip.Address) (ccip.OnRampReader, error) {
	return s.onRampReader, nil
}

// NewPriceRegistryReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error) {
	return s.priceRegistryReader, nil
}

// NewTokenDataReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewTokenDataReader(ctx context.Context, tokenAddress ccip.Address) (ccip.TokenDataReader, error) {
	return s.tokenDataReader, nil
}

// NewTokenPoolBatchedReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewTokenPoolBatchedReader(ctx context.Context) (ccip.TokenPoolBatchedReader, error) {
	return s.tokenPoolBatchedReader, nil
}

// OffchainConfigDigester implements ExecProviderEvaluator.
func (s staticExecProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return s.offchainDigester
}

// Ready implements ExecProviderEvaluator.
func (s staticExecProvider) Ready() error {
	return nil
}

// SourceNativeToken implements ExecProviderEvaluator.
func (s staticExecProvider) SourceNativeToken(ctx context.Context) (ccip.Address, error) {
	return s.sourceNativeTokenResponse, nil
}

// Start implements ExecProviderEvaluator.
func (s staticExecProvider) Start(context.Context) error {
	return nil
}

// AssertEqual implements ExecProviderTester.
func (s staticExecProvider) AssertEqual(ctx context.Context, t *testing.T, other types.CCIPExecProvider) {
	t.Run("StaticExecProvider", func(t *testing.T) {
		// OnRampReader test case
		t.Run(onRampComponent, func(t *testing.T) {
			other, err := other.NewOnRampReader(ctx, "ignored")
			require.NoError(t, err)
			assert.NoError(t, s.onRampReader.Evaluate(ctx, other))
		})

		// OffRampReader test case
		t.Run(offRampComponent, func(t *testing.T) {
			other, err := other.NewOffRampReader(ctx, "ignored")
			require.NoError(t, err)
			assert.NoError(t, s.offRampReader.Evaluate(ctx, other))
		})

		// PriceRegistryReader test case
		t.Run(priceRegistryComponent, func(t *testing.T) {
			other, err := other.NewPriceRegistryReader(ctx, "ignored")
			require.NoError(t, err)
			assert.NoError(t, s.priceRegistryReader.Evaluate(ctx, other))
		})

		// SourceNativeToken test case
		t.Run("SourceNativeToken", func(t *testing.T) {
			other, err := other.SourceNativeToken(ctx)
			require.NoError(t, err)
			assert.Equal(t, s.sourceNativeTokenResponse, other)
		})

		// TokenDataReader test case
		t.Run("TokenDataReader", func(t *testing.T) {
			other, err := other.NewTokenDataReader(ctx, "ignored")
			require.NoError(t, err)
			assert.NoError(t, s.tokenDataReader.Evaluate(ctx, other))
		})
	})
}

type evaluationError struct {
	err       error
	component string
}

func (e evaluationError) Error() string {
	return fmt.Sprintf("error evaluating %s: %s", e.component, e.err)
}

const (
	offRampComponent       = "offRamp"
	onRampComponent        = "onRamp"
	priceRegistryComponent = "priceRegistry"
	priceGetterComponent   = "priceGetter"
)
