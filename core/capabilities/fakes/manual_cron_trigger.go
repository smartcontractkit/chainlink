package fakes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	crontypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	cronserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

var _ services.Service = (*ManualCronTriggerService)(nil)
var _ cronserver.CronCapability = (*ManualCronTriggerService)(nil)

const ServiceName = "CronTriggerService"
const ID = "cron-trigger@1.0.0"
const defaultFastestScheduleIntervalSeconds = 1
const allowSeconds = true

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
	workflowIDs      map[string]string                                           // triggerID -> workflowID mapping
	triggerConfigs   map[string]*crontypedapi.Config
	scheduler        gocron.Scheduler
}

func NewManualCronTriggerService(parentLggr logger.Logger) (*ManualCronTriggerService, error) {
	lggr := logger.Named(parentLggr, "CronTriggerService") // ManualCronTriggerService

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create cron scheduler: %w", err)
	}

	return &ManualCronTriggerService{
		CapabilityInfo:   manualCronTriggerInfo,
		config:           ManualCronConfig{FastestScheduleIntervalSeconds: 1},
		lggr:             lggr,
		callbackCh:       make(map[string]chan capabilities.TriggerAndId[*crontypedapi.Payload]),
		legacyCallbackCh: make(chan capabilities.TriggerAndId[*crontypedapi.LegacyPayload]), //nolint:staticcheck // LegacyPayload intentionally used for backward compatibility
		workflowIDs:      make(map[string]string),
		triggerConfigs:   make(map[string]*crontypedapi.Config),
		scheduler:        scheduler,
	}, nil
}

func (f *ManualCronTriggerService) Initialise(ctx context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	f.lggr.Debugf("Initialising %s", ServiceName)

	var cronConfig ManualCronConfig
	if len(dependencies.Config) > 0 {
		err := json.Unmarshal([]byte(dependencies.Config), &cronConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %s %w", dependencies.Config, err)
		}
	}

	if cronConfig.FastestScheduleIntervalSeconds == 0 {
		cronConfig.FastestScheduleIntervalSeconds = defaultFastestScheduleIntervalSeconds
	}

	f.config = cronConfig

	if err := f.Start(ctx); err != nil {
		return fmt.Errorf("error when starting trigger service: %w", err)
	}

	return nil
}

func (f *ManualCronTriggerService) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) (<-chan capabilities.TriggerAndId[*crontypedapi.Payload], caperrors.Error) {
	f.callbackCh[triggerID] = make(chan capabilities.TriggerAndId[*crontypedapi.Payload], 1)
	f.workflowIDs[triggerID] = metadata.WorkflowID
	f.triggerConfigs[triggerID] = input
	return f.callbackCh[triggerID], nil
}

func (f *ManualCronTriggerService) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) caperrors.Error {
	return nil
}

func (f *ManualCronTriggerService) RegisterLegacyTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) (<-chan capabilities.TriggerAndId[*crontypedapi.LegacyPayload], caperrors.Error) { //nolint:staticcheck // LegacyPayload intentionally used for backward compatibility
	return f.legacyCallbackCh, nil
}

func (f *ManualCronTriggerService) UnregisterLegacyTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) caperrors.Error {
	return nil
}

func (f *ManualCronTriggerService) AckEvent(ctx context.Context, triggerID string, eventID string, method string) caperrors.Error {
	return nil
}

func (f *ManualCronTriggerService) ManualTrigger(ctx context.Context, triggerID string, skipWait <-chan struct{}) error {
	config, exists := f.triggerConfigs[triggerID]
	if !exists {
		return fmt.Errorf(`trigger config "%s" not found`, triggerID)
	}

	jobFired := make(chan struct{}, 1)
	job, err := f.scheduler.NewJob(
		gocron.CronJob(config.Schedule, allowSeconds),
		gocron.NewTask(func() {
			defer close(jobFired)
			jobFired <- struct{}{}
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create cron job: %w", err)
	}
	scheduledExecutionTime, err := job.NextRun()
	if err != nil {
		return fmt.Errorf("failed to get next scheduled execution time: %w", err)
	}

	f.lggr.Debugf("ManualTrigger: %s", scheduledExecutionTime.Format(time.RFC3339Nano))

	triggerEvent := f.createManualTriggerEvent(scheduledExecutionTime)

	// Get the workflowID for this trigger
	workflowID, exists := f.workflowIDs[triggerID]
	if !exists {
		f.lggr.Errorw("workflowID not found for triggerID", "triggerID", triggerID)
		workflowID = "unknownWorkflow"
	}

	// Emit trigger execution started event with real workflowExecutionID
	workflowExecutionID, err := events.GenerateExecutionID(workflowID, triggerEvent.Id)
	if err != nil {
		f.lggr.Errorw("failed to generate execution ID", "err", err)
		workflowExecutionID = ""
	}
	err = events.EmitTriggerExecutionStarted(ctx, map[string]string{}, triggerEvent.Id, workflowExecutionID)
	if err != nil {
		f.lggr.Errorw("failed to emit trigger execution started event", "err", err)
	}

	defer func() {
		_ = f.scheduler.RemoveJob(job.ID())
	}()

	// Either wait for cron scheduler or skip wait signal
	select {
	case <-skipWait:
		break
	case <-jobFired:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	// Sent trigger response
	f.callbackCh[triggerID] <- triggerEvent
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
	f.scheduler.Start()
	return nil
}

func (f *ManualCronTriggerService) Close() error {
	f.lggr.Debug("Closing ManualCronTriggerService")
	if err := f.scheduler.Shutdown(); err != nil {
		f.lggr.Errorw("failed to close scheduler", "err", err)
	}
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
