package fluxmonitorv2

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	corestore "github.com/smartcontractkit/chainlink/core/store"
	"gorm.io/gorm"
)

// Delegate represents a Flux Monitor delegate
type Delegate struct {
	db             *gorm.DB
	store          *corestore.Store
	ethKeyStore    *keystore.Eth
	jobORM         job.ORM
	pipelineORM    pipeline.ORM
	pipelineRunner pipeline.Runner
	ethClient      eth.Client
	logBroadcaster log.Broadcaster
	cfg            Config
}

var _ job.Delegate = (*Delegate)(nil)

// NewDelegate constructs a new delegate
func NewDelegate(
	store *corestore.Store,
	ethKeyStore *keystore.Eth,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	pipelineRunner pipeline.Runner,
	db *gorm.DB,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
	cfg Config,
) *Delegate {
	return &Delegate{
		db,
		store,
		ethKeyStore,
		jobORM,
		pipelineORM,
		pipelineRunner,
		ethClient,
		logBroadcaster,
		cfg,
	}
}

// JobType implements the job.Delegate interface
func (d *Delegate) JobType() job.Type {
	return job.FluxMonitor
}

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

// ServicesForSpec returns the flux monitor service for the job spec
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.FluxMonitorSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.FluxMonitorSpec to be present, got %v", spec)
	}

	fm, err := NewFromJobSpec(
		spec,
		d.db,
		NewORM(d.store.DB),
		d.jobORM,
		d.pipelineORM,
		NewKeyStore(d.ethKeyStore),
		d.ethClient,
		d.logBroadcaster,
		d.pipelineRunner,
		d.cfg,
	)
	if err != nil {
		return nil, err
	}

	return []job.Service{fm}, nil
}
