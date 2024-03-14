package test

import (
	"context"
	"fmt"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

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

var ExecutionProvider = &staticExecProvider{
	staticExecProviderConfig: staticExecProviderConfig{
		addr:                ccip.Address("some address"),
		offchainDigester:    testpluginprovider.OffchainConfigDigester,
		contractTracker:     testpluginprovider.ContractConfigTracker,
		contractTransmitter: testpluginprovider.ContractTransmitter,
		onrampreader:        OnRamp,
		offrampreader:       OffRamp,
	},
}

var _ ExecProviderTester = (*staticExecProvider)(nil)

type staticExecProviderConfig struct {
	addr                ccip.Address
	offchainDigester    testtypes.OffchainConfigDigesterEvaluator
	contractTracker     testtypes.ContractConfigTrackerEvaluator
	contractTransmitter testtypes.ContractTransmitterEvaluator
	onrampreader        OnRampEvaluator
	offrampreader       OffRampEvaluator
	// TODO BCF-2979 fill in the rest of exec provider components
}

type staticExecProvider struct {
	staticExecProviderConfig
}

// ChainReader implements ExecProviderEvaluator.
func (s *staticExecProvider) ChainReader() types.ChainReader {
	return nil
}

// Close implements ExecProviderEvaluator.
func (s *staticExecProvider) Close() error {
	return nil
}

// Codec implements ExecProviderEvaluator.
func (s *staticExecProvider) Codec() types.Codec {
	return nil
}

// ContractConfigTracker implements ExecProviderEvaluator.
func (s *staticExecProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return s.contractTracker
}

// ContractTransmitter implements ExecProviderEvaluator.
func (s *staticExecProvider) ContractTransmitter() libocr.ContractTransmitter {
	return s.contractTransmitter
}

// Evaluate implements ExecProviderEvaluator.
func (s *staticExecProvider) Evaluate(ctx context.Context, other types.CCIPExecProvider) error {
	otherOnRamp, err := other.NewOnRampReader(ctx, "ignored")
	if err != nil {
		return fmt.Errorf("failed to create other on ramp reader: %w", err)
	}
	err = s.onrampreader.Evaluate(ctx, otherOnRamp)
	if err != nil {
		return fmt.Errorf("on ramp reader evaluation failed: %w", err)
	}
	// TODO BCF-2979 other components of exec provider
	return nil
}

// HealthReport implements ExecProviderEvaluator.
func (s *staticExecProvider) HealthReport() map[string]error {
	panic("unimplemented")
}

// Name implements ExecProviderEvaluator.
func (s *staticExecProvider) Name() string {
	panic("unimplemented")
}

// NewCommitStoreReader implements ExecProviderEvaluator.
func (s *staticExecProvider) NewCommitStoreReader(ctx context.Context, addr ccip.Address) (ccip.CommitStoreReader, error) {
	panic("unimplemented")
}

// NewOffRampReader implements ExecProviderEvaluator.
func (s *staticExecProvider) NewOffRampReader(ctx context.Context, addr ccip.Address) (ccip.OffRampReader, error) {
	return s.offrampreader, nil
}

// NewOnRampReader implements ExecProviderEvaluator.
func (s *staticExecProvider) NewOnRampReader(ctx context.Context, addr ccip.Address) (ccip.OnRampReader, error) {
	return s.onrampreader, nil
}

// NewPriceRegistryReader implements ExecProviderEvaluator.
func (s *staticExecProvider) NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error) {
	panic("unimplemented")
}

// NewTokenDataReader implements ExecProviderEvaluator.
func (s *staticExecProvider) NewTokenDataReader(ctx context.Context, tokenAddress ccip.Address) (ccip.TokenDataReader, error) {
	panic("unimplemented")
}

// NewTokenPoolBatchedReader implements ExecProviderEvaluator.
func (s *staticExecProvider) NewTokenPoolBatchedReader(ctx context.Context) (ccip.TokenPoolBatchedReader, error) {
	panic("unimplemented")
}

// OffchainConfigDigester implements ExecProviderEvaluator.
func (s *staticExecProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return s.offchainDigester
}

// Ready implements ExecProviderEvaluator.
func (s *staticExecProvider) Ready() error {
	return nil
}

// SourceNativeToken implements ExecProviderEvaluator.
func (s *staticExecProvider) SourceNativeToken(ctx context.Context) (ccip.Address, error) {
	panic("unimplemented")
}

// Start implements ExecProviderEvaluator.
func (s *staticExecProvider) Start(context.Context) error {
	return nil
}

// AssertEqual implements ExecProviderTester.
func (s *staticExecProvider) AssertEqual(ctx context.Context, t *testing.T, other types.CCIPExecProvider) {
	// TODO BCF-2979 other components of exec provider
	assert.NoError(t, s.Evaluate(ctx, other))
}
