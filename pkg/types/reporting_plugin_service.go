package types

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/grpc"
)

type ReportingPluginServiceConfig struct {
	ProviderType  string
	Command       string
	PluginName    string
	TelemetryType string
	PluginConfig  string
}

// ReportingPluginClient is the client interface to a plugin running
// as a generic job (job type = GenericPlugin) inside the core node.
type ReportingPluginClient interface {
	NewReportingPluginFactory(ctx context.Context, config ReportingPluginServiceConfig, grpcProvider grpc.ClientConnInterface, pipelineRunner PipelineRunnerService, telemetry TelemetryService, errorLog ErrorLog) (ReportingPluginFactory, error)
}

// ReportingPluginServer is the server interface to a plugin running
// as a generic job (job type = GenericPlugin) inside the core node,
// with the passthrough provider connection converted to the provider
// expected by the plugin.
type ReportingPluginServer[T PluginProvider] interface {
	NewReportingPluginFactory(ctx context.Context, config ReportingPluginServiceConfig, provider T, pipelineRunner PipelineRunnerService, telemetry TelemetryClient, errorLog ErrorLog) (ReportingPluginFactory, error)
}

type OCR3ReportingPluginClient interface {
	NewReportingPluginFactory(ctx context.Context, config ReportingPluginServiceConfig, grpcProvider grpc.ClientConnInterface, pipelineRunner PipelineRunnerService, telemetry TelemetryService, errorLog ErrorLog, capRegistry CapabilitiesRegistry) (OCR3ReportingPluginFactory, error)
}

type OCR3ReportingPluginServer[T PluginProvider] interface {
	NewReportingPluginFactory(ctx context.Context, config ReportingPluginServiceConfig, provider T, pipelineRunner PipelineRunnerService, telemetry TelemetryClient, errorLog ErrorLog, capRegistry CapabilitiesRegistry) (OCR3ReportingPluginFactory, error)
}

type ReportingPluginFactory interface {
	Service
	libocr.ReportingPluginFactory
}

type OCR3ReportingPluginFactory interface {
	Service
	ocr3types.ReportingPluginFactory[[]byte]
}
