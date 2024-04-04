package ocr3_test

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pipelinetest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline/test"
	telemetrytest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry/test"
	validationtest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/validation/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	mediantest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median/test"
	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var MedianServer = medianServer{
	medianGeneratorConfig: medianGeneratorConfig{
		medianProvider:    mediantest.MedianProvider,
		pipeline:          pipelinetest.PipelineRunner,
		telemetry:         telemetrytest.Telemetry,
		validationService: validationtest.ValidationService,
	},
}

const OCR3ReportingPluginWithMedianProviderName = "ocr3-reporting-plugin-with-median-provider"

type medianGeneratorConfig struct {
	medianProvider    testtypes.MedianProviderTester
	pipeline          testtypes.Evaluator[types.PipelineRunnerService]
	telemetry         testtypes.Evaluator[types.TelemetryClient]
	validationService testtypes.ValidationEvaluator
}

type medianServer struct {
	medianGeneratorConfig
}

func (s medianServer) NewValidationService(ctx context.Context) (types.ValidationService, error) {
	return s.validationService, nil
}
func (s medianServer) ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerConfig net.BrokerConfig) types.MedianProvider {
	return s.medianProvider
}

func (s medianServer) NewReportingPluginFactory(ctx context.Context, config types.ReportingPluginServiceConfig, provider types.MedianProvider, pipelineRunner types.PipelineRunnerService, telemetry types.TelemetryClient, errorLog types.ErrorLog, capRegistry types.CapabilitiesRegistry, keyValueStore types.KeyValueStore) (types.OCR3ReportingPluginFactory, error) {
	err := s.medianProvider.Evaluate(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate median provider: %w", err)
	}

	err = s.pipeline.Evaluate(ctx, pipelineRunner)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate pipeline runner: %w", err)
	}

	err = s.telemetry.Evaluate(ctx, telemetry)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate telemetry: %w", err)
	}

	return Factory, nil
}

var AgnosticPluginServer = agnosticPluginServer{
	provider:          ocr2test.AgnosticPluginProvider,
	pipelineRunner:    pipelinetest.PipelineRunner,
	telemetry:         telemetrytest.Telemetry,
	validationService: validationtest.ValidationService,
}

type agnosticPluginServer struct {
	provider          testtypes.PluginProviderTester
	pipelineRunner    testtypes.PipelineEvaluator
	telemetry         testtypes.TelemetryEvaluator
	validationService testtypes.ValidationEvaluator
}

func (s agnosticPluginServer) NewValidationService(ctx context.Context) (types.ValidationService, error) {
	return s.validationService, nil
}

func (s agnosticPluginServer) ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerConfig net.BrokerConfig) types.PluginProvider {
	return s.provider
}

func (s agnosticPluginServer) NewReportingPluginFactory(ctx context.Context, config types.ReportingPluginServiceConfig, provider types.PluginProvider, pipelineRunner types.PipelineRunnerService, telemetry types.TelemetryClient, errorLog types.ErrorLog, capRegistry types.CapabilitiesRegistry, keyValueStore types.KeyValueStore) (types.OCR3ReportingPluginFactory, error) {
	err := s.provider.Evaluate(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate agnostic provider: %w", err)
	}

	err = s.pipelineRunner.Evaluate(ctx, pipelineRunner)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate pipeline runner: %w", err)
	}

	err = s.telemetry.Evaluate(ctx, telemetry)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate telemetry: %w", err)
	}

	return Factory, nil
}
