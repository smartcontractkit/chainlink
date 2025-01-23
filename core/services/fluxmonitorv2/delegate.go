package fluxmonitorv2

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type DelegateConfig interface {
	FluxMonitor() config.FluxMonitor
	JobPipeline() config.JobPipeline
}

// Delegate represents a Flux Monitor delegate
type Delegate struct {
	cfg            DelegateConfig
	ds             sqlutil.DataSource
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
	cfg DelegateConfig,
	ethKeyStore keystore.Eth,
	jobORM job.ORM,
	pipelineORM pipeline.ORM,
	pipelineRunner pipeline.Runner,
	ds sqlutil.DataSource,
	legacyChains legacyevm.LegacyChainContainer,
	lggr logger.Logger,
) *Delegate {
	return &Delegate{
		cfg:            cfg,
		ds:             ds,
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

func (d *Delegate) BeforeJobCreated(spec job.Job)              {}
func (d *Delegate) AfterJobCreated(spec job.Job)               {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)              {}
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec returns the flux monitor service for the job spec
func (d *Delegate) ServicesForSpec(ctx context.Context, jb job.Job) (services []job.ServiceCtx, err error) {
	if jb.FluxMonitorSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.FluxMonitorSpec to be present, got %v", jb)
	}
	chain, err := d.legacyChains.Get(jb.FluxMonitorSpec.EVMChainID.String())
	if err != nil {
		return nil, err
	}
	strategy := txmgrcommon.NewQueueingTxStrategy(jb.ExternalJobID, d.cfg.FluxMonitor().DefaultTransactionQueueDepth())
	var checker txmgr.TransmitCheckerSpec
	if d.cfg.FluxMonitor().SimulateTransactions() {
		checker.CheckerType = txmgr.TransmitCheckerTypeSimulate
	}

	fm, err := NewFromJobSpec(
		jb,
		d.ds,
		NewORM(d.ds, d.lggr, chain.TxManager(), strategy, checker),
		d.jobORM,
		d.pipelineORM,
		NewKeyStore(d.ethKeyStore),
		chain.Client(),
		chain.LogBroadcaster(),
		d.pipelineRunner,
		chain.Config().EVM(),
		chain.Config().EVM().GasEstimator(),
		d.cfg.JobPipeline(),
		d.lggr,
	)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{fm}, nil
}
