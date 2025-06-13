package fakes

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	httptypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ services.Service = (*FakeManualHttpTriggerService)(nil)
var _ ManualTriggerCapability = (*FakeManualHttpTriggerService)(nil)
var _ httpserver.HTTPCapability = (*FakeManualHttpTriggerService)(nil)

const HTTPTriggerServiceName = "HttpTriggerService"
const HTTPTriggerID = "http-trigger@1.0.0"

var fakeHttpTriggerInfo = capabilities.MustNewCapabilityInfo(
	HTTPTriggerID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that uses an HTTP request to run periodically at fixed times, dates, or intervals.",
)

type FakeManualHttpTriggerService struct {
	capabilities.CapabilityInfo
	lggr       logger.Logger
	callbackCh chan capabilities.TriggerAndId[*httptypedapi.Payload]
}

func NewFakeManualHttpTriggerService(parentLggr logger.Logger) *FakeManualHttpTriggerService {
	lggr := logger.Named(parentLggr, "HttpTriggerService")

	return &FakeManualHttpTriggerService{
		CapabilityInfo: fakeHttpTriggerInfo,
		lggr:           lggr,
		callbackCh:     make(chan capabilities.TriggerAndId[*httptypedapi.Payload]),
	}
}

// HTTPCapability interface methods
func (f *FakeManualHttpTriggerService) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *httptypedapi.Config) (<-chan capabilities.TriggerAndId[*httptypedapi.Payload], error) {
	return f.callbackCh, nil
}

func (f *FakeManualHttpTriggerService) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *httptypedapi.Config) error {
	return nil
}

func (f *FakeManualHttpTriggerService) Initialise(ctx context.Context, config string, telemetryService core.TelemetryService, store core.KeyValueStore, errorLog core.ErrorLog, pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet, oracleFactory core.OracleFactory) error {
	f.lggr.Debugf("Initialising %s", HTTPTriggerServiceName)
	return f.Start(ctx)
}

// ManualTriggerCapability interface method
func (f *FakeManualHttpTriggerService) ManualTrigger(ctx context.Context) error {
	// Run in a goroutine to avoid blocking
	go func() {
		// Send the trigger response
		f.callbackCh <- createFakeHttpTriggerResponse(time.Now())
	}()

	return nil
}

func createFakeHttpTriggerResponse(scheduledExecutionTime time.Time) capabilities.TriggerAndId[*httptypedapi.Payload] {
	return capabilities.TriggerAndId[*httptypedapi.Payload]{
		Trigger: &httptypedapi.Payload{
			Input: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"url": {
						Kind: &structpb.Value_StringValue{
							StringValue: "https://por.chainlink.com/por",
						},
					},
					"method": {
						Kind: &structpb.Value_StringValue{
							StringValue: "GET",
						},
					},
					"headers": {
						Kind: &structpb.Value_StructValue{
							StructValue: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"Content-Type": {
										Kind: &structpb.Value_StringValue{
											StringValue: "application/json",
										},
									},
								},
							},
						},
					},
				},
			},
			Key: &httptypedapi.AuthorizedKey{},
		},
		Id: "fake-http-trigger-id",
	}
}

// Service interface methods
func (f *FakeManualHttpTriggerService) Start(ctx context.Context) error {
	f.lggr.Info("Starting FakeManualHttpTriggerService")
	return nil
}

func (f *FakeManualHttpTriggerService) Close() error {
	f.lggr.Info("Closing FakeManualHttpTriggerService")
	return nil
}

func (f *FakeManualHttpTriggerService) Ready() error {
	return nil
}

func (f *FakeManualHttpTriggerService) HealthReport() map[string]error {
	return map[string]error{f.Name(): nil}
}

func (f *FakeManualHttpTriggerService) Name() string {
	return f.lggr.Name()
}

func (f *FakeManualHttpTriggerService) Description() string {
	return "Fake HTTP Trigger Service"
}
