package types_test

import (
	"context"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
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
	types.TelemetryClient
	Evaluator[types.TelemetryClient]
}

// PipelineEvaluator is a helper interface for testing PipelineRunnerServices
type PipelineEvaluator interface {
	types.PipelineRunnerService
	Evaluator[types.PipelineRunnerService]
}

// CodecEvaluator is a helper interface for testing Codecs
type CodecEvaluator interface {
	types.Codec
	Evaluator[types.Codec]
}

// ErrorLogEvaluator is a helper interface for testing ErrorLogs
type ErrorLogEvaluator interface {
	types.ErrorLog
	Evaluator[types.ErrorLog]
}

type MedianProviderTester interface {
	types.MedianProvider
	Evaluator[types.MedianProvider]
	AssertEqualer[types.MedianProvider]
}

type RelayerTester interface {
	internal.PluginRelayer
	internal.Relayer
	// implements all the possible providers as a one-stop shop for testing
	internal.MercuryProvider
	internal.MedianProvider
	internal.CCIPExecProvider

	AssertEqualer[loop.Relayer]
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
