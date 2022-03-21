package ocrbootstrap

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/relay"
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
	relayer       types.RelayerCtx
}

// NewDelegateBootstrap creates a new Delegate
func NewDelegateBootstrap(
	db *sqlx.DB,
	jobORM job.ORM,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	lggr logger.Logger,
	cfg validate.Config,
	relayer types.RelayerCtx,
) *Delegate {
	return &Delegate{
		db:          db,
		jobORM:      jobORM,
		peerWrapper: peerWrapper,
		lggr:        lggr,
		cfg:         cfg,
		relayer:     relayer,
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

	ocr2Provider, err := d.relayer.NewOCR2Provider(jobSpec.ExternalJobID, &relay.OCR2ProviderArgs{
		ID:              spec.ID,
		ContractID:      spec.ContractID,
		TransmitterID:   null.String{},
		Relay:           spec.Relay,
		RelayConfig:     spec.RelayConfig,
		IsBootstrapPeer: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relayer.NewOCR2Provider'")
	}
	services = append(services, ocr2Provider)

	configDB := NewDB(d.db.DB, spec.ID, d.lggr)
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	}

	loggerWith := d.lggr.Named("OCR").With(
		"contractID", spec.ContractID,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := logger.NewOCRWrapper(loggerWith, true, func(msg string) {
		d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
	})

	ocr2Spec := spec.AsOCR2Spec()
	lc := validate.ToLocalConfig(d.cfg, ocr2Spec)
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
	tracker := ocr2Provider.ContractConfigTracker()
	offchainConfigDigester := ocr2Provider.OffchainConfigDigester()

	bootstrapNodeArgs := ocr.BootstrapperArgs{
		BootstrapperFactory:    peerWrapper.Peer2,
		ContractConfigTracker:  tracker,
		Database:               configDB,
		LocalConfig:            lc,
		Logger:                 ocrLogger,
		OffchainConfigDigester: offchainConfigDigester,
	}

	d.lggr.Debugw("Launching new bootstrap node", "args", bootstrapNodeArgs)
	bootstrapper, err := ocr.NewBootstrapper(bootstrapNodeArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewBootstrapNode")
	}
	services = append(services, job.NewServiceAdapter(bootstrapper))

	return services, nil
}

// AfterJobCreated satisfies the job.Delegate interface.
func (d Delegate) AfterJobCreated(spec job.Job) {
}

// BeforeJobDeleted satisfies the job.Delegate interface.
func (d Delegate) BeforeJobDeleted(spec job.Job) {
}
