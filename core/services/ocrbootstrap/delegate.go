package ocrbootstrap

import (
	"context"

	"github.com/pkg/errors"

	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// Delegate creates Bootstrap jobs
type Delegate struct {
	db                *sqlx.DB
	jobORM            job.ORM
	peerWrapper       *ocrcommon.SingletonPeerWrapper
	cfg               validate.Config
	insecureCfg       validate.InsecureConfig
	lggr              logger.SugaredLogger
	relayers          map[relay.Network]loop.Relayer
	isNewlyCreatedJob bool
}

// NewDelegateBootstrap creates a new Delegate
func NewDelegateBootstrap(
	db *sqlx.DB,
	jobORM job.ORM,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	lggr logger.Logger,
	cfg validate.Config,
	insecureCfg validate.InsecureConfig,
	relayers map[relay.Network]loop.Relayer,
) *Delegate {
	return &Delegate{
		db:          db,
		jobORM:      jobORM,
		peerWrapper: peerWrapper,
		lggr:        logger.Sugared(lggr),
		cfg:         cfg,
		insecureCfg: insecureCfg,
		relayers:    relayers,
	}
}

// JobType satisfies the job.Delegate interface.
func (d *Delegate) JobType() job.Type {
	return job.Bootstrap
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	d.isNewlyCreatedJob = true
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(jobSpec job.Job) (services []job.ServiceCtx, err error) {
	spec := jobSpec.BootstrapSpec
	if spec == nil {
		return nil, errors.Errorf("Bootstrap.Delegate expects an *job.BootstrapSpec to be present, got %v", jobSpec)
	}
	if d.peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !d.peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running p2p v2 peer")
	}

	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return nil, errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}
	if spec.FeedID != nil {
		spec.RelayConfig["feedID"] = *spec.FeedID
	}

	ctxVals := loop.ContextValues{
		JobID:      jobSpec.ID,
		JobName:    jobSpec.Name.ValueOrZero(),
		ContractID: spec.ContractID,
		FeedID:     spec.FeedID,
	}
	ctx := ctxVals.ContextWithValues(context.Background())

	configProvider, err := relayer.NewConfigProvider(ctx, types.RelayArgs{
		ExternalJobID: jobSpec.ExternalJobID,
		JobID:         spec.ID,
		ContractID:    spec.ContractID,
		New:           d.isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relayer.NewConfigWatcher'")
	}
	lc := validate.ToLocalConfig(d.cfg, d.insecureCfg, spec.AsOCR2Spec())
	if err = ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	lggr := d.lggr.With(ctxVals.Args()...)
	lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)
	bootstrapNodeArgs := ocr.BootstrapperArgs{
		BootstrapperFactory:   d.peerWrapper.Peer2,
		ContractConfigTracker: configProvider.ContractConfigTracker(),
		Database:              NewDB(d.db.DB, spec.ID, lggr),
		LocalConfig:           lc,
		Logger: logger.NewOCRWrapper(lggr.Named("OCRBootstrap"), true, func(msg string) {
			logger.Sugared(lggr).ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
		}),
		OffchainConfigDigester: configProvider.OffchainConfigDigester(),
	}
	lggr.Debugw("Launching new bootstrap node", "args", bootstrapNodeArgs)
	bootstrapper, err := ocr.NewBootstrapper(bootstrapNodeArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewBootstrapNode")
	}
	return []job.ServiceCtx{configProvider, job.NewServiceAdapter(bootstrapper)}, nil
}

// OnCreateJob satisfies the job.Delegate interface.
func (d *Delegate) OnCreateJob(jb job.Job, q pg.Queryer) error {
	spec := jb.BootstrapSpec
	if spec == nil {
		return errors.Errorf("Bootstrap.Delegate expects an *job.BootstrapSpec to be present, got %v", spec)
	}

	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}

	rargs := types.RelayArgs{
		JobID:       spec.ID,
		ContractID:  spec.ContractID,
		RelayConfig: spec.RelayConfig.Bytes(),
	}
	r, ok := relayer.(relay.LogPollerCapable)
	if !ok {
		return nil
	}
	err := r.RegisterLogFilters(rargs, q)
	if err != nil && !errors.Is(err, relay.ErrLogFiltersNotSupported{}) {
		return errors.Wrapf(err, "Failed to register required log filters for OCRBootstrap job, with relay args: %v", rargs)
	}
	return nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d *Delegate) AfterJobCreated(spec job.Job) {
}

// BeforeJobDeleted satisfies the job.Delegate interface.
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// OnDeleteJob satisfies the job.Delegate interface.
func (d *Delegate) OnDeleteJob(jb job.Job, q pg.Queryer) error {
	spec := jb.BootstrapSpec
	if spec == nil {
		return errors.Errorf("Bootstrap.Delegate expects an *job.BootstrapSpec to be present, got %v", spec)
	}

	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}

	rargs := types.RelayArgs{
		JobID:       spec.ID,
		ContractID:  spec.ContractID,
		RelayConfig: spec.RelayConfig.Bytes(),
	}
	r, ok := relayer.(relay.LogPollerCapable)
	if !ok {
		return nil
	}
	err := r.UnregisterLogFilters(rargs, q)
	if err != nil && !errors.Is(err, relay.ErrLogFiltersNotSupported{}) {
		return errors.Wrapf(err, "Failed to unregister log filters for OCRBootstrap job, with relay args: %v", rargs)
	}
	return nil
}
