package webhook

import (
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	Delegate struct {
		webhookJobRunner         *webhookJobRunner
		externalInitiatorManager ExternalInitiatorManager
	}

	JobRunner interface {
		RunJob(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta pipeline.JSONSerializable) (int64, error)
	}
)

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(runner pipeline.Runner, externalInitiatorManager ExternalInitiatorManager) *Delegate {
	return &Delegate{
		externalInitiatorManager: externalInitiatorManager,
		webhookJobRunner: &webhookJobRunner{
			specsByUUID: make(map[uuid.UUID]registeredJob),
			runner:      runner,
		},
	}
}

func (d *Delegate) WebhookJobRunner() JobRunner {
	return d.webhookJobRunner
}

func (d *Delegate) JobType() job.Type {
	return job.Webhook
}

func (d *Delegate) AfterJobCreated(jb job.Job) {
	err := d.externalInitiatorManager.NotifyV2(*jb.WebhookSpecID)
	if err != nil {
		logger.Errorw("Webhook delegate AfterJobCreated errored",
			"error", err,
			"jobID", jb.ID,
		)
	}
}

func (d *Delegate) BeforeJobDeleted(jb job.Job) {
	err := d.externalInitiatorManager.DeleteJobV2(*jb.WebhookSpecID)
	if err != nil {
		logger.Errorw("Webhook delegate BeforeJobDeleted errored",
			"error", err,
			"jobID", jb.ID,
		)
	}
}

func (d *Delegate) ServicesForSpec(spec job.Job) ([]job.Service, error) {
	// TODO: we need to fill these out manually, find a better fix
	spec.PipelineSpec.JobName = spec.Name.ValueOrZero()
	spec.PipelineSpec.JobID = spec.ID

	service := &pseudoService{
		spec:             spec,
		webhookJobRunner: d.webhookJobRunner,
	}
	return []job.Service{service}, nil
}

type pseudoService struct {
	spec             job.Job
	webhookJobRunner *webhookJobRunner
}

func (s pseudoService) Start() error {
	// add the spec to the webhookJobRunner
	return s.webhookJobRunner.addSpec(s.spec)
}

func (s pseudoService) Close() error {
	// remove the spec from the webhookJobRunner
	s.webhookJobRunner.rmSpec(s.spec)
	return nil
}

type webhookJobRunner struct {
	specsByUUID   map[uuid.UUID]registeredJob
	muSpecsByUUID sync.RWMutex
	runner        pipeline.Runner
}

type registeredJob struct {
	job.Job
	chRemove chan struct{}
}

func (r *webhookJobRunner) addSpec(spec job.Job) error {
	r.muSpecsByUUID.Lock()
	defer r.muSpecsByUUID.Unlock()
	_, exists := r.specsByUUID[spec.ExternalJobID]
	if exists {
		return errors.Errorf("a webhook job with that UUID already exists (uuid: %v)", spec.ExternalJobID)
	}
	r.specsByUUID[spec.ExternalJobID] = registeredJob{spec, make(chan struct{})}
	return nil
}

func (r *webhookJobRunner) rmSpec(spec job.Job) {
	r.muSpecsByUUID.Lock()
	defer r.muSpecsByUUID.Unlock()
	j, exists := r.specsByUUID[spec.ExternalJobID]
	if exists {
		close(j.chRemove)
		delete(r.specsByUUID, spec.ExternalJobID)
	}
}

func (r *webhookJobRunner) spec(externalJobID uuid.UUID) (registeredJob, bool) {
	r.muSpecsByUUID.RLock()
	defer r.muSpecsByUUID.RUnlock()
	spec, exists := r.specsByUUID[externalJobID]
	return spec, exists
}

var ErrJobNotExists = errors.New("job does not exist")

func (r *webhookJobRunner) RunJob(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta pipeline.JSONSerializable) (int64, error) {
	spec, exists := r.spec(jobUUID)
	if !exists {
		return 0, ErrJobNotExists
	}

	logger := logger.CreateLogger(
		logger.Default.With(
			"jobID", spec.ID,
			"uuid", spec.ExternalJobID,
		),
	)

	ctx, cancel := utils.CombinedContext(ctx, spec.chRemove)
	defer cancel()

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    spec.ID,
			"externalJobID": spec.ExternalJobID,
			"name":          spec.Name.ValueOrZero(),
		},
		"jobRun": map[string]interface{}{
			"requestBody": requestBody,
			"meta":        meta.Val,
		},
	})

	run := pipeline.NewRun(*spec.PipelineSpec, vars)

	_, err := r.runner.Run(ctx, &run, *logger, true)
	if err != nil {
		logger.Errorw("Error running pipeline for webhook job", "error", err)
		return 0, err
	}
	return run.ID, nil
}
