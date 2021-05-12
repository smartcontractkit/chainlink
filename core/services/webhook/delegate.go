package webhook

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	Delegate struct {
		webhookJobRunner *webhookJobRunner
	}

	JobRunner interface {
		RunJob(jobUUID common.Hash, pipelineInputs []pipeline.Result, meta pipeline.JSONSerializable) (int64, error)
	}
)

func NewDelegate(runner pipeline.Runner) *Delegate {
	return &Delegate{
		webhookJobRunner: &webhookJobRunner{
			specsByUUID: make(map[common.Hash]registeredJob),
			runner:      runner,
		},
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Webhook
}

func (d *Delegate) ServicesForSpec(spec job.Job) ([]job.Service, error) {
	service := &pseudoService{
		spec:             spec,
		webhookJobRunner: d.webhookJobRunner,
	}
	return []job.Service{service}, nil
}

func (d *Delegate) WebhookJobRunner() JobRunner {
	return d.webhookJobRunner
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
	specsByUUID   map[common.Hash]registeredJob
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

func (r *webhookJobRunner) spec(jobID common.Hash) (registeredJob, bool) {
	r.muSpecsByUUID.RLock()
	defer r.muSpecsByUUID.RUnlock()
	spec, exists := r.specsByUUID[jobID]
	return spec, exists
}

var ErrJobNotExists = errors.New("job does not exist")

func (r *webhookJobRunner) RunJob(jobUUID common.Hash, pipelineInputs []pipeline.Result, meta pipeline.JSONSerializable) (int64, error) {
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

	ctx, cancel := utils.ContextFromChan(spec.chRemove)
	defer cancel()

	runID, _, err := r.runner.ExecuteAndInsertFinishedRun(ctx, *spec.PipelineSpec, pipelineInputs, meta, *logger, true)
	if err != nil {
		logger.Errorw("Error running pipeline for webhook job", "error", err)
		return 0, err
	}
	return runID, nil
}
