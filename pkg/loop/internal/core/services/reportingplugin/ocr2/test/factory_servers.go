package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pipelinetest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/pipeline/test"
	telemetrytest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/telemetry/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	mediantest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median/test"
	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
	reportingplugintest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/test"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var MedianProviderServer = medianFactoryServer{
	medianGeneratorConfig: medianGeneratorConfig{
		medianProvider: mediantest.MedianProvider,
		pipeline:       pipelinetest.PipelineRunner,
		telemetry:      telemetrytest.Telemetry,
	},
}

const MedianID = "ocr2-reporting-plugin-with-median-provider"

type medianGeneratorConfig struct {
	medianProvider testtypes.MedianProviderTester
	pipeline       testtypes.Evaluator[types.PipelineRunnerService]
	telemetry      testtypes.Evaluator[types.TelemetryClient]
}

type medianFactoryServer struct {
	medianGeneratorConfig
}

var _ reportingplugins.ProviderServer[types.MedianProvider] = medianFactoryServer{}

func (s medianFactoryServer) ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerConfig net.BrokerConfig) types.MedianProvider {
	return s.medianProvider
}

func (s medianFactoryServer) NewReportingPluginFactory(ctx context.Context, config types.ReportingPluginServiceConfig, provider types.MedianProvider, pipelineRunner types.PipelineRunnerService, telemetry types.TelemetryClient, errorLog types.ErrorLog, keyValueStore types.KeyValueStore) (types.ReportingPluginFactory, error) {
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

	return reportingplugintest.Factory, nil
}

var AgnosticProviderServer = agnosticPluginFactoryServer{
	provider:       ocr2test.AgnosticPluginProvider,
	pipelineRunner: pipelinetest.PipelineRunner,
	telemetry:      telemetrytest.Telemetry,
}

var _ reportingplugins.ProviderServer[types.PluginProvider] = agnosticPluginFactoryServer{}

type agnosticPluginFactoryServer struct {
	provider       testtypes.PluginProviderTester
	pipelineRunner testtypes.PipelineEvaluator
	telemetry      testtypes.TelemetryEvaluator
}

func (s agnosticPluginFactoryServer) ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerConfig net.BrokerConfig) types.PluginProvider {
	return s.provider
}

func (s agnosticPluginFactoryServer) NewReportingPluginFactory(ctx context.Context, config types.ReportingPluginServiceConfig, provider types.PluginProvider, pipelineRunner types.PipelineRunnerService, telemetry types.TelemetryClient, errorLog types.ErrorLog, keyValueStore types.KeyValueStore) (types.ReportingPluginFactory, error) {
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

	return reportingplugintest.Factory, nil
}
