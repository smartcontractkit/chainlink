package webhook

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type (
	Delegate struct {
		webhookJobRunner         *webhookJobRunner
		externalInitiatorManager ExternalInitiatorManager
		lggr                     logger.Logger
		stopCh                   services.StopChan
	}

	JobRunner interface {
		RunJob(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error)
	}
)

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(runner pipeline.Runner, externalInitiatorManager ExternalInitiatorManager, lggr logger.Logger) *Delegate {
	lggr = lggr.Named("Webhook")
	return &Delegate{
		externalInitiatorManager: externalInitiatorManager,
		webhookJobRunner:         newWebhookJobRunner(runner, lggr),
		lggr:                     lggr,
		stopCh:                   make(services.StopChan),
	}
}

func (d *Delegate) WebhookJobRunner() JobRunner {
	return d.webhookJobRunner
}

func (d *Delegate) JobType() job.Type {
	return job.Webhook
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}
func (d *Delegate) AfterJobCreated(jb job.Job) {
	ctx, cancel := d.stopCh.NewCtx()
	defer cancel()
	err := d.externalInitiatorManager.Notify(ctx, *jb.WebhookSpecID)
	if err != nil {
		d.lggr.Errorw("Webhook delegate AfterJobCreated errored",
			"err", err,
			"jobID", jb.ID,
		)
	}
}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {
	ctx, cancel := d.stopCh.NewCtx()
	defer cancel()
	err := d.externalInitiatorManager.DeleteJob(ctx, *spec.WebhookSpecID)
	if err != nil {
		d.lggr.Errorw("Webhook delegate OnDeleteJob errored",
			"err", err,
			"jobID", spec.ID,
		)
	}
}
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	service := &pseudoService{
		spec:             spec,
		webhookJobRunner: d.webhookJobRunner,
	}
	return []job.ServiceCtx{service}, nil
}

type pseudoService struct {
	spec             job.Job
	webhookJobRunner *webhookJobRunner
}

// Start starts PseudoService.
func (s pseudoService) Start(context.Context) error {
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
	lggr          logger.Logger
}

func newWebhookJobRunner(runner pipeline.Runner, lggr logger.Logger) *webhookJobRunner {
	return &webhookJobRunner{
		specsByUUID: make(map[uuid.UUID]registeredJob),
		runner:      runner,
		lggr:        lggr.Named("JobRunner"),
	}
}

type registeredJob struct {
	job.Job
	chRemove services.StopChan
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

func (r *webhookJobRunner) RunJob(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error) {
	spec, exists := r.spec(jobUUID)
	if !exists {
		return 0, ErrJobNotExists
	}

	jobLggr := r.lggr.With(
		"jobID", spec.ID,
		"uuid", spec.ExternalJobID,
	)

	ctx, cancel := spec.chRemove.Ctx(ctx)
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

	_, err := r.runner.Run(ctx, run, jobLggr, true, nil)
	if err != nil {
		jobLggr.Errorw("Error running pipeline for webhook job", "err", err)
		return 0, err
	}
	if run.ID == 0 {
		panic("expected run to have non-zero id")
	}
	return run.ID, nil
}
