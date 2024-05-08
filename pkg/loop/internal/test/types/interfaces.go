package types_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	looptypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

// Evaluator is a helper interface for testing types
// implementations must run all the methods of the other type and
// 1. return on the first error encounter 2. return an error the first time
// that the result of an executed method does not match the expected result
type Evaluator[T any] interface {
	Evaluate(ctx context.Context, other T) error
}

// AssertEqualer is a helper interface for testing types
// implementations must run all the methods asserting for equality with the other type
// and no error. Implementations should use parallel subcomponent testing eg
// t.Parallel() and t.Run("subcomponentX", func(t *testing.T) { ... })
type AssertEqualer[T any] interface {
	AssertEqual(ctx context.Context, t *testing.T, other T)
}

// ChainReaderEvaluator is a helper interface for testing ChainReaders
type ChainReaderEvaluator interface {
	types.ChainReader
	Evaluator[types.ChainReader]
}

// ContractTransmitterEvaluator is a helper interface for testing ContractTransmitters
type ContractTransmitterEvaluator interface {
	libocr.ContractTransmitter
	Evaluator[libocr.ContractTransmitter]
}

// OCR3ContractTransmitterEvaluator is a helper interface for testing OCR3 ContractTransmitters
type OCR3ContractTransmitterEvaluator interface {
	ocr3types.ContractTransmitter[[]byte]
	Evaluator[ocr3types.ContractTransmitter[[]byte]]
}

// ContractConfigTrackerEvaluator is a helper interface for testing ContractConfigTrackers
type ContractConfigTrackerEvaluator interface {
	libocr.ContractConfigTracker
	Evaluator[libocr.ContractConfigTracker]
}

// OffchainConfigDigesterEvaluator is a helper interface for testing OffchainConfigDigesters
type OffchainConfigDigesterEvaluator interface {
	libocr.OffchainConfigDigester
	Evaluator[libocr.OffchainConfigDigester]
}

// TelemetryEvaluator is a helper interface for testing TelemetryClients
type TelemetryEvaluator interface {
	core.TelemetryClient
	Evaluator[core.TelemetryClient]
}

// PipelineEvaluator is a helper interface for testing PipelineRunnerServices
type PipelineEvaluator interface {
	core.PipelineRunnerService
	Evaluator[core.PipelineRunnerService]
}

// CodecEvaluator is a helper interface for testing Codecs
type CodecEvaluator interface {
	types.Codec
	Evaluator[types.Codec]
}

// ErrorLogEvaluator is a helper interface for testing ErrorLogs
type ErrorLogEvaluator interface {
	core.ErrorLog
	Evaluator[core.ErrorLog]
}

// ValidationEvaluator is a helper interface for testing ValidationService
type ValidationEvaluator interface {
	core.ValidationService
	Evaluator[core.ValidationService]
}

type MedianProviderTester interface {
	types.MedianProvider
	Evaluator[types.MedianProvider]
	AssertEqualer[types.MedianProvider]
}

type RelayerTester interface {
	looptypes.PluginRelayer
	looptypes.Relayer
	// implements all the possible providers as a one-stop shop for testing
	looptypes.MercuryProvider
	looptypes.MedianProvider
	looptypes.CCIPExecProvider
	looptypes.CCIPCommitProvider
	looptypes.OCR3CapabilityProvider

	AssertEqualer[looptypes.Relayer]
}

type ReportingPluginTester interface {
	libocr.ReportingPlugin
	AssertEqualer[libocr.ReportingPlugin]
}

type PluginProviderTester interface {
	types.PluginProvider
	AssertEqualer[types.PluginProvider]
	Evaluator[types.PluginProvider]
}

type OCR3CapabilityProviderTester interface {
	types.OCR3CapabilityProvider
	AssertEqualer[types.OCR3CapabilityProvider]
	Evaluator[types.OCR3CapabilityProvider]
}
