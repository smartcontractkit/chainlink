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

type CommitProviderEvaluator interface {
	types.CCIPCommitProvider
	testtypes.Evaluator[types.CCIPCommitProvider]
}

type CommitProviderTester interface {
	types.CCIPCommitProvider
	testtypes.Evaluator[types.CCIPCommitProvider]
	testtypes.AssertEqualer[types.CCIPCommitProvider]
}

// CommitProvider is a static implementation of the CommitProviderTester interface.
// It is to be used in tests the verify grpc implementations of the CommitProvider interface.
var CommitProvider = staticCommitProvider{
	staticCommitProviderConfig: staticCommitProviderConfig{
		addr:                      ccip.Address("some address"),
		offchainDigester:          ocr2test.OffchainConfigDigester,
		contractTracker:           ocr2test.ContractConfigTracker,
		contractTransmitter:       ocr2test.ContractTransmitter,
		commitStoreReader:         CommitStoreReader,
		offRampReader:             OffRampReader,
		onRampReader:              OnRampReader,
		priceGetter:               PriceGetter,
		priceRegistryReader:       PriceRegistryReader,
		sourceNativeTokenResponse: ccip.Address("source native token response"),
	},
}

var _ CommitProviderTester = staticCommitProvider{}

type staticCommitProviderConfig struct {
	addr                ccip.Address
	offchainDigester    testtypes.OffchainConfigDigesterEvaluator
	contractTracker     testtypes.ContractConfigTrackerEvaluator
	contractTransmitter testtypes.ContractTransmitterEvaluator

	commitStoreReader         CommitStoreReaderEvaluator
	offRampReader             OffRampEvaluator
	onRampReader              OnRampEvaluator
	priceGetter               PriceGetterEvaluator
	priceRegistryReader       PriceRegistryReaderEvaluator
	sourceNativeTokenResponse ccip.Address
}

type staticCommitProvider struct {
	staticCommitProviderConfig
}

// ChainReader implements CommitProviderEvaluator.
func (s staticCommitProvider) ChainReader() types.ContractReader {
	return nil
}

// Close implements CommitProviderEvaluator.
func (s staticCommitProvider) Close() error {
	return nil
}

// Codec implements CommitProviderEvaluator.
func (s staticCommitProvider) Codec() types.Codec {
	return nil
}

// ContractConfigTracker implements CommitProviderEvaluator.
func (s staticCommitProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return s.contractTracker
}

// ContractTransmitter implements CommitProviderEvaluator.
func (s staticCommitProvider) ContractTransmitter() libocr.ContractTransmitter {
	return s.contractTransmitter
}

// Evaluate implements CommitProviderEvaluator.
func (s staticCommitProvider) Evaluate(ctx context.Context, other types.CCIPCommitProvider) error {
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

	// PriceGetter test case
	otherPriceGetter, err := other.NewPriceGetter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create other price getter: %w", err)
	}
	err = s.priceGetter.Evaluate(ctx, otherPriceGetter)
	if err != nil {
		return evaluationError{err: err, component: priceGetterComponent}
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

// HealthReport implements CommitProviderEvaluator.
func (s staticCommitProvider) HealthReport() map[string]error {
	panic("unimplemented")
}

// Name implements CommitProviderEvaluator.
func (s staticCommitProvider) Name() string {
	panic("unimplemented")
}

// NewCommitStoreReader implements CommitProviderEvaluator.
func (s staticCommitProvider) NewCommitStoreReader(ctx context.Context, addr ccip.Address) (ccip.CommitStoreReader, error) {
	return s.commitStoreReader, nil
}

// NewOffRampReader implements CommitProviderEvaluator.
func (s staticCommitProvider) NewOffRampReader(ctx context.Context, addr ccip.Address) (ccip.OffRampReader, error) {
	return s.offRampReader, nil
}

// NewOnRampReader implements CommitProviderEvaluator.
func (s staticCommitProvider) NewOnRampReader(ctx context.Context, addr ccip.Address) (ccip.OnRampReader, error) {
	return s.onRampReader, nil
}

// NewPriceGetter implements CommitProviderEvaluator.
func (s staticCommitProvider) NewPriceGetter(ctx context.Context) (ccip.PriceGetter, error) {
	return s.priceGetter, nil
}

// NewPriceRegistryReader implements CommitProviderEvaluator.
func (s staticCommitProvider) NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error) {
	return s.priceRegistryReader, nil
}

// OffchainConfigDigester implements CommitProviderEvaluator.
func (s staticCommitProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return s.offchainDigester
}

// Ready implements CommitProviderEvaluator.
func (s staticCommitProvider) Ready() error {
	return nil
}

// SourceNativeToken implements CommitProviderEvaluator.
func (s staticCommitProvider) SourceNativeToken(ctx context.Context) (ccip.Address, error) {
	return s.sourceNativeTokenResponse, nil
}

// Start implements CommitProviderEvaluator.
func (s staticCommitProvider) Start(context.Context) error {
	return nil
}

// AssertEqual implements CommitProviderTester.
func (s staticCommitProvider) AssertEqual(ctx context.Context, t *testing.T, other types.CCIPCommitProvider) {
	t.Run("StaticCommitProvider", func(t *testing.T) {
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
	})
}
