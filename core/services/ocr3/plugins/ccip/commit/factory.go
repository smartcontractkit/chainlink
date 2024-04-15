package commit

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// PluginFactoryConstructor implements common OCR3ReportingPluginClient and is used for initializing a plugin factory
// and a validation service.
type PluginFactoryConstructor struct{}

func NewPluginFactoryConstructor() *PluginFactoryConstructor {
	return &PluginFactoryConstructor{}
}
func (p PluginFactoryConstructor) NewReportingPluginFactory(
	ctx context.Context,
	config types.ReportingPluginServiceConfig,
	grpcProvider grpc.ClientConnInterface,
	pipelineRunner types.PipelineRunnerService,
	telemetry types.TelemetryService,
	errorLog types.ErrorLog,
	capRegistry types.CapabilitiesRegistry,
	keyValueStore types.KeyValueStore,
) (types.OCR3ReportingPluginFactory, error) {
	return NewPluginFactory(), nil
}

func (p PluginFactoryConstructor) NewValidationService(ctx context.Context) (types.ValidationService, error) {
	panic("implement me")
}

// PluginFactory implements common ReportingPluginFactory and is used for (re-)initializing commit plugin instances.
type PluginFactory struct{}

func NewPluginFactory() *PluginFactory {
	return &PluginFactory{}
}

func (p PluginFactory) NewReportingPlugin(config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	return NewPlugin(), ocr3types.ReportingPluginInfo{}, nil
}

func (p PluginFactory) Name() string {
	panic("implement me")
}

func (p PluginFactory) Start(ctx context.Context) error {
	panic("implement me")
}

func (p PluginFactory) Close() error {
	panic("implement me")
}

func (p PluginFactory) Ready() error {
	panic("implement me")
}

func (p PluginFactory) HealthReport() map[string]error {
	panic("implement me")
}

// Interface compatibility checks.
var _ types.OCR3ReportingPluginClient = &PluginFactoryConstructor{}
var _ types.OCR3ReportingPluginFactory = &PluginFactory{}
