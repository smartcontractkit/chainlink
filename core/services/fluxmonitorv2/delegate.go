package fluxmonitorv2

import (
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// Delegate represents a Flux Monitor delegate
type Delegate struct {
	db             *sqlx.DB
	ethKeyStore    keystore.Eth
	jobORM         job.ORM
	pipelineORM    pipeline.ORM
	pipelineRunner pipeline.Runner
	legacyChains   legacyevm.LegacyChainContainer
	lggr           logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

// NewDelegate constructs a new delegate
func NewDelegate(
	ethKeyStore keystore.Eth,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	pipelineRunner pipeline.Runner,
	db *sqlx.DB,
	legacyChains legacyevm.LegacyChainContainer,
	lggr logger.Logger,
) *Delegate {
	return &Delegate{
		db:             db,
		ethKeyStore:    ethKeyStore,
		jobORM:         jobORM,
		pipelineORM:    pipelineORM,
		pipelineRunner: pipelineRunner,
		legacyChains:   legacyChains,
		lggr:           lggr.Named("FluxMonitor"),
	}
}

// JobType implements the job.Delegate interface
func (d *Delegate) JobType() job.Type {
	return job.FluxMonitor
}

func (d *Delegate) BeforeJobCreated(spec job.Job)                {}
func (d *Delegate) AfterJobCreated(spec job.Job)                 {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)                {}
func (d *Delegate) OnDeleteJob(spec job.Job, q pg.Queryer) error { return nil }

// ServicesForSpec returns the flux monitor service for the job spec
func (d *Delegate) ServicesForSpec(jb job.Job, qopts ...pg.QOpt) (services []job.ServiceCtx, err error) {
	if jb.FluxMonitorSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.FluxMonitorSpec to be present, got %v", jb)
	}
	chain, err := d.legacyChains.Get(jb.FluxMonitorSpec.EVMChainID.String())
	if err != nil {
		return nil, err
	}
	cfg := chain.Config()
	strategy := txmgrcommon.NewQueueingTxStrategy(jb.ExternalJobID, cfg.FluxMonitor().DefaultTransactionQueueDepth(), cfg.Database().DefaultQueryTimeout())
	var checker txmgr.TransmitCheckerSpec
	if chain.Config().FluxMonitor().SimulateTransactions() {
		checker.CheckerType = txmgr.TransmitCheckerTypeSimulate
	}

	fm, err := NewFromJobSpec(
		jb,
		d.db,
		NewORM(d.db, d.lggr, chain.Config().Database(), chain.TxManager(), strategy, checker),
		d.jobORM,
		d.pipelineORM,
		NewKeyStore(d.ethKeyStore),
		chain.Client(),
		chain.LogBroadcaster(),
		d.pipelineRunner,
		chain.Config().EVM(),
		chain.Config().EVM().GasEstimator(),
		chain.Config().EVM().Transactions(),
		chain.Config().FluxMonitor(),
		chain.Config().JobPipeline(),
		chain.Config().Database(),
		d.lggr,
	)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{fm}, nil
}
