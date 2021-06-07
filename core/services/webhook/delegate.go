package webhook

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	Delegate struct {
		webhookJobRunner         *webhookJobRunner
		externalInitiatorManager ExternalInitiatorManager
	}

	JobRunner interface {
		RunJob(ctx context.Context, jobUUID models.JobID, pipelineInput interface{}, meta pipeline.JSONSerializable) (int64, error)
	}
)

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(runner pipeline.Runner, externalInitiatorManager ExternalInitiatorManager) *Delegate {
	return &Delegate{
		externalInitiatorManager: externalInitiatorManager,
		webhookJobRunner: &webhookJobRunner{
			specsByUUID: make(map[models.JobID]registeredJob),
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

func (d *Delegate) OnJobCreated(spec job.Job) {
	if spec.WebhookSpec.ExternalInitiatorName != "" {
		err := d.externalInitiatorManager.NotifyV2(
			spec.WebhookSpec.OnChainJobSpecID,
			spec.WebhookSpec.ExternalInitiatorName,
			spec.WebhookSpec.ExternalInitiatorSpec,
		)
		if err != nil {
			logger.Errorw("Webhook delegate OnJobCreated errored",
				"error", err,
				"jobID", spec.ID,
			)
		}
	}
}

func (d *Delegate) OnJobDeleted(spec job.Job) {
	if spec.WebhookSpec.ExternalInitiatorName != "" {
		err := d.externalInitiatorManager.DeleteJobV2(spec.WebhookSpec.OnChainJobSpecID)
		if err != nil {
			logger.Errorw("Webhook delegate OnJobDeleted errored",
				"error", err,
				"jobID", spec.ID,
			)
		}
	}
}

func (d *Delegate) ServicesForSpec(spec job.Job) ([]job.Service, error) {
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
	specsByUUID   map[models.JobID]registeredJob
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
	_, exists := r.specsByUUID[spec.WebhookSpec.OnChainJobSpecID]
	if exists {
		return errors.Errorf("a webhook job with that UUID already exists (uuid: %v)", spec.WebhookSpec.OnChainJobSpecID)
	}
	r.specsByUUID[spec.WebhookSpec.OnChainJobSpecID] = registeredJob{spec, make(chan struct{})}
	return nil
}

func (r *webhookJobRunner) rmSpec(spec job.Job) {
	r.muSpecsByUUID.Lock()
	defer r.muSpecsByUUID.Unlock()
	j, exists := r.specsByUUID[spec.WebhookSpec.OnChainJobSpecID]
	if exists {
		close(j.chRemove)
		delete(r.specsByUUID, spec.WebhookSpec.OnChainJobSpecID)
	}
}

func (r *webhookJobRunner) spec(jobID models.JobID) (registeredJob, bool) {
	r.muSpecsByUUID.RLock()
	defer r.muSpecsByUUID.RUnlock()
	spec, exists := r.specsByUUID[jobID]
	return spec, exists
}

var ErrJobNotExists = errors.New("job does not exist")

func (r *webhookJobRunner) RunJob(ctx context.Context, jobUUID models.JobID, pipelineInput interface{}, meta pipeline.JSONSerializable) (int64, error) {
	spec, exists := r.spec(jobUUID)
	if !exists {
		return 0, ErrJobNotExists
	}

	logger := logger.CreateLogger(
		logger.Default.With(
			"jobID", spec.ID,
			"uuid", spec.WebhookSpec.OnChainJobSpecID,
		),
	)

	ctx, cancel := utils.CombinedContext(ctx, spec.chRemove)
	defer cancel()

	runID, _, err := r.runner.ExecuteAndInsertFinishedRun(ctx, *spec.PipelineSpec, pipelineInput, meta, *logger, true)
	if err != nil {
		logger.Errorw("Error running pipeline for webhook job", "error", err)
		return 0, err
	}
	return runID, nil
}
