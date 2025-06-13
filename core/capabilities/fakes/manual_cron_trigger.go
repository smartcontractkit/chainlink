package fakes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	crontypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	cronserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ services.Service = (*FakeCronTriggerService)(nil)
var _ ManualTriggerCapability = (*FakeCronTriggerService)(nil)
var _ cronserver.CronCapability = (*FakeCronTriggerService)(nil)

const ServiceName = "CronTriggerService"
const ID = "cron-trigger@1.0.0"
const defaultFastestScheduleIntervalSeconds = 1

var fakeCronTriggerInfo = capabilities.MustNewCapabilityInfo(
	ID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that uses a cron schedule to run periodically at fixed times, dates, or intervals.",
)

type FakeCronConfig struct {
	FastestScheduleIntervalSeconds int `json:"fastestScheduleIntervalSeconds"`
}

type FakeCronTriggerService struct {
	capabilities.CapabilityInfo
	config     FakeCronConfig
	lggr       logger.Logger
	callbackCh chan capabilities.TriggerAndId[*crontypedapi.Payload]
}

func NewFakeCronTriggerService(parentLggr logger.Logger) *FakeCronTriggerService {
	lggr := logger.Named(parentLggr, "CronTriggerService") // FakeCronTriggerService

	return &FakeCronTriggerService{
		CapabilityInfo: fakeCronTriggerInfo,
		config:         FakeCronConfig{FastestScheduleIntervalSeconds: 1},
		lggr:           lggr,
		callbackCh:     make(chan capabilities.TriggerAndId[*crontypedapi.Payload]),
	}
}

func (f *FakeCronTriggerService) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory) error {
	f.lggr.Debugf("Initialising %s", ServiceName)

	var cronConfig FakeCronConfig
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

func (f *FakeCronTriggerService) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) (<-chan capabilities.TriggerAndId[*crontypedapi.Payload], error) {
	return f.callbackCh, nil
}

func (s *FakeCronTriggerService) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) error {
	return nil
}

func (f *FakeCronTriggerService) ManualTrigger(ctx context.Context) error {
	// Run in a goroutine to avoid blocking
	f.lggr.Debugf("ManualTrigger: %s", time.Now().Format(time.RFC3339Nano))
	go func() {
		// Send the trigger response
		f.callbackCh <- createFakeTriggerResponse(time.Now())
	}()

	return nil
}

func createFakeTriggerResponse(scheduledExecutionTime time.Time) capabilities.TriggerAndId[*crontypedapi.Payload] {
	// Ensure UTC time is used for consistency across nodes.
	scheduledExecutionTimeUTC := scheduledExecutionTime.UTC()

	// Use the scheduled execution time as a deterministic identifier.
	// Since cron schedules only go to second granularity this should never have ms.
	// Just in case, truncate on seconds by formatting to ensure consistency across nodes.
	scheduledExecutionTimeFormatted := scheduledExecutionTimeUTC.Format(time.RFC3339)
	triggerEventID := scheduledExecutionTimeFormatted

	return capabilities.TriggerAndId[*crontypedapi.Payload]{
		Trigger: &crontypedapi.Payload{
			ScheduledExecutionTime: scheduledExecutionTimeUTC.Format(time.RFC3339Nano),
		},
		Id: triggerEventID,
	}
}

func (f *FakeCronTriggerService) Start(ctx context.Context) error {
	f.lggr.Info("Starting FakeCronTriggerService")
	return nil
}

func (f *FakeCronTriggerService) Close() error {
	f.lggr.Info("Closing FakeCronTriggerService")
	return nil
}

func (f *FakeCronTriggerService) Ready() error {
	return nil
}

func (f *FakeCronTriggerService) HealthReport() map[string]error {
	return map[string]error{f.Name(): nil}
}

func (f *FakeCronTriggerService) Name() string {
	return f.lggr.Name()
}

func (f *FakeCronTriggerService) Description() string {
	return "Fake Cron Trigger Service"
}
