package fakes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/cron"
	crontypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

// TODO(CAPPL-866): remove this copy of cron trigger implementation from fakes
const ID = "cron-trigger@1.0.0"
const ServiceName = "CronCapabilities"

const (
	defaultSendChannelBufferSize          = 1
	defaultFastestScheduleIntervalSeconds = 30
)

var cronTriggerInfo = capabilities.MustNewCapabilityInfo(
	ID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that uses a cron schedule to run periodically at fixed times, dates, or intervals.",
)

type Config struct {
	FastestScheduleIntervalSeconds int `json:"fastestScheduleIntervalSeconds"`
}

type Response struct {
	capabilities.TriggerEvent
	Payload cron.Payload
}

type cronTrigger struct {
	ch      chan<- capabilities.TriggerAndId[*crontypedapi.Payload]
	job     gocron.Job
	nextRun time.Time
}

type Service struct {
	capabilities.CapabilityInfo
	config    Config
	clock     clockwork.Clock
	lggr      logger.Logger
	scheduler gocron.Scheduler
	triggers  *cronStore
}

var _ services.Service = &Service{}

// NewTriggerService creates a new trigger service.  Optionally, a clock can be passed in for testing, if nil
// the system clock will be used.
func NewTriggerService(parentLggr logger.Logger, clock clockwork.Clock) *Service {
	lggr := logger.Named(parentLggr, "Service")

	var options []gocron.SchedulerOption
	// Set scheduler location to UTC for consistency across nodes.
	options = append(options, gocron.WithLocation(time.UTC))

	// Allow injecting a clock for testing. Otherwise use system clock.
	if clock != nil {
		options = append(options, gocron.WithClock(clock))
	} else {
		clock = clockwork.NewRealClock()
	}

	scheduler, err := gocron.NewScheduler(options...)
	if err != nil {
		return nil
	}

	return &Service{
		lggr:           lggr,
		CapabilityInfo: cronTriggerInfo,
		triggers:       NewCronStore(),
		scheduler:      scheduler,
		clock:          clock,
	}
}

func (s *Service) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory) error {
	s.lggr.Debugf("Initialising %s", ServiceName)

	var cronConfig Config
	if len(config) > 0 {
		err := json.Unmarshal([]byte(config), &cronConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %s %w", config, err)
		}
	}

	if cronConfig.FastestScheduleIntervalSeconds == 0 {
		cronConfig.FastestScheduleIntervalSeconds = defaultFastestScheduleIntervalSeconds
	}

	s.config = cronConfig

	err := s.Start(ctx)
	if err != nil {
		return fmt.Errorf("error when starting trigger service: %w", err)
	}

	return nil
}

func (s *Service) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) (<-chan capabilities.TriggerAndId[*crontypedapi.Payload], error) {
	_, ok := s.triggers.Read(triggerID)
	if ok {
		return nil, fmt.Errorf("triggerId %s already registered", triggerID)
	}

	var job gocron.Job
	callbackCh := make(chan capabilities.TriggerAndId[*crontypedapi.Payload], defaultSendChannelBufferSize)

	allowSeconds := true
	jobDef := gocron.CronJob(input.Schedule, allowSeconds)

	err := enforceFastestSchedule(s.lggr, s.clock, jobDef, time.Second*time.Duration(s.config.FastestScheduleIntervalSeconds))
	if err != nil {
		return nil, err
	}

	task := gocron.NewTask(
		// Task callback, executed at next run time
		func() {
			trigger, ok := s.triggers.Read(triggerID)
			if !ok {
				// Invariant: The trigger should always exist, as unregistering the trigger removes the job
				s.lggr.Errorw("task callback invariant: trigger no longer exists", "triggerID", triggerID)
				return
			}

			scheduledExecutionTimeUTC := trigger.nextRun.UTC()
			currentTimeUTC := s.clock.Now().UTC()

			response := createTriggerResponse(scheduledExecutionTimeUTC)

			s.lggr.Debugw("task callback sending trigger response", "executionID", metadata.WorkflowExecutionID, "triggerID", triggerID, "scheduledExecTimeUTC", scheduledExecutionTimeUTC.Format(time.RFC3339Nano), "actualExecTimeUTC", currentTimeUTC.Format(time.RFC3339Nano))

			nextExecutionTime, nextRunErr := job.NextRun()
			if nextRunErr != nil {
				// .NextRun() will error if the job no longer exists
				// or if there is no next run to schedule, which shouldn't happen with cron jobs
				s.lggr.Errorw("task callback failed to schedule next run", "executionID", metadata.WorkflowExecutionID, "triggerID", triggerID)
			}

			s.triggers.Write(triggerID, cronTrigger{
				ch:      callbackCh,
				job:     job,
				nextRun: nextExecutionTime,
			})

			select {
			case callbackCh <- response:
			default:
				s.lggr.Errorw("callback channel full, dropping event", "executionID", metadata.WorkflowExecutionID, "triggerID", triggerID, "eventID", response.Id)
			}
		})

	if s.scheduler == nil {
		return nil, errors.New("cannot register a new trigger, service has been closed")
	}

	// If service has already started, job will be scheduled immediately
	job, err = s.scheduler.NewJob(jobDef, task, gocron.WithName(triggerID))
	if err != nil {
		s.lggr.Errorw("failed to create new job", "err", err)
		return nil, err
	}

	firstRunTime, err := job.NextRun()
	if err != nil {
		// errors if job no longer exists on scheduler
		s.lggr.Errorw("failed to get next run time", "err", err)
		// ensure that it is out of scheduler
		err := s.scheduler.RemoveJob(job.ID())
		return nil, err
	}

	s.triggers.Write(triggerID, cronTrigger{
		ch:      callbackCh,
		job:     job,
		nextRun: firstRunTime,
	})

	s.lggr.Debugw("Trigger registered", "workflowId", metadata.WorkflowID, "triggerId", triggerID, "jobId", job.ID())

	return callbackCh, nil
}

