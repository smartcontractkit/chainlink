package fluxmonitorv2

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gorm.io/gorm"
)

// Delegate represents a Flux Monitor delegate
type Delegate struct {
	db              *gorm.DB
	ethKeyStore     *keystore.Eth
	jobORM          job.ORM
	pipelineORM     pipeline.ORM
	pipelineRunner  pipeline.Runner
	chainCollection evm.ChainCollection
}

var _ job.Delegate = (*Delegate)(nil)

// NewDelegate constructs a new delegate
func NewDelegate(
	ethKeyStore *keystore.Eth,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	pipelineRunner pipeline.Runner,
	db *gorm.DB,
	chainCollection evm.ChainCollection,
) *Delegate {
	return &Delegate{
		db,
		ethKeyStore,
		jobORM,
		pipelineORM,
		pipelineRunner,
		chainCollection,
	}
}

// JobType implements the job.Delegate interface
func (d *Delegate) JobType() job.Type {
	return job.FluxMonitor
}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the flux monitor service for the job spec
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.FluxMonitorSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.FluxMonitorSpec to be present, got %v", spec)
	}
	// TODO: Fix with https://app.clubhouse.io/chainlinklabs/story/14615/add-ability-to-set-chain-id-in-all-pipeline-tasks-that-interact-with-evm
	chain, err := d.chainCollection.Default()
	if err != nil {
		return nil, err
	}
	strategy := bulletprooftxmanager.NewQueueingTxStrategy(spec.ExternalJobID, chain.Config().FMDefaultTransactionQueueDepth())

	fm, err := NewFromJobSpec(
		spec,
		d.db,
		NewORM(d.db, chain.TxManager(), strategy),
		d.jobORM,
		d.pipelineORM,
		NewKeyStore(d.ethKeyStore),
		chain.Client(),
		chain.LogBroadcaster(),
		d.pipelineRunner,
		chain.Config(),
	)
	if err != nil {
		return nil, err
	}

	return []job.Service{fm}, nil
}
