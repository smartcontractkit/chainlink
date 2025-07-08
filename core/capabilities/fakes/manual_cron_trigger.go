package fakes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	crontypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	cronserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ services.Service = (*ManualCronTriggerService)(nil)
var _ cronserver.CronCapability = (*ManualCronTriggerService)(nil)

const ServiceName = "CronTriggerService"
const ID = "cron-trigger@1.0.0"
const defaultFastestScheduleIntervalSeconds = 1

var manualCronTriggerInfo = capabilities.MustNewCapabilityInfo(
	ID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that uses a cron schedule to run periodically at fixed times, dates, or intervals.",
)

type ManualCronConfig struct {
	FastestScheduleIntervalSeconds int `json:"fastestScheduleIntervalSeconds"`
}

type ManualCronTriggerService struct {
	capabilities.CapabilityInfo
	config           ManualCronConfig
	lggr             logger.Logger
	callbackCh       map[string]chan capabilities.TriggerAndId[*crontypedapi.Payload]
	legacyCallbackCh chan capabilities.TriggerAndId[*crontypedapi.LegacyPayload] //nolint:staticcheck // LegacyPayload intentionally used for backward compatibility
}

func NewManualCronTriggerService(parentLggr logger.Logger) *ManualCronTriggerService {
	lggr := logger.Named(parentLggr, "CronTriggerService") // ManualCronTriggerService

	return &ManualCronTriggerService{
		CapabilityInfo:   manualCronTriggerInfo,
		config:           ManualCronConfig{FastestScheduleIntervalSeconds: 1},
		lggr:             lggr,
		callbackCh:       make(map[string]chan capabilities.TriggerAndId[*crontypedapi.Payload]),
		legacyCallbackCh: make(chan capabilities.TriggerAndId[*crontypedapi.LegacyPayload]), //nolint:staticcheck // LegacyPayload intentionally used for backward compatibility
	}
}

func (f *ManualCronTriggerService) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory,
	_ core.GatewayConnector,
	_ core.Keystore) error {
	f.lggr.Debugf("Initialising %s", ServiceName)

	var cronConfig ManualCronConfig
	if len(config) > 0 {
		err := json.Unmarshal([]byte(config), &cronConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %s %w", config, err)
		}
	}

	if cronConfig.FastestScheduleIntervalSeconds == 0 {
		cronConfig.FastestScheduleIntervalSeconds = defaultFastestScheduleIntervalSeconds
	}

	f.config = cronConfig

	err := f.Start(ctx)
	if err != nil {
		return fmt.Errorf("error when starting trigger service: %w", err)
	}

	return nil
}

func (f *ManualCronTriggerService) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) (<-chan capabilities.TriggerAndId[*crontypedapi.Payload], error) {
	f.callbackCh[triggerID] = make(chan capabilities.TriggerAndId[*crontypedapi.Payload])
	return f.callbackCh[triggerID], nil
}

func (f *ManualCronTriggerService) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) error {
	return nil
}

func (f *ManualCronTriggerService) RegisterLegacyTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) (<-chan capabilities.TriggerAndId[*crontypedapi.LegacyPayload], error) { //nolint:staticcheck // LegacyPayload intentionally used for backward compatibility
	return f.legacyCallbackCh, nil
}

func (f *ManualCronTriggerService) UnregisterLegacyTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) error {
	return nil
}

func (f *ManualCronTriggerService) ManualTrigger(ctx context.Context, triggerID string, scheduledExecutionTime time.Time) error {
	f.lggr.Debugf("ManualTrigger: %s", scheduledExecutionTime.Format(time.RFC3339Nano))

	go func() {
		select {
		case f.callbackCh[triggerID] <- f.createManualTriggerEvent(scheduledExecutionTime):
			// Successfully sent trigger response
		case <-ctx.Done():
			// Context cancelled, cleanup goroutine
			f.lggr.Debug("ManualTrigger goroutine cancelled due to context cancellation")
		}
	}()

	return nil
}

func (f *ManualCronTriggerService) createManualTriggerEvent(scheduledExecutionTime time.Time) capabilities.TriggerAndId[*crontypedapi.Payload] {
	// Ensure UTC time is used for consistency across nodes.
	scheduledExecutionTimeUTC := scheduledExecutionTime.UTC()

	// Use the scheduled execution time as a deterministic identifier.
	// Since cron schedules only go to second granularity this should never have ms.
	// Just in case, truncate on seconds by formatting to ensure consistency across nodes.
	scheduledExecutionTimeFormatted := scheduledExecutionTimeUTC.Format(time.RFC3339)
	triggerEventID := scheduledExecutionTimeFormatted

	return capabilities.TriggerAndId[*crontypedapi.Payload]{
		Trigger: &crontypedapi.Payload{
			ScheduledExecutionTime: timestamppb.New(scheduledExecutionTimeUTC),
		},
		Id: triggerEventID,
	}
}

func (f *ManualCronTriggerService) Start(ctx context.Context) error {
	f.lggr.Debugw("Starting ManualCronTriggerService")
	return nil
}

func (f *ManualCronTriggerService) Close() error {
	f.lggr.Debug("Closing ManualCronTriggerService")
	return nil
}

func (f *ManualCronTriggerService) Ready() error {
	return nil
}

func (f *ManualCronTriggerService) HealthReport() map[string]error {
	return map[string]error{f.Name(): nil}
}

func (f *ManualCronTriggerService) Name() string {
	return f.lggr.Name()
}

func (f *ManualCronTriggerService) Description() string {
	return "Manual Cron Trigger Service"
}