func createTriggerResponse(scheduledExecutionTime time.Time) capabilities.TriggerAndId[*crontypedapi.Payload] {
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

func (s *Service) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *crontypedapi.Config) error {
	s.lggr.Debug("got call to unregister trigger")
	trigger, ok := s.triggers.Read(triggerID)
	if !ok {
		return fmt.Errorf("triggerId %s not found", triggerID)
	}

	jobID := trigger.job.ID()

	// Remove job from scheduler
	if s.scheduler == nil {
		return errors.New("cannot unregister a new trigger, service has been closed")
	}
	err := s.scheduler.RemoveJob(jobID)
	if err != nil {
		return fmt.Errorf("UnregisterTrigger failed to remove job from scheduler: %w", err)
	}

	// Close callback channel
	s.lggr.Debug("closing event channel")
	close(trigger.ch)

	// Remove from triggers context
	s.triggers.Delete(triggerID)

	s.lggr.Debugw("UnregisterTrigger", "triggerId", triggerID, "jobId", jobID)
	return nil
}

// Start the service.
func (s *Service) Start(ctx context.Context) error {
	if s.scheduler == nil {
		return errors.New("service has shutdown, it must be built again to restart")
	}

	s.scheduler.Start()

	for triggerID, trigger := range s.triggers.ReadAll() {
		nextExecutionTime, err := trigger.job.NextRun()
		s.triggers.Write(triggerID, cronTrigger{
			ch:      trigger.ch,
			job:     trigger.job,
			nextRun: nextExecutionTime,
		})
		if err != nil {
			s.lggr.Errorw("Unable to get next run time", "err", err, "triggerID", triggerID)
		}
	}

	s.lggr.Info(s.Name() + " started")

	return nil
}

// Close stops the Service.
// After this call the Service cannot be started again,
// The service will need to be re-built to start scheduling again.
func (s *Service) Close() error {
	if s.scheduler == nil {
		return errors.New("service has shutdown, it must be built again to restart")
	}

	err := s.scheduler.Shutdown()
	if err != nil {
		return fmt.Errorf("scheduler shutdown encountered a problem: %w", err)
	}

	// After .Shutdown() the scheduler cannot be started again,
	// but calling .Start() on it will not error. Set to nil to mark closed.
	s.scheduler = nil

	s.lggr.Info(s.Name() + " closed")

	return nil
}

func (s *Service) Ready() error {
	return nil
}

func (s *Service) HealthReport() map[string]error {
	return map[string]error{s.Name(): nil}
}

func (s *Service) Name() string {
	return s.lggr.Name()
}

func (s *Service) Description() string {
	return "Cron Trigger Capability"
}

func enforceFastestSchedule(lggr logger.Logger, clock clockwork.Clock, jobDef gocron.JobDefinition, maximumFastest time.Duration) error {
	var options []gocron.SchedulerOption
	// Set scheduler location to UTC for consistency across nodes.
	options = append(options, gocron.WithLocation(time.UTC))
	// Use passed in clock
	options = append(options, gocron.WithClock(clock))

	tempScheduler, err := gocron.NewScheduler(options...)
	if err != nil {
		return err
	}
	tempJob, err := tempScheduler.NewJob(jobDef, gocron.NewTask(func() {}))
	if err != nil {
		return err
	}
	tempScheduler.Start()
	defer func() {
		if err = tempScheduler.Shutdown(); err != nil {
			lggr.Errorw("error shutting down enforceFastestSchedule temporary scheduler")
		}
	}()

	nextRuns, err := tempJob.NextRuns(2)
	if err != nil {
		return err
	}

	if len(nextRuns) != 2 {
		return errors.New("could not determine next two scheduled runs")
	}

	if nextRuns[1].Before(nextRuns[0].Add(maximumFastest)) {
		return fmt.Errorf("maximum fastest cron schedule is %s", maximumFastest.String())
	}

	return nil
}
