package test

import (
	"context"
	"fmt"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testpluginprovider "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/ocr2/plugin_provider"
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

var ExecutionConfig = types.CCIPExecFactoryGeneratorConfig{
	OnRampAddress:      ccip.Address("onramp"),
	OffRampAddress:     ccip.Address("offramp"),
	CommitStoreAddress: ccip.Address("commitstore"),
	TokenReaderAddress: ccip.Address("tokenreader"),
}

// ExecutionProvider is a static implementation of the ExecProviderTester interface.
// It is to be used in tests the verify grpc implementations of the ExecProvider interface.
var ExecutionProvider = staticExecProvider{
	staticExecProviderConfig: staticExecProviderConfig{
		addr:                ccip.Address("some address"),
		offchainDigester:    testpluginprovider.OffchainConfigDigester,
		contractTracker:     testpluginprovider.ContractConfigTracker,
		contractTransmitter: testpluginprovider.ContractTransmitter,
		onRampReader:        OnRamp,
		offRampReader:       OffRampReader,
		priceRegistryReader: PriceRegistryReader,
	},
}

var _ ExecProviderTester = staticExecProvider{}

type staticExecProviderConfig struct {
	addr                ccip.Address
	offchainDigester    testtypes.OffchainConfigDigesterEvaluator
	contractTracker     testtypes.ContractConfigTrackerEvaluator
	contractTransmitter testtypes.ContractTransmitterEvaluator
	onRampReader        OnRampEvaluator
	offRampReader       OffRampEvaluator
	priceRegistryReader PriceRegistryReaderEvaluator
	// TODO BCF-2979 fill in the rest of exec provider components
}

type staticExecProvider struct {
	staticExecProviderConfig
}

// ChainReader implements ExecProviderEvaluator.
func (s staticExecProvider) ChainReader() types.ChainReader {
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
	// OnRampReader test case
	otherOnRamp, err := other.NewOnRampReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other on ramp reader: %w", err)
	}
	err = s.onRampReader.Evaluate(ctx, otherOnRamp)
	if err != nil {
		return evaluationError{err: err, component: onRampComponent}
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

	// PriceRegistryReader test case
	otherPriceRegistry, err := other.NewPriceRegistryReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other price registry reader: %w", err)
	}
	err = s.priceRegistryReader.Evaluate(ctx, otherPriceRegistry)
	if err != nil {
		return evaluationError{err: err, component: priceRegistryComponent}
	}

	// TODO BCF-2979 other components of exec provider
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
	panic("unimplemented")
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
	panic("unimplemented")
}

// NewTokenPoolBatchedReader implements ExecProviderEvaluator.
func (s staticExecProvider) NewTokenPoolBatchedReader(ctx context.Context) (ccip.TokenPoolBatchedReader, error) {
	panic("unimplemented")
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
	panic("unimplemented")
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

		// TODO BCF-2979 other components of exec provider
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
)
