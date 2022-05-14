package ocrbootstrap

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/relay/types"
)

// Delegate creates Bootstrap jobs
type Delegate struct {
	bootstrappers []commontypes.BootstrapperLocator
	db            *sqlx.DB
	jobORM        job.ORM
	peerWrapper   *ocrcommon.SingletonPeerWrapper
	cfg           validate.Config
	lggr          logger.Logger
	relayers      map[types.Network]types.Relayer
}

// NewDelegateBootstrap creates a new Delegate
func NewDelegateBootstrap(
	db *sqlx.DB,
	jobORM job.ORM,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	lggr logger.Logger,
	cfg validate.Config,
	relayers map[types.Network]types.Relayer,
) *Delegate {
	return &Delegate{
		db:          db,
		jobORM:      jobORM,
		peerWrapper: peerWrapper,
		lggr:        lggr,
		cfg:         cfg,
		relayers:    relayers,
	}
}

// JobType satisfies the job.Delegate interface.
func (d Delegate) JobType() job.Type {
	return job.Bootstrap
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d Delegate) ServicesForSpec(jobSpec job.Job) (services []job.ServiceCtx, err error) {
	spec := jobSpec.BootstrapSpec
	if spec == nil {
		return nil, errors.Errorf("Bootstrap.Delegate expects an *job.BootstrapSpec to be present, got %v", jobSpec)
	}
	if d.peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !d.peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	}
	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return nil, errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}
	configWatcher, err := relayer.NewConfigWatcher(types.ConfigWatcherArgs{
		JobID:       spec.ID,
		ContractID:  spec.ContractID,
		RelayConfig: spec.RelayConfig.Bytes(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relayer.NewOCR2Provider'")
	}
	lc := validate.ToLocalConfig(d.cfg, spec.AsOCR2Spec())
	if err = ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	d.lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)
	bootstrapNodeArgs := ocr.BootstrapperArgs{
		BootstrapperFactory:   d.peerWrapper.Peer2,
		ContractConfigTracker: configWatcher.ContractConfigTracker(),
		Database:              NewDB(d.db.DB, spec.ID, d.lggr),
		LocalConfig:           lc,
		Logger: logger.NewOCRWrapper(d.lggr.Named("OCR").With(
			"contractID", spec.ContractID,
			"jobName", jobSpec.Name.ValueOrZero(),
			"jobID", jobSpec.ID), true, func(msg string) {
			d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
		}),
		OffchainConfigDigester: configWatcher.OffchainConfigDigester(),
	}
	d.lggr.Debugw("Launching new bootstrap node", "args", bootstrapNodeArgs)
	bootstrapper, err := ocr.NewBootstrapper(bootstrapNodeArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewBootstrapNode")
	}
	return []job.ServiceCtx{configWatcher, job.NewServiceAdapter(bootstrapper)}, nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d Delegate) AfterJobCreated(spec job.Job) {
}

// BeforeJobDeleted satisfies the job.Delegate interface.
func (d Delegate) BeforeJobDeleted(spec job.Job) {
}
